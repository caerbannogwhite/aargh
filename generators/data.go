package main

import (
	"fmt"
	"go/ast"

	"github.com/caerbannogwhite/enchanter/meta"
)

type MakeOperationType func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr

type OperationApplyTo struct {
	SeriesName    string
	SeriesType    meta.BaseType
	MakeOperation MakeOperationType
}

type Operation struct {
	OpCode  meta.OPCODE
	ApplyTo []OperationApplyTo
}

type SeriesFile struct {
	SeriesName            string
	SeriesType            meta.BaseType
	SeriesTypeStr         string
	SeriesGoTypeStr       string
	SeriesGoOuterTypeStr  string
	SeriesNullableTypeStr string
	DefaultValue          string
	IsGoTypePtr           bool
	IsTimeType            bool
	Operations            map[string]Operation
}

var DATA_BASE_METHODS = map[string]SeriesFile{
	"bool_base.go": {
		SeriesName:            "Bools",
		SeriesTypeStr:         "BoolType",
		SeriesGoTypeStr:       "bool",
		SeriesGoOuterTypeStr:  "bool",
		SeriesNullableTypeStr: "enchanter.NullableBool",
		DefaultValue:          "false",
	},

	"int_base.go": {
		SeriesName:            "Ints",
		SeriesTypeStr:         "IntType",
		SeriesGoTypeStr:       "int",
		SeriesGoOuterTypeStr:  "int",
		SeriesNullableTypeStr: "enchanter.NullableInt",
		DefaultValue:          "0",
	},

	"int64_base.go": {
		SeriesName:            "Int64s",
		SeriesTypeStr:         "Int64Type",
		SeriesGoTypeStr:       "int64",
		SeriesGoOuterTypeStr:  "int64",
		SeriesNullableTypeStr: "enchanter.NullableInt64",
		DefaultValue:          "0",
	},

	"float64_base.go": {
		SeriesName:            "Float64s",
		SeriesTypeStr:         "Float64Type",
		SeriesGoTypeStr:       "float64",
		SeriesGoOuterTypeStr:  "float64",
		SeriesNullableTypeStr: "enchanter.NullableFloat64",
		DefaultValue:          "0",
	},

	"string_base.go": {
		SeriesName:            "Strings",
		SeriesTypeStr:         "StringType",
		SeriesGoTypeStr:       "*string",
		SeriesGoOuterTypeStr:  "string",
		SeriesNullableTypeStr: "enchanter.NullableString",
		DefaultValue:          "s.Ctx_.StringPool.Put(enchanter.NA_TEXT)",
		IsGoTypePtr:           true,
	},

	"time_base.go": {
		SeriesName:            "Times",
		SeriesTypeStr:         "TimeType",
		SeriesGoTypeStr:       "time.Time",
		SeriesGoOuterTypeStr:  "time.Time",
		SeriesNullableTypeStr: "enchanter.NullableTime",
		DefaultValue:          "time.Time{}",
		IsTimeType:            true,
	},

	"duration_base.go": {
		SeriesName:            "Durations",
		SeriesTypeStr:         "DurationType",
		SeriesGoTypeStr:       "time.Duration",
		SeriesGoOuterTypeStr:  "time.Duration",
		SeriesNullableTypeStr: "enchanter.NullableDuration",
		DefaultValue:          "time.Duration(0)",
	},
}

func GenerateOperationsData() map[string]SeriesFile {
	var data = map[string]SeriesFile{
		"na_ops.go": {
			SeriesName: "NAs",
			SeriesType: meta.NullType,
			Operations: map[string]Operation{},
		},

		"bool_ops.go": {
			SeriesName: "Bools",
			SeriesType: meta.BoolType,
			Operations: map[string]Operation{
				"Mul": {
					OpCode: meta.OP_BINARY_MUL,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] && %s.Data_[%s] { %s[%s] = 1 }", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Div": {
					OpCode: meta.OP_BINARY_DIV,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nb2 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 / b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 / float64(%s.Data_[%s])", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 / float64(%s.Data_[%s])", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 / %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Mod": {
					OpCode: meta.OP_BINARY_MOD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nb2 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = math.Mod(b1, b2)", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = math.Mod(b1, float64(%s.Data_[%s]))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = math.Mod(b1, float64(%s.Data_[%s]))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = math.Mod(b1, float64(%s.Data_[%s]))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Exp": {
					OpCode: meta.OP_BINARY_EXP,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nb2 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = int64(math.Pow(b1, b2))", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = int64(math.Pow(b1, float64(%s.Data_[%s])))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = int64(math.Pow(b1, float64(%s.Data_[%s])))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = float64(math.Pow(b1, float64(%s.Data_[%s])))", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 + b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 + %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 + %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 + %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(boolToString(%s.Data_[%s]) + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 - b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 - %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 - %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 - %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 < b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 < %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 < %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 < %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 <= b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 <= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 <= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 <= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 > b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 > %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 > %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 > %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nb2 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = b1 >= b2", op1, op1Index, op2, op2Index, res, resIndex)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 >= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := int64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 >= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b1 := float64(0)\nif %s.Data_[%s] { b1 = 1 }\n%s[%s] = b1 >= %s.Data_[%s]", op1, op1Index, res, resIndex, op2, op2Index)}
							},
						},
					},
				},

				"And": {
					OpCode: meta.OP_BINARY_AND,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] && %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Or": {
					OpCode: meta.OP_BINARY_OR,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] || %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"int_ops.go": {
			SeriesName: "Ints",
			SeriesType: meta.IntType,
			Operations: map[string]Operation{
				"Mul": {
					OpCode: meta.OP_BINARY_MUL,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Div": {
					OpCode: meta.OP_BINARY_DIV,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = float64(%s.Data_[%s]) / b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Mod": {
					OpCode: meta.OP_BINARY_MOD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = math.Mod(float64(%s.Data_[%s]), b2)", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Exp": {
					OpCode: meta.OP_BINARY_EXP,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), b2))", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s])))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s])))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] + b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(intToString(int64(%s.Data_[%s])) + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] - b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] < b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] < %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) < %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) < %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] <= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] <= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) <= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) <= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] > b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] > %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) > %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) > %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] >= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] >= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(%s.Data_[%s]) >= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) >= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"int64_ops.go": {
			SeriesName: "Int64s",
			SeriesType: meta.Int64Type,
			Operations: map[string]Operation{
				"Mul": {
					OpCode: meta.OP_BINARY_MUL,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * int64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Div": {
					OpCode: meta.OP_BINARY_DIV,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = float64(%s.Data_[%s]) / b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) / %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Mod": {
					OpCode: meta.OP_BINARY_MOD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = math.Mod(float64(%s.Data_[%s]), b2)", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Exp": {
					OpCode: meta.OP_BINARY_EXP,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), b2))", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s])))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = int64(math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s])))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Pow(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] + b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + int64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(intToString(%s.Data_[%s]) + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] - b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - int64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] == int64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] != int64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = float64(%s.Data_[%s]) != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] < b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] < int64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] < %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = float64(%s.Data_[%s]) < %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] <= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] <= int64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] <= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = float64(%s.Data_[%s]) <= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] > b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] > int64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] > %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = float64(%s.Data_[%s]) > %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := int64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] >= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] >= int64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] >= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = float64(%s.Data_[%s]) >= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"float64_ops.go": {
			SeriesName: "Float64s",
			SeriesType: meta.Float64Type,
			Operations: map[string]Operation{
				"Mul": {
					OpCode: meta.OP_BINARY_MUL,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("if %s.Data_[%s] { %s[%s] = %s.Data_[%s] }", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] * %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Div": {
					OpCode: meta.OP_BINARY_DIV,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] / b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] / float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] / %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Mod": {
					OpCode: meta.OP_BINARY_MOD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = math.Mod(%s.Data_[%s], b2)", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Mod(float64(%s.Data_[%s]), float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Exp": {
					OpCode: meta.OP_BINARY_EXP,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = math.Pow(%s.Data_[%s], b2)", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Pow(%s.Data_[%s], float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Pow(%s.Data_[%s], float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = math.Pow(%s.Data_[%s], float64(%s.Data_[%s]))", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] + b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(floatToString(%s.Data_[%s]) + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] - b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - float64(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] == float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] == float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] == %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] != float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] != float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] != %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] < b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] < float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] < float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] < %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] <= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] <= float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] <= float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] <= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] > b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] > float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] > float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] > %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("b2 := float64(0)\nif %s.Data_[%s] { b2 = 1 }\n%s[%s] = %s.Data_[%s] >= b2", op2, op2Index, res, resIndex, op1, op1Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] >= float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] >= float64(%s.Data_[%s])`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = %s.Data_[%s] >= %s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"string_ops.go": {
			SeriesName: "Strings",
			SeriesType: meta.StringType,
			Operations: map[string]Operation{
				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Bools",
							SeriesType: meta.BoolType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + boolToString(%s.Data_[%s]))", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Ints",
							SeriesType: meta.IntType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + intToString(int64(%s.Data_[%s])))", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Int64s",
							SeriesType: meta.Int64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + intToString(%s.Data_[%s]))", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Float64s",
							SeriesType: meta.Float64Type,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + floatToString(%s.Data_[%s]))", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + *%s.Data_[%s])", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + %s.Data_[%s].String())", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + %s.Data_[%s].String())", res, resIndex, op1, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] == *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] != *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] < *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] <= *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] > *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf(`%s[%s] = *%s.Data_[%s] >= *%s.Data_[%s]`, res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"time_ops.go": {
			SeriesName: "Times",
			SeriesType: meta.TimeType,
			Operations: map[string]Operation{
				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(%s.Data_[%s].String() + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].AddDate(%s.Data_[%s].Year(), int(%s.Data_[%s].Month()), %s.Data_[%s].Day())", res, resIndex, op1, op1Index, op2, op2Index, op2, op2Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Add(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Sub(%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Add(-%s.Data_[%s])", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) == 0", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) != 0", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) == -1", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) <= 0", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) == 1", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Compare(%s.Data_[%s]) >= 1", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},

		"duration_ops.go": {
			SeriesName: "Durations",
			SeriesType: meta.DurationType,
			Operations: map[string]Operation{
				"Add": {
					OpCode: meta.OP_BINARY_ADD,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Strings",
							SeriesType: meta.StringType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(%s.Data_[%s].String() + *%s.Data_[%s])", res, resIndex, op2, op1, op1Index, op2, op2Index)}
							},
						},
						{
							SeriesName: "Times",
							SeriesType: meta.TimeType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s].Add(%s.Data_[%s])", res, resIndex, op2, op2Index, op1, op1Index)}
							},
						},
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] + %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Sub": {
					OpCode: meta.OP_BINARY_SUB,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] - %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Eq": {
					OpCode: meta.OP_BINARY_EQ,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] == %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ne": {
					OpCode: meta.OP_BINARY_NE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",

							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] != %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Lt": {
					OpCode: meta.OP_BINARY_LT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] < %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Le": {
					OpCode: meta.OP_BINARY_LE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] <= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Gt": {
					OpCode: meta.OP_BINARY_GT,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] > %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},

				"Ge": {
					OpCode: meta.OP_BINARY_GE,
					ApplyTo: []OperationApplyTo{
						{
							SeriesName: "Durations",
							SeriesType: meta.DurationType,
							MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
								return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s] >= %s.Data_[%s]", res, resIndex, op1, op1Index, op2, op2Index)}
							},
						},
					},
				},
			},
		},
	}

	// Generate entries for NAs
	opNames := []string{"Mul", "Div", "Mod", "Exp", "Add", "Sub", "Eq", "Ne", "Lt", "Le", "Gt", "Ge", "And", "Or"}
	opCodes := []meta.OPCODE{
		meta.OP_BINARY_MUL, meta.OP_BINARY_DIV, meta.OP_BINARY_MOD, meta.OP_BINARY_EXP, meta.OP_BINARY_ADD, meta.OP_BINARY_SUB,
		meta.OP_BINARY_EQ, meta.OP_BINARY_NE, meta.OP_BINARY_LT, meta.OP_BINARY_LE, meta.OP_BINARY_GT, meta.OP_BINARY_GE,
		meta.OP_BINARY_AND, meta.OP_BINARY_OR,
	}

	seriesNames := []string{"NAs", "Bools", "Ints", "Int64s", "Float64s", "Strings", "Times", "Durations"}
	seriesTypes := []meta.BaseType{meta.NullType, meta.BoolType, meta.IntType, meta.Int64Type, meta.Float64Type, meta.StringType, meta.TimeType, meta.DurationType}

	for i, opName := range opNames {
		applyTo := []OperationApplyTo{}
		for j, seriesName := range seriesNames {
			resType := ComputeResInnerType(opCodes[i], seriesTypes[j], seriesTypes[j])

			// Special case for string concatenation
			if opCodes[i] == meta.OP_BINARY_ADD && seriesTypes[j] == meta.StringType {
				applyTo = append(applyTo, OperationApplyTo{
					SeriesName: seriesName,
					SeriesType: seriesTypes[j],
					MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
						return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(enchanter.NA_TEXT + *%s.Data_[%s])", res, resIndex, op1, op2, op2Index)}
					},
				})
			} else

			// Special case for logical OR
			if opCodes[i] == meta.OP_BINARY_OR && seriesTypes[j] == meta.BoolType {
				applyTo = append(applyTo, OperationApplyTo{
					SeriesName: seriesName,
					SeriesType: seriesTypes[j],
					MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
						return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s]", res, resIndex, op2, op2Index)}
					},
				})
			} else if resType != meta.ErrorType {
				applyTo = append(applyTo, OperationApplyTo{
					SeriesName: seriesName,
					SeriesType: seriesTypes[j],
					MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
						return &ast.Ident{Name: ""}
					},
				})
			}
		}

		data["na_ops.go"].Operations[opName] = Operation{OpCode: opCodes[i], ApplyTo: applyTo}
	}

	// Append NAs to all other series
	fileNames := []string{"bool_ops.go", "int_ops.go", "int64_ops.go", "float64_ops.go", "string_ops.go", "time_ops.go", "duration_ops.go"}
	seriesTypes = []meta.BaseType{meta.BoolType, meta.IntType, meta.Int64Type, meta.Float64Type, meta.StringType, meta.TimeType, meta.DurationType}

	for i, fileName := range fileNames {
		for j, opName := range opNames {
			resType := ComputeResInnerType(opCodes[j], seriesTypes[i], meta.NullType)

			// Special case for string concatenation
			if opCodes[j] == meta.OP_BINARY_ADD && seriesTypes[i] == meta.StringType {
				data[fileName].Operations[opName] = Operation{
					OpCode: data[fileName].Operations[opName].OpCode,
					ApplyTo: append(data[fileName].Operations[opName].ApplyTo, OperationApplyTo{
						SeriesName: "NAs",
						SeriesType: meta.NullType,
						MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
							return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Ctx_.StringPool.Put(*%s.Data_[%s] + enchanter.NA_TEXT)", res, resIndex, op1, op1, op1Index)}
						},
					}),
				}
			} else

			// Special case for logical OR
			if opCodes[j] == meta.OP_BINARY_OR && seriesTypes[i] == meta.BoolType {
				data[fileName].Operations[opName] = Operation{
					OpCode: data[fileName].Operations[opName].OpCode,
					ApplyTo: append(data[fileName].Operations[opName].ApplyTo, OperationApplyTo{
						SeriesName: "NAs",
						SeriesType: meta.NullType,
						MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
							return &ast.Ident{Name: fmt.Sprintf("%s[%s] = %s.Data_[%s]", res, resIndex, op1, op1Index)}
						},
					}),
				}
			} else if resType != meta.ErrorType {
				data[fileName].Operations[opName] = Operation{
					OpCode: data[fileName].Operations[opName].OpCode,
					ApplyTo: append(data[fileName].Operations[opName].ApplyTo, OperationApplyTo{
						SeriesName: "NAs",
						SeriesType: meta.NullType,
						MakeOperation: func(res, resIndex, op1, op1Index, op2, op2Index string) ast.Expr {
							return &ast.Ident{Name: ""}
						},
					}),
				}
			}
		}
	}

	return data
}
