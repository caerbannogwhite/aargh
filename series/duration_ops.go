package series

import (
	"fmt"
	"time"

	"github.com/caerbannogwhite/gandalff/utils"
)

func (s Durations) And(other any) Series {
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

func (s Durations) Or(other any) Series {
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

func (s Durations) Mul(other any) Series {
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

func (s Durations) Div(other any) Series {
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

func (s Durations) Mod(other any) Series {
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

func (s Durations) Exp(other any) Series {
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

func (s Durations) Add(other any) Series {
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
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[0])
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[0])
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
							result[i] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[i])
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
							result[i] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[0].String() + *o.Data_[i])
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
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[0])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[0])
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
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[i])
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
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[i])
						}
						return Strings{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Ctx_.StringPool.Put(s.Data_[i].String() + *o.Data_[i])
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
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = o.Data_[0].Add(s.Data_[0])
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = o.Data_[0].Add(s.Data_[0])
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = o.Data_[0].Add(s.Data_[0])
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = o.Data_[0].Add(s.Data_[0])
						return Times{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[0])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[0])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[0])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[0])
						}
						return Times{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[0].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[0].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[0].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[0].Add(s.Data_[i])
						}
						return Times{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[i])
						}
						return Times{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Time, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.Data_[i].Add(s.Data_[i])
						}
						return Times{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
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
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] + o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] + o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] + o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] + o.Data_[0]
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] + o.Data_[i]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[0]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] + o.Data_[i]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
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
			return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Durations) Sub(other any) Series {
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
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] - o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] - o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] - o.Data_[0]
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] - o.Data_[0]
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] - o.Data_[i]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVS(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[0]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[0]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			} else if s.Len() == o.Len() {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrVV(s.NullMask_, o.NullMask_, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, o.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[i]
						}
						return Durations{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]time.Duration, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] - o.Data_[i]
						}
						return Durations{IsNullable_: false, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				}
			}
			return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return Errors{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s Durations) Eq(other any) Series {
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
	case Durations:
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

func (s Durations) Ne(other any) Series {
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
	case Durations:
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

func (s Durations) Gt(other any) Series {
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
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] > o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] > o.Data_[0]
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
							result[i] = s.Data_[0] > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] > o.Data_[i]
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
							result[i] = s.Data_[0] > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] > o.Data_[i]
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
							result[i] = s.Data_[i] > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] > o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] > o.Data_[0]
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
							result[i] = s.Data_[i] > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] > o.Data_[i]
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
							result[i] = s.Data_[i] > o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] > o.Data_[i]
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

func (s Durations) Ge(other any) Series {
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
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] >= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] >= o.Data_[0]
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
							result[i] = s.Data_[0] >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] >= o.Data_[i]
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
							result[i] = s.Data_[0] >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] >= o.Data_[i]
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
							result[i] = s.Data_[i] >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] >= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] >= o.Data_[0]
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
							result[i] = s.Data_[i] >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] >= o.Data_[i]
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
							result[i] = s.Data_[i] >= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] >= o.Data_[i]
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

func (s Durations) Lt(other any) Series {
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
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] < o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] < o.Data_[0]
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
							result[i] = s.Data_[0] < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] < o.Data_[i]
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
							result[i] = s.Data_[0] < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] < o.Data_[i]
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
							result[i] = s.Data_[i] < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] < o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] < o.Data_[0]
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
							result[i] = s.Data_[i] < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] < o.Data_[i]
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
							result[i] = s.Data_[i] < o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] < o.Data_[i]
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

func (s Durations) Le(other any) Series {
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
	case Durations:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.IsNullable_ {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						utils.BinVecOrSS(s.NullMask_, o.NullMask_, resultNullMask)
						result[0] = s.Data_[0] <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						result[0] = s.Data_[0] <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						result[0] = s.Data_[0] <= o.Data_[0]
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						result[0] = s.Data_[0] <= o.Data_[0]
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
							result[i] = s.Data_[0] <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, s.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] <= o.Data_[i]
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
							result[i] = s.Data_[0] <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[0] <= o.Data_[i]
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
							result[i] = s.Data_[i] <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					}
				} else {
					if o.IsNullable_ {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, o.NullMask_[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] <= o.Data_[0]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] <= o.Data_[0]
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
							result[i] = s.Data_[i] <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(resultSize, false)
						copy(resultNullMask, s.NullMask_)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] <= o.Data_[i]
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
							result[i] = s.Data_[i] <= o.Data_[i]
						}
						return Bools{IsNullable_: true, NullMask_: resultNullMask, Data_: result, Ctx_: s.Ctx_}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := utils.BinVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.Data_[i] <= o.Data_[i]
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
