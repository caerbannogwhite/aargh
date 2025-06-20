package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/caerbannogwhite/aargh/meta"
)

const (
	GOROUTINES                = 4
	RESULT_VAR_NAME           = "result"
	RESULT_SIZE_VAR_NAME      = "resultSize"
	RESULT_NULL_MASK_VAR_NAME = "resultNullMask"
	FINAL_RETURN_FMT          = "Errors{fmt.Sprintf(\"Cannot %s %%s and %%s\", s.Type().String(), o.Type().String())}"
)

type BuildInfo struct {
	OpCode        meta.OPCODE
	Op1Nullable   bool
	Op1Scalar     bool
	Op2Nullable   bool
	Op2Scalar     bool
	Op1VarName    string
	Op1SeriesType string
	Op1InnerType  meta.BaseType
	Op2VarName    string
	Op2SeriesType string
	Op2InnerType  meta.BaseType
	ResInnerType  meta.BaseType
	MakeOperation MakeOperationType
}

func (bi BuildInfo) UpdateScalarInfo(Op1Scalar, Op2Scalar bool) BuildInfo {
	bi.Op1Scalar = Op1Scalar
	bi.Op2Scalar = Op2Scalar
	return bi
}

func (bi BuildInfo) UpdateNullableInfo(Op1Nullable, Op2Nullable bool) BuildInfo {
	bi.Op1Nullable = Op1Nullable
	bi.Op2Nullable = Op2Nullable
	return bi
}

// Generate the code to define the result inner array
// and to compute the result size and null mask
func generateMakeResultStmt(info BuildInfo) []ast.Stmt {
	var resSizeVariable string

	if info.ResInnerType == info.Op1InnerType {
		if info.Op1Scalar {
			resSizeVariable = info.Op2VarName
		} else {
			resSizeVariable = info.Op1VarName
		}
	} else {
		if info.Op1Scalar {
			resSizeVariable = info.Op2VarName
		} else {
			resSizeVariable = info.Op1VarName
		}
	}

	sizeCase := 0
	if info.Op1Scalar {
		if info.Op2Scalar {
			sizeCase = 0
		} else {
			sizeCase = 1
		}
	} else {
		if info.Op2Scalar {
			sizeCase = 2
		} else {
			sizeCase = 3
		}
	}

	resultGoType := info.ResInnerType.ToGoType()

	// Special case for the result type
	if info.ResInnerType == meta.StringType {
		resultGoType = "[]*string"
	}

	// assign the result size
	stmts := []ast.Stmt{&ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: RESULT_SIZE_VAR_NAME},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.Ident{Name: fmt.Sprintf("%s.Len()", resSizeVariable)},
		},
	}}

	// One of the operands is NAs, take the null mask of the other operand
	// if info.Op1InnerType == meta.NullType || info.Op2InnerType == meta.NullType {
	// 	stmts = append(stmts, &ast.AssignStmt{
	// 		Lhs: []ast.Expr{
	// 			&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
	// 		},
	// 		Tok: token.DEFINE,
	// 		Rhs: []ast.Expr{
	// 			&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", resSizeVariable)},
	// 		},
	// 	})
	// }

	if info.ResInnerType != meta.NullType {
		stmts = append(stmts,

			// make the result array
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{Name: RESULT_VAR_NAME},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{Name: "make"},
						Args: []ast.Expr{
							&ast.Ident{Name: resultGoType},
							&ast.Ident{Name: RESULT_SIZE_VAR_NAME},
						},
					},
				},
			})

		// Make the result null mask

		// Special case: one of the operands is NAs
		if info.Op1InnerType == meta.NullType || info.Op2InnerType == meta.NullType {

			nonNullOperand := info.Op1VarName
			nonNullOperandIsScalar := info.Op1Scalar
			if info.Op1InnerType == meta.NullType {
				nonNullOperand = info.Op2VarName
				nonNullOperandIsScalar = info.Op2Scalar
			}

			stmts = append(stmts, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								{Name: RESULT_NULL_MASK_VAR_NAME},
							},
							Type: &ast.Ident{Name: "[]uint8"},
						},
					},
				},
			})

			// The non-null operand is a scalar
			if nonNullOperandIsScalar {
				stmts = append(stmts, &ast.IfStmt{
					Cond: ast.NewIdent(fmt.Sprintf("%s.IsNullable_", nonNullOperand)),
					Body: &ast.BlockStmt{List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{Name: fmt.Sprintf("utils.BinVecInit(%s, %s.NullMask_[0] == 1)", RESULT_SIZE_VAR_NAME, nonNullOperand)},
							}},
					}},
					Else: &ast.BlockStmt{List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{Name: "make([]uint8, 0)"},
							}},
					}},
				})
			} else

			// The non-null operand is vector
			{
				stmts = append(stmts,
					&ast.IfStmt{
						Cond: ast.NewIdent(fmt.Sprintf("%s.IsNullable_", nonNullOperand)),
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: fmt.Sprintf("utils.BinVecInit(%s, %s)", RESULT_SIZE_VAR_NAME, "false")},
								}},
							&ast.ExprStmt{X: &ast.CallExpr{
								Fun: &ast.Ident{Name: "copy"},
								Args: []ast.Expr{
									&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
									&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", nonNullOperand)},
								}}},
						}},
						Else: &ast.BlockStmt{List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: "make([]uint8, 0)"},
								}},
						}},
					})
			}
		} else

		// Default: check the nullability of the operands
		if info.Op1Nullable {
			if info.Op2Nullable {

				// Both operands are nullable:
				// call the binary vector or function to merge the null masks
				stmts = append(stmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.Ident{Name: fmt.Sprintf("utils.BinVecInit(%s, false)", RESULT_SIZE_VAR_NAME)},
					},
				})

				funcName := "utils.BinVecOrSS"
				switch sizeCase {
				case 0:
					funcName = "utils.BinVecOrSS"
				case 1:
					funcName = "utils.BinVecOrSV"
				case 2:
					funcName = "utils.BinVecOrVS"
				case 3:
					funcName = "utils.BinVecOrVV"
				}

				stmts = append(stmts, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: funcName},
						Args: []ast.Expr{
							&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", info.Op1VarName)},
							&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", info.Op2VarName)},
							&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
						},
					},
				})
			} else {

				// Only the second operand is nullable, so the reuslt null mask
				// depends on the value of the second operand null mask

				// 	1 - initialize the null mask to 0 or, if the second operand is a scalar,
				// 		to the value of its null mask
				nullMaskInitFlag := "false"
				if info.Op1Scalar {
					nullMaskInitFlag = fmt.Sprintf("%s.NullMask_[0] == 1", info.Op1VarName)
				}

				// 	2 - call the binary vector init function
				stmts = append(stmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.Ident{Name: fmt.Sprintf("utils.BinVecInit(%s, %s)", RESULT_SIZE_VAR_NAME, nullMaskInitFlag)},
					},
				})

				// 	3 - if the first operand is not a scalar, copy its null mask
				if !info.Op1Scalar {
					stmts = append(stmts, &ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.Ident{Name: "copy"},
						Args: []ast.Expr{
							&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
							&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", info.Op1VarName)},
						}},
					})
				}
			}
		} else {
			if info.Op2Nullable {

				// Only the second operand is nullable, so the reuslt null mask
				// depends on the value of the second operand null mask

				// 	1 - initialize the null mask to 0 or, if the second operand is a scalar,
				// 		to the value of its null mask
				nullMaskInitFlag := "false"
				if info.Op2Scalar {
					nullMaskInitFlag = fmt.Sprintf("%s.NullMask_[0] == 1", info.Op2VarName)
				}

				// 	2 - call the binary vector init function
				stmts = append(stmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.Ident{Name: fmt.Sprintf("utils.BinVecInit(%s, %s)", RESULT_SIZE_VAR_NAME, nullMaskInitFlag)},
					},
				})

				// 	3 - if the second operand is not a scalar, copy its null mask
				if !info.Op2Scalar {
					stmts = append(stmts, &ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.Ident{Name: "copy"},
						Args: []ast.Expr{
							&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
							&ast.Ident{Name: fmt.Sprintf("%s.NullMask_", info.Op2VarName)},
						}},
					})
				}
			} else {

				// None of the operands is nullable:
				// initialize the null mask to 0
				stmts = append(stmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.Ident{Name: "utils.BinVecInit(0, false)"},
					},
				})
			}
		}
	}

	return stmts
}

// Generate the code to compute the operation
func generateOperationLoop(info BuildInfo) []ast.Stmt {

	statements := make([]ast.Stmt, 0)

	if info.Op1Scalar && info.Op2Scalar {
		statements = append(statements, &ast.ExprStmt{
			X: info.MakeOperation(RESULT_VAR_NAME, "0", info.Op1VarName, "0", info.Op2VarName, "0"),
		})
	} else {

		op1Index := "i"
		op2Index := "i"
		if info.Op1Scalar {
			op1Index = "0"
		}

		if info.Op2Scalar {
			op2Index = "0"
		}

		statements = append(statements, &ast.ForStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{Name: "i"},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.Ident{Name: "0"},
				},
			},
			Cond: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "i"},
				Op: token.LSS,
				Y:  &ast.Ident{Name: RESULT_SIZE_VAR_NAME},
			},
			Post: &ast.IncDecStmt{
				X:   &ast.Ident{Name: "i"},
				Tok: token.INC,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: info.MakeOperation(RESULT_VAR_NAME, "i", info.Op1VarName, op1Index, info.Op2VarName, op2Index),
					},
				},
			},
		})
	}

	return statements
}

func generateOperation(info BuildInfo) []ast.Stmt {
	resIsNullable := info.Op1Nullable || info.Op2Nullable
	resSeriesType := computeResSeriesType(info.OpCode, info.Op1InnerType, info.Op2InnerType)

	statements := make([]ast.Stmt, 0)

	// 1 - Generate the result inner data array
	statements = append(statements, generateMakeResultStmt(info)...)

	// 2 - Generate the loop to compute the operation
	if resSeriesType != "NAs" {
		statements = append(statements, generateOperationLoop(info)...)
	}

	// 3 - Generate the return statement with the result series
	params := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "IsNullable_"},
			Value: &ast.Ident{Name: fmt.Sprintf("%v", resIsNullable)},
		},
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "NullMask_"},
			Value: &ast.Ident{Name: RESULT_NULL_MASK_VAR_NAME},
		},
	}

	switch resSeriesType {

	// NA: the only parameter is the size of the result series
	case "NAs":
		params = []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "size"},
				Value: &ast.Ident{Name: RESULT_SIZE_VAR_NAME},
			},
		}

	// BOOL Memory optimized: convert the result to a binary vector and add the size to the result series
	case "SeriesBoolMemOpt":
		params = append(params, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "Data_"},
			Value: &ast.Ident{Name: fmt.Sprintf("boolVecToBinVec(%s)", RESULT_VAR_NAME)},
		})

		params = append(params, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "size"},
			Value: &ast.Ident{Name: RESULT_SIZE_VAR_NAME},
		})

	// Default: just add the data to the result series
	default:
		params = append(params, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "Data_"},
			Value: &ast.Ident{Name: RESULT_VAR_NAME},
		})

		params = append(params, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "Ctx_"},
			Value: &ast.Ident{Name: fmt.Sprintf("%s.Ctx_", info.Op1VarName)},
		})
	}

	statements = append(statements, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.Ident{Name: resSeriesType},
				Elts: params,
			},
		},
	})

	return statements
}

// Generate the if statement to check the nullability of the operands
func generateNullabilityCheck(info BuildInfo) []ast.Stmt {

	// If one of the operands is nullable, just generate the operation
	// There is no need to check the nullability of the operands
	if info.Op1InnerType == meta.NullType || info.Op2InnerType == meta.NullType {
		return generateOperation(info)
	} else {
		return []ast.Stmt{
			&ast.IfStmt{
				Cond: ast.NewIdent(fmt.Sprintf("%s.IsNullable_", info.Op1VarName)),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: ast.NewIdent(fmt.Sprintf("%s.IsNullable_", info.Op2VarName)),
							Body: &ast.BlockStmt{
								List: generateOperation(info.UpdateNullableInfo(true, true)),
							},
							Else: &ast.BlockStmt{
								List: generateOperation(info.UpdateNullableInfo(true, false)),
							},
						},
					},
				},
				Else: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: ast.NewIdent(fmt.Sprintf("%s.IsNullable_", info.Op2VarName)),
							Body: &ast.BlockStmt{
								List: generateOperation(info.UpdateNullableInfo(false, true)),
							},
							Else: &ast.BlockStmt{
								List: generateOperation(info.UpdateNullableInfo(false, false)),
							},
						},
					},
				},
			},
		}
	}
}

// Generate the if statement to check the size of the series
func generateSizeCheck(info BuildInfo, defaultReturn ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: fmt.Sprintf("%s.Len()", info.Op1VarName)},
			Op: token.EQL,
			Y:  &ast.Ident{Name: "1"},
		},

		// CASE OP1_SIZE == 1
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: fmt.Sprintf("%s.Len()", info.Op2VarName)},
						Op: token.EQL,
						Y:  &ast.Ident{Name: "1"},
					},

					// CASE OP1_SIZE == 1 AND OP2_SIZE == 1
					Body: &ast.BlockStmt{
						List: generateNullabilityCheck(info.UpdateScalarInfo(true, true)),
					},

					// CASE OP1_SIZE == 1 AND OP2_SIZE != 1
					Else: &ast.BlockStmt{
						List: generateNullabilityCheck(info.UpdateScalarInfo(true, false)),
					},
				},
			},
		},

		// CASE OP1_SIZE != 1
		Else: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: fmt.Sprintf("%s.Len()", info.Op2VarName)},
						Op: token.EQL,
						Y:  &ast.Ident{Name: "1"},
					},

					// CASE OP1_SIZE != 1 AND OP2_SIZE == 1
					Body: &ast.BlockStmt{
						List: generateNullabilityCheck(info.UpdateScalarInfo(false, true)),
					},

					// CASE OP1_SIZE != 1 AND OP2_SIZE != 1
					Else: &ast.IfStmt{
						Cond: &ast.BinaryExpr{
							X:  &ast.Ident{Name: fmt.Sprintf("%s.Len()", info.Op1VarName)},
							Op: token.EQL,
							Y:  &ast.Ident{Name: fmt.Sprintf("%s.Len()", info.Op2VarName)},
						},

						Body: &ast.BlockStmt{
							List: generateNullabilityCheck(info.UpdateScalarInfo(false, false)),
						},
					},
				},

				defaultReturn,
			},
		},
	}
}

// Generate the switch statement to handle the different types of the second operand
func generateSwitchType(
	operation Operation, op1SeriesType string, op1InnerType meta.BaseType,
	op1VarName, op2VarName string, defaultReturn ast.Stmt) []ast.Stmt {

	// Generate the preliminary type check, to check the type of second operand
	// is a Series or a raw value
	otherSeriesDefiniton := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{Name: "otherSeries"}},
					Type:  &ast.Ident{Name: "Series"},
				},
			},
		},
	}

	typeCheck := &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("_, ok")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{ast.NewIdent("other.(Series)")},
		},
		Cond: &ast.Ident{Name: "ok"},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("otherSeries")},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{ast.NewIdent("other.(Series)")},
				},
			},
		},
		Else: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("otherSeries")},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{ast.NewIdent("NewSeries(other, nil, false, false, s.Ctx_)")},
				},
			},
		},
	}

	// Generate the context check
	contextCheck := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: fmt.Sprintf("%s.Ctx_", op1VarName)},
			Op: token.NEQ,
			Y:  &ast.Ident{Name: fmt.Sprintf("%s.GetContext()", "otherSeries")},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{ast.NewIdent(fmt.Sprintf("Errors{fmt.Sprintf(\"Cannot operate on series with different contexts: %%v and %%v\", s.Ctx_, %s.GetContext())}", "otherSeries"))},
				},
			},
		},
	}

	op2VarNameTyped := "o"
	bigSwitch := &ast.TypeSwitchStmt{
		Assign: &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(op2VarNameTyped)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{ast.NewIdent(fmt.Sprintf("%s.(type)", "otherSeries"))},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	// Generate the switch cases for each type of the second operand
	for _, op2 := range operation.ApplyTo {
		bigSwitch.Body.List = append(bigSwitch.Body.List,
			&ast.CaseClause{
				List: []ast.Expr{ast.NewIdent(op2.SeriesName)},
				Body: []ast.Stmt{
					generateSizeCheck(BuildInfo{
						OpCode:        operation.OpCode,
						Op1VarName:    op1VarName,
						Op1SeriesType: op1SeriesType,
						Op1InnerType:  op1InnerType,
						Op2VarName:    op2VarNameTyped,
						Op2SeriesType: op2.SeriesName,
						Op2InnerType:  op2.SeriesType,
						ResInnerType:  ComputeResInnerType(operation.OpCode, op1InnerType, op2.SeriesType),
						MakeOperation: op2.MakeOperation,
					}, defaultReturn),
				},
			},
		)
	}

	bigSwitch.Body.List = append(bigSwitch.Body.List, &ast.CaseClause{
		List: nil,
		Body: []ast.Stmt{defaultReturn},
	})

	return []ast.Stmt{otherSeriesDefiniton, typeCheck, contextCheck, bigSwitch}
}

func computeResSeriesType(opCode meta.OPCODE, op1, op2 meta.BaseType) string {
	switch ComputeResInnerType(opCode, op1, op2) {
	case meta.NullType:
		return "NAs"
	case meta.BoolType:
		return "Bools"
	case meta.IntType:
		return "Ints"
	case meta.Int64Type:
		return "Int64s"
	case meta.Float32Type:
		return "SeriesFloat32"
	case meta.Float64Type:
		return "Float64s"
	case meta.StringType:
		return "Strings"
	case meta.TimeType:
		return "Times"
	case meta.DurationType:
		return "Durations"
	}
	return "Errors"
}

func ComputeResInnerType(opCode meta.OPCODE, op1, op2 meta.BaseType) meta.BaseType {
	return opCode.GetBinaryOpResultType(meta.Primitive{Base: op1}, meta.Primitive{Base: op2}).Base
}

func generateOperations() {
	for filename, info := range GenerateOperationsData() {

		src, err := os.ReadFile(filepath.Join(SERIES_FOLDER, filename))
		if err != nil {
			panic(err)
		}

		// Parse the file.
		fset := token.NewFileSet()
		fast, err := parser.ParseFile(fset, filepath.Join(SERIES_FOLDER, filename), src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// // Add the utils package
		// fast.Decls = append(fast.Decls, &ast.GenDecl{
		// 	Tok: token.IMPORT,
		// 	Specs: []ast.Spec{
		// 		&ast.ImportSpec{Path: &ast.BasicLit{Value: `"github.com/caerbannogwhite/aargh/utils"`}},
		// 	},
		// })

		for i, decl := range fast.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				switch funcDecl.Name.Name {
				case "And":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["And"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "AND"))},
						})

				case "Or":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Or"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "OR"))},
						})

				case "Mul":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Mul"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "multiply"))},
						})

				case "Div":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Div"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "divide"))},
						})

				case "Mod":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Mod"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "use modulo"))},
						})

				case "Exp":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Exp"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "use exponentiation"))},
						})

				case "Add":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Add"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "sum"))},
						})

				case "Sub":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Sub"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "subtract"))},
						})

				case "Eq":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Eq"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for equality"))},
						})

				case "Ne":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Ne"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for inequality"))},
						})

				case "Lt":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Lt"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for less than"))},
						})

				case "Le":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Le"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for less than or equal to"))},
						})

				case "Gt":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Gt"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for greater than"))},
						})

				case "Ge":
					fast.Decls[i].(*ast.FuncDecl).Body.List = generateSwitchType(
						info.Operations["Ge"], info.SeriesName, info.SeriesType, "s", "other",
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent(fmt.Sprintf(FINAL_RETURN_FMT, "compare for greater than or equal to"))},
						})
				}
			}

			buf := new(bytes.Buffer)
			err = format.Node(buf, fset, fast)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(filepath.Join(SERIES_FOLDER, filename), buf.Bytes(), 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}

func generateBase() {
	for filename, info := range DATA_BASE_METHODS {
		tmplBase, err := template.New("base").Parse(TEMPLATE_BASIC_ACCESSORS)
		if err != nil {
			panic(err)
		}

		tmplFilters, err := template.New("filters").Parse(TEMPLATE_FILTERS)
		if err != nil {
			panic(err)
		}

		tmplMaps, err := template.New("maps").Parse(TEMPLATE_MAPS)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(SERIES_FOLDER, filename))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = tmplBase.Execute(f, info)
		if err != nil {
			panic(err)
		}

		err = tmplFilters.Execute(f, info)
		if err != nil {
			panic(err)
		}

		err = tmplMaps.Execute(f, info)
		if err != nil {
			panic(err)
		}
	}
}

const SERIES_FOLDER = "..\\series"

func main() {
	generateBase()
	generateOperations()
}
