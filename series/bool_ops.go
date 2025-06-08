package series

import (
	"fmt"
	"math"

	"github.com/caerbannogwhite/aargh/utils"
)

// Not performs logical NOT operation on series
func (s Bools) Not() Series {
	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = !s.Data_[i]
	}

	return s
}

func (s Bools) All() bool {
	if s.IsNullable_ {
		for i := 0; i < len(s.Data_); i++ {
			if s.NullMask_[i>>3]&(1<<uint(i%8)) == 0 && !s.Data_[i] {
				return false
			}
		}

		return true
	} else {
		for i := 0; i < len(s.Data_); i++ {
			if !s.Data_[i] {
				return false
			}
		}

		return true
	}
}

func (s Bools) Any() bool {
	if s.IsNullable_ {
		for i := 0; i < len(s.Data_); i++ {
			if s.NullMask_[i>>3]&(1<<uint(i%8)) == 0 && s.Data_[i] {
				return true
			}
		}

		return false
	} else {
		for i := 0; i < len(s.Data_); i++ {
			if s.Data_[i] {
				return true
			}
		}

		return false
	}
}

func (s Bools) And(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] && o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] && o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] && o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] && o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] && o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] && o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot AND %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return NAs{size: resultSize}
			} else {
				resultSize := o.Len()
				return NAs{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return NAs{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return NAs{size: resultSize}
			}
			return Errors{fmt.Sprintf("Cannot AND %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot AND %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Or(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] || o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] || o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] || o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] || o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] || o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] || o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot OR %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				result := make([]bool, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				result[0] = s.Data_[0]
				return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			} else {
				resultSize := o.Len()
				result := make([]bool, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Data_[0]
				}
				return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				result := make([]bool, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, false)
					copy(resultNullMask, s.NullMask_)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Data_[i]
				}
				return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				result := make([]bool, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, false)
					copy(resultNullMask, s.NullMask_)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Data_[i]
				}
				return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			}
			return Errors{fmt.Sprintf("Cannot OR %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot OR %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Mul(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						if s.Data_[0] && o.Data_[0] {
							result[0] = 1
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						if s.Data_[0] && o.Data_[0] {
							result[0] = 1
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						if s.Data_[0] && o.Data_[0] {
							result[0] = 1
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						if s.Data_[0] && o.Data_[0] {
							result[0] = 1
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[0] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[0] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[0] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[0] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] && o.Data_[i] {
								result[i] = 1
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						if s.Data_[0] {
							result[0] = o.Data_[0]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[0] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[0]
							}
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if s.Data_[i] {
								result[i] = o.Data_[i]
							}
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return NAs{size: resultSize}
			} else {
				resultSize := o.Len()
				return NAs{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return NAs{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return NAs{size: resultSize}
			}
			return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Div(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 / b2
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 / b2
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 / b2
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 / b2
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 / b2
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / float64(o.Data_[0])
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[0])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / float64(o.Data_[i])
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 / o.Data_[0]
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[0]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 / o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Mod(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = math.Mod(b1, b2)
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = math.Mod(b1, b2)
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = math.Mod(b1, b2)
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = math.Mod(b1, b2)
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = math.Mod(b1, b2)
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = math.Mod(b1, float64(o.Data_[0]))
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[0]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = math.Mod(b1, float64(o.Data_[i]))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Exp(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(b1, b2))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(b1, b2))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(b1, b2))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						b2 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(b1, b2))
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							b2 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(b1, b2))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = int64(math.Pow(b1, float64(o.Data_[0])))
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = int64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = float64(math.Pow(b1, float64(o.Data_[0])))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = float64(math.Pow(b1, float64(o.Data_[0])))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = float64(math.Pow(b1, float64(o.Data_[0])))
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = float64(math.Pow(b1, float64(o.Data_[0])))
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[0])))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = float64(math.Pow(b1, float64(o.Data_[i])))
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Add(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 + b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 + b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 + b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 + b2
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 + b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 + o.Data_[0]
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[0]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 + o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[0])
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[0]) + *o.Data_[i])
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[0])
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(boolToString(s.Data_[i]) + *o.Data_[i])
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return NAs{size: resultSize}
			} else {
				resultSize := o.Len()
				return NAs{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return NAs{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return NAs{size: resultSize}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Sub(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 - b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 - b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 - b2
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 - b2
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 - b2
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Ints{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Int64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 - o.Data_[0]
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[0]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 - o.Data_[i]
						}
						return Float64s{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Eq(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] == o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] == o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] == o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] == o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] == o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] == o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return NAs{size: resultSize}
			} else {
				resultSize := o.Len()
				return NAs{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return NAs{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return NAs{size: resultSize}
			}
			return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Ne(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] != o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] != o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] != o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] != o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] != o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] != o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().String(), o.Type().String())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return NAs{size: resultSize}
			} else {
				resultSize := o.Len()
				return NAs{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return NAs{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return NAs{size: resultSize}
			}
			return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Gt(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 > b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 > b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 > b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 > b2
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 > b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 > o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 > o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Ge(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 >= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 >= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 >= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 >= b2
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 >= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 >= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 >= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Lt(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 < b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 < b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 < b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 < b2
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 < b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 < o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 < o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().String(), o.Type().String())}
	}

}

func (s Bools) Le(other any) Series {
	var otherSeries Series
	if _, ok := other.(Series); ok {
		otherSeries = other.(Series)
	} else {
		otherSeries = NewSeries(other, nil, false, false, s.Ctx_)
	}
	if s.Ctx_ != otherSeries.GetContext() {
		return Errors{fmt.Sprintf("Cannot operate on series with different contexts: %v and %v", s.Ctx_, otherSeries.GetContext())}
	}
	switch o := otherSeries.(type) {
	case Bools:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 <= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 <= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 <= b2
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						b2 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						if o.Data_[0] {
							b2 = 1
						}
						result[0] = b1 <= b2
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[0] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							b2 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							if o.Data_[i] {
								b2 = 1
							}
							result[i] = b1 <= b2
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := int64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := int64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						b1 := float64(0)
						if s.Data_[0] {
							b1 = 1
						}
						result[0] = b1 <= o.Data_[0]
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[0] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[0]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b1 := float64(0)
							if s.Data_[i] {
								b1 = 1
							}
							result[i] = b1 <= o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().String(), o.Type().String())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().String(), o.Type().String())}
	}

}
