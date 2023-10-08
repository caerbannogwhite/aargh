package gandalff

import (
	"fmt"
	"math"
)

func (s SeriesInt) Neg() Series {
	for i, v := range s.data {
		s.data[i] = -v
	}
	return s
}

func (s SeriesInt) Mul(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						if o.data[0] {
							result[0] = s.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						if o.data[0] {
							result[0] = s.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						if o.data[0] {
							result[0] = s.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						if o.data[0] {
							result[0] = s.data[0]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[0]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[0]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[0]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[0]
							}
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[0] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[0] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							if o.data[0] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if o.data[0] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							if o.data[i] {
								result[i] = s.data[i]
							}
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] * o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] * o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] * o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] * o.data[0]
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] * o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[0]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] * o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) * o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) * o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) * o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) * o.data[0]
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) * o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[0]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) * o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) * o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) * o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) * o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) * o.data[0]
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) * o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[0]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) * o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesNA:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			} else {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			}
			return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot multiply %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Div(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = float64(s.data[0]) / b2
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = float64(s.data[0]) / b2
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = float64(s.data[0]) / b2
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = float64(s.data[0]) / b2
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[0]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[0]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[0]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[0]) / b2
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = float64(s.data[i]) / b2
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) / float64(o.data[0])
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[0])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / float64(o.data[i])
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) / o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) / o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) / o.data[0]
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) / o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[0]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) / o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot divide %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Mod(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = math.Mod(float64(s.data[0]), b2)
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = math.Mod(float64(s.data[0]), b2)
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = math.Mod(float64(s.data[0]), b2)
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = math.Mod(float64(s.data[0]), b2)
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[0]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[0]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[0]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[0]), b2)
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = math.Mod(float64(s.data[i]), b2)
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = math.Mod(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Mod(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot use modulo %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Exp(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(float64(s.data[0]), b2))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(float64(s.data[0]), b2))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(float64(s.data[0]), b2))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := float64(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = int64(math.Pow(float64(s.data[0]), b2))
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[0]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[0]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[0]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[0]), b2))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := float64(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = int64(math.Pow(float64(s.data[i]), b2))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(math.Pow(float64(s.data[0]), float64(o.data[0])))
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[0]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[0])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(math.Pow(float64(s.data[i]), float64(o.data[i])))
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = math.Pow(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = math.Pow(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = math.Pow(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = math.Pow(float64(s.data[0]), float64(o.data[0]))
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[0]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[0]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = math.Pow(float64(s.data[i]), float64(o.data[i]))
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot use exponentiation %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Add(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] + b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] + b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] + b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] + b2
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] + b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] + b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] + o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] + o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] + o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] + o.data[0]
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] + o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[0]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] + o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) + o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) + o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) + o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) + o.data[0]
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) + o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[0]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) + o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) + o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) + o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) + o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) + o.data[0]
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) + o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[0]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) + o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesString:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[0])
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[0])
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[0])
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[0])
						return SeriesString{isNullable: false, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := o.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[0])) + *o.data[i])
						}
						return SeriesString{isNullable: false, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[0])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[0])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[0])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[0])
						}
						return SeriesString{isNullable: false, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[i])
						}
						return SeriesString{isNullable: true, nullMask: resultNullMask, pool: o.pool, data: result}
					} else {
						resultSize := s.Len()
						result := make([]*string, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = o.pool.Put(intToString(int64(s.data[i])) + *o.data[i])
						}
						return SeriesString{isNullable: false, nullMask: resultNullMask, pool: o.pool, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesNA:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			} else {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			}
			return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot sum %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Sub(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] - b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] - b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] - b2
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] - b2
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] - b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] - b2
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] - o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] - o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] - o.data[0]
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] - o.data[0]
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] - o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[0]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[0]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[i]
						}
						return SeriesInt{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] - o.data[i]
						}
						return SeriesInt{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) - o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) - o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) - o.data[0]
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) - o.data[0]
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) - o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[0]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[0]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[i]
						}
						return SeriesInt64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]int64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) - o.data[i]
						}
						return SeriesInt64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) - o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) - o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) - o.data[0]
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) - o.data[0]
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) - o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[0]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[0]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[i]
						}
						return SeriesFloat64{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]float64, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) - o.data[i]
						}
						return SeriesFloat64{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot subtract %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Eq(other Series) Series {
	switch o := other.(type) {
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] == o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) == o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) == o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesNA:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			} else {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for equality %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Ne(other Series) Series {
	switch o := other.(type) {
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] != o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) != o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) != o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesNA:
		if s.Len() == 1 {
			if o.Len() == 1 {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			} else {
				resultSize := o.Len()
				return SeriesNA{size: resultSize}
			}
		} else {
			if o.Len() == 1 {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			} else if s.Len() == o.Len() {
				resultSize := s.Len()
				return SeriesNA{size: resultSize}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for inequality %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Gt(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] > b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] > b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] > b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] > b2
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] > b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] > b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] > o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) > o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) > o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for greater than %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Ge(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] >= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] >= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] >= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] >= b2
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] >= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] >= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] >= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) >= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) >= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for greater than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Lt(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] < b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] < b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] < b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] < b2
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] < b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] < b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] < o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) < o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) < o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for less than %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}

func (s SeriesInt) Le(other Series) Series {
	switch o := other.(type) {
	case SeriesBool:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] <= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] <= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] <= b2
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						b2 := int(0)
						if o.data[0] {
							b2 = 1
						}
						result[0] = s.data[0] <= b2
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[0] <= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[0] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							b2 := int(0)
							if o.data[i] {
								b2 = 1
							}
							result[i] = s.data[i] <= b2
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = s.data[0] <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = s.data[0] <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = s.data[0] <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = s.data[0] <= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[0] <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = s.data[i] <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesInt64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = int64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = int64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = int64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = int64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = int64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	case SeriesFloat64:
		if s.Len() == 1 {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSS(s.nullMask, o.nullMask, resultNullMask)
						result[0] = float64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						result[0] = float64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						result[0] = float64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						result[0] = float64(s.data[0]) <= o.data[0]
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else {
				if s.isNullable {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrSV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, s.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := o.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[0]) <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
		} else {
			if o.Len() == 1 {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVS(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, o.nullMask[0] == 1)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[0]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			} else if s.Len() == o.Len() {
				if s.isNullable {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						__binVecOrVV(s.nullMask, o.nullMask, resultNullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, s.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					}
				} else {
					if o.isNullable {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(resultSize, false)
						copy(resultNullMask, o.nullMask)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: true, nullMask: resultNullMask, data: result}
					} else {
						resultSize := s.Len()
						result := make([]bool, resultSize)
						resultNullMask := __binVecInit(0, false)
						for i := 0; i < resultSize; i++ {
							result[i] = float64(s.data[i]) <= o.data[i]
						}
						return SeriesBool{isNullable: false, nullMask: resultNullMask, data: result}
					}
				}
			}
			return SeriesError{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
		}
	default:
		return SeriesError{fmt.Sprintf("Cannot compare for less than or equal to %s and %s", s.Type().ToString(), o.Type().ToString())}
	}

}
