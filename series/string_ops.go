package series

import (
	"fmt"

	"github.com/caerbannogwhite/gandalff"
	"github.com/caerbannogwhite/gandalff/utils"
)

func (s Strings) And(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot AND %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Or(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot OR %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Mul(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Div(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Mod(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Exp(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Add(other any) Series {
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
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + boolToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + boolToString(o.Data_[i]))
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case Ints:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[0])))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[0])))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[0])))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[0])))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[i])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[i])))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[i])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(int64(o.Data_[i])))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[0])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[0])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[0])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[0])))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[i])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[i])))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[i])))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(int64(o.Data_[i])))
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case Int64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + intToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + intToString(o.Data_[i]))
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case Float64s:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[0]))
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + floatToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[0]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[0]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[i]))
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[i]))
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + floatToString(o.Data_[i]))
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
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
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[0])
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[i])
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + *o.Data_[i])
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[0])
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[i])
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + *o.Data_[i])
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case Times:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[0].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[0].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
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
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + o.Data_[i].String())
						}
						return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case NAs:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				result := make([]*string, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				result[0] = s.Ctx_.StringPool.Put(*s.Data_[0] + gandalff.NA_TEXT)
				return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			} else {
				resultSize := o.Len()
				result := make([]*string, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Ctx_.StringPool.Put(*s.Data_[0] + gandalff.NA_TEXT)
				}
				return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				result := make([]*string, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, false)
					copy(resultNullMask, s.NullMask_)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + gandalff.NA_TEXT)
				}
				return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				result := make([]*string, resultSize)
				var resultNullMask []uint8
				if s.IsNullable_ {
					resultNullMask = utils.BinVecInit(resultSize, false)
					copy(resultNullMask, s.NullMask_)
				} else {
					resultNullMask = make([]uint8, 0)
				}
				for i := 0; i < resultSize; i++ {
					result[i] = s.Ctx_.StringPool.Put(*s.Data_[i] + gandalff.NA_TEXT)
				}
				return Strings{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Sub(other any) Series {
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
	default:
		return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Eq(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] == *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] == *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] == *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] == *o.Data_[0]
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
							result[i] = *s.Data_[0] == *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] == *o.Data_[i]
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
							result[i] = *s.Data_[0] == *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] == *o.Data_[i]
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
							result[i] = *s.Data_[i] == *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] == *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] == *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] == *o.Data_[0]
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
							result[i] = *s.Data_[i] == *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] == *o.Data_[i]
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
							result[i] = *s.Data_[i] == *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] == *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
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
			return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Ne(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] != *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] != *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] != *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] != *o.Data_[0]
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
							result[i] = *s.Data_[0] != *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] != *o.Data_[i]
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
							result[i] = *s.Data_[0] != *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] != *o.Data_[i]
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
							result[i] = *s.Data_[i] != *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] != *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] != *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] != *o.Data_[0]
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
							result[i] = *s.Data_[i] != *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] != *o.Data_[i]
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
							result[i] = *s.Data_[i] != *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] != *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
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
			return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Gt(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] > *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] > *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] > *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] > *o.Data_[0]
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
							result[i] = *s.Data_[0] > *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] > *o.Data_[i]
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
							result[i] = *s.Data_[0] > *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] > *o.Data_[i]
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
							result[i] = *s.Data_[i] > *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] > *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] > *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] > *o.Data_[0]
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
							result[i] = *s.Data_[i] > *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] > *o.Data_[i]
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
							result[i] = *s.Data_[i] > *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] > *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Ge(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] >= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] >= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] >= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] >= *o.Data_[0]
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
							result[i] = *s.Data_[0] >= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] >= *o.Data_[i]
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
							result[i] = *s.Data_[0] >= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] >= *o.Data_[i]
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
							result[i] = *s.Data_[i] >= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] >= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] >= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] >= *o.Data_[0]
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
							result[i] = *s.Data_[i] >= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] >= *o.Data_[i]
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
							result[i] = *s.Data_[i] >= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] >= *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Lt(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] < *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] < *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] < *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] < *o.Data_[0]
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
							result[i] = *s.Data_[0] < *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] < *o.Data_[i]
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
							result[i] = *s.Data_[0] < *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] < *o.Data_[i]
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
							result[i] = *s.Data_[i] < *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] < *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] < *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] < *o.Data_[0]
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
							result[i] = *s.Data_[i] < *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] < *o.Data_[i]
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
							result[i] = *s.Data_[i] < *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] < *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Strings) Le(other any) Series {
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
	case Strings:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = *s.Data_[0] <= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = *s.Data_[0] <= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = *s.Data_[0] <= *o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = *s.Data_[0] <= *o.Data_[0]
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
							result[i] = *s.Data_[0] <= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] <= *o.Data_[i]
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
							result[i] = *s.Data_[0] <= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[0] <= *o.Data_[i]
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
							result[i] = *s.Data_[i] <= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] <= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] <= *o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] <= *o.Data_[0]
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
							result[i] = *s.Data_[i] <= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] <= *o.Data_[i]
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
							result[i] = *s.Data_[i] <= *o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = *s.Data_[i] <= *o.Data_[i]
						}
						return Bools{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}
