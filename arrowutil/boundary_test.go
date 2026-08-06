package arrowutil

import (
	"fmt"
	"testing"
)

// Boundary tests for the null-mask <-> validity-bitmap conversion at byte
// boundaries. enchanter uses bit-set=null, Arrow uses bit-set=valid; trailing
// bits past the logical length must always be zero.

func bit(b []byte, i int) bool { return b[i>>3]&(1<<uint(i%8)) != 0 }

var blLengths = []int{1, 7, 8, 9, 15, 16, 63, 64, 65}

var blPatterns = []struct {
	name string
	fn   func(i, n int) bool // true = null at index i
}{
	{"none", func(i, n int) bool { return false }},
	{"all", func(i, n int) bool { return true }},
	{"alternating", func(i, n int) bool { return i%2 == 0 }},
	{"firstonly", func(i, n int) bool { return i == 0 }},
	{"lastonly", func(i, n int) bool { return i == n-1 }},
}

func TestBitmapConversionBoundaries(t *testing.T) {
	for _, n := range blLengths {
		for _, p := range blPatterns {
			t.Run(fmt.Sprintf("n=%d_%s", n, p.name), func(t *testing.T) {
				nbytes := (n + 7) / 8
				nullMask := make([]uint8, nbytes)
				boolMask := make([]bool, n)
				for i := 0; i < n; i++ {
					if p.fn(i, n) {
						nullMask[i>>3] |= 1 << uint(i%8)
						boolMask[i] = true
					}
				}

				validity := NullMaskToValidityBitmap(nullMask, n)
				if len(validity) != nbytes {
					t.Fatalf("validity length: got %d, want %d", len(validity), nbytes)
				}
				for i := 0; i < n; i++ {
					if bit(validity, i) != !p.fn(i, n) {
						t.Fatalf("validity bit %d: got %v, want %v", i, bit(validity, i), !p.fn(i, n))
					}
				}
				for i := n; i < nbytes*8; i++ {
					if bit(validity, i) {
						t.Fatalf("validity trailing bit %d (length %d) must be clear", i, n)
					}
				}

				back := ValidityBitmapToNullMask(validity, n)
				for i := 0; i < nbytes; i++ {
					if back[i] != nullMask[i] {
						t.Fatalf("null mask round trip byte %d: got %08b, want %08b", i, back[i], nullMask[i])
					}
				}

				fromBools := BoolSliceToValidityBitmap(boolMask, n)
				for i := 0; i < nbytes; i++ {
					if fromBools[i] != validity[i] {
						t.Fatalf("BoolSliceToValidityBitmap byte %d: got %08b, want %08b", i, fromBools[i], validity[i])
					}
				}
			})
		}
	}
}

func TestAllValidBoundaries(t *testing.T) {
	for _, n := range blLengths {
		buf := AllValid(n)
		for i := 0; i < n; i++ {
			if !bit(buf, i) {
				t.Fatalf("n=%d: bit %d must be set", n, i)
			}
		}
		for i := n; i < len(buf)*8; i++ {
			if bit(buf, i) {
				t.Fatalf("n=%d: trailing bit %d must be clear", n, i)
			}
		}
	}
}
