package arrowutil

import (
	"testing"
)

func TestNullMaskToValidityBitmapRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		nullMask []uint8
		length   int
	}{
		{"empty", []uint8{}, 0},
		{"no nulls 8", []uint8{0x00}, 8},
		{"all nulls 8", []uint8{0xFF}, 8},
		{"some nulls 8", []uint8{0b10100101}, 8},
		{"partial byte", []uint8{0b00000101}, 5},
		{"two bytes", []uint8{0xFF, 0x0F}, 12},
		{"16 elements", []uint8{0xAA, 0x55}, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validity := NullMaskToValidityBitmap(tt.nullMask, tt.length)
			got := ValidityBitmapToNullMask(validity, tt.length)

			// Compare bit-by-bit up to length
			for i := 0; i < tt.length; i++ {
				byteIdx := i / 8
				bitIdx := uint(i % 8)
				var origNull bool
				if byteIdx < len(tt.nullMask) {
					origNull = tt.nullMask[byteIdx]&(1<<bitIdx) != 0
				}
				gotNull := got[byteIdx]&(1<<bitIdx) != 0
				if origNull != gotNull {
					t.Errorf("bit %d: origNull=%v, gotNull=%v", i, origNull, gotNull)
				}
			}
		})
	}
}

func TestAllValid(t *testing.T) {
	tests := []struct {
		length int
	}{
		{0},
		{1},
		{7},
		{8},
		{9},
		{16},
		{17},
	}

	for _, tt := range tests {
		buf := AllValid(tt.length)
		for i := 0; i < tt.length; i++ {
			byteIdx := i / 8
			bitIdx := uint(i % 8)
			if buf[byteIdx]&(1<<bitIdx) == 0 {
				t.Errorf("AllValid(%d): bit %d not set", tt.length, i)
			}
		}
		// Trailing bits should be clear
		nbytes := (tt.length + 7) / 8
		if nbytes > 0 {
			trail := tt.length % 8
			if trail != 0 {
				mask := byte(0xFF) << uint(trail)
				if buf[nbytes-1]&mask != 0 {
					t.Errorf("AllValid(%d): trailing bits set: %08b", tt.length, buf[nbytes-1])
				}
			}
		}
	}
}

func TestBoolSliceToValidityBitmap(t *testing.T) {
	// nullMask: true means null
	nullMask := []bool{false, true, false, true, false}
	length := 5

	validity := BoolSliceToValidityBitmap(nullMask, length)

	// bit 0 should be valid (set)
	if validity[0]&(1<<0) == 0 {
		t.Error("bit 0 should be valid")
	}
	// bit 1 should be null (clear)
	if validity[0]&(1<<1) != 0 {
		t.Error("bit 1 should be null")
	}
	// bit 2 should be valid (set)
	if validity[0]&(1<<2) == 0 {
		t.Error("bit 2 should be valid")
	}
	// bit 3 should be null (clear)
	if validity[0]&(1<<3) != 0 {
		t.Error("bit 3 should be null")
	}
	// bit 4 should be valid (set)
	if validity[0]&(1<<4) == 0 {
		t.Error("bit 4 should be valid")
	}

	// Round-trip back to null mask
	recovered := ValidityBitmapToNullMask(validity, length)
	for i := 0; i < length; i++ {
		isNull := recovered[i/8]&(1<<uint(i%8)) != 0
		if isNull != nullMask[i] {
			t.Errorf("bit %d: expected null=%v, got null=%v", i, nullMask[i], isNull)
		}
	}
}

func TestNullMaskToValidityInversion(t *testing.T) {
	// Verify that each bit is properly inverted
	nullMask := []uint8{0b10101010} // bits 1,3,5,7 are null
	length := 8

	validity := NullMaskToValidityBitmap(nullMask, length)

	// Expected: 0b01010101 (bits 0,2,4,6 are valid)
	if validity[0] != 0b01010101 {
		t.Errorf("expected 0b01010101, got %08b", validity[0])
	}
}
