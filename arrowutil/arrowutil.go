// Package arrowutil provides conversion utilities between enchanter's null
// mask convention (bit-set=null) and Arrow's validity bitmap convention
// (bit-set=valid).
package arrowutil

// NullMaskToValidityBitmap converts an enchanter null mask (bit-set=null) to an
// Arrow validity bitmap (bit-set=valid). The result is a byte slice suitable
// for use with arrow.ArrayData or builders.
func NullMaskToValidityBitmap(nullMask []uint8, length int) []byte {
	nbytes := (length + 7) / 8
	validity := make([]byte, nbytes)
	// Start with all bits set (all valid), then clear bits that are null.
	for i := 0; i < nbytes; i++ {
		if i < len(nullMask) {
			validity[i] = ^nullMask[i]
		} else {
			validity[i] = 0xFF
		}
	}
	// Clear trailing bits beyond length so they don't contribute.
	if trail := length % 8; trail != 0 {
		validity[nbytes-1] &= (1 << uint(trail)) - 1
	}
	return validity
}

// ValidityBitmapToNullMask converts an Arrow validity bitmap (bit-set=valid)
// to an enchanter null mask (bit-set=null).
func ValidityBitmapToNullMask(validBitmap []byte, length int) []uint8 {
	nbytes := (length + 7) / 8
	nullMask := make([]uint8, nbytes)
	for i := 0; i < nbytes; i++ {
		if i < len(validBitmap) {
			nullMask[i] = ^validBitmap[i]
		} else {
			nullMask[i] = 0xFF
		}
	}
	// Clear trailing bits beyond length.
	if trail := length % 8; trail != 0 {
		nullMask[nbytes-1] &= (1 << uint(trail)) - 1
	}
	return nullMask
}

// AllValid returns an Arrow validity bitmap where all bits are set (all valid).
func AllValid(length int) []byte {
	nbytes := (length + 7) / 8
	buf := make([]byte, nbytes)
	for i := range buf {
		buf[i] = 0xFF
	}
	if trail := length % 8; trail != 0 {
		buf[nbytes-1] = (1 << uint(trail)) - 1
	}
	return buf
}

// BoolSliceToValidityBitmap converts a bool slice (true=null, matching enchanter
// convention) to an Arrow validity bitmap (bit-set=valid).
func BoolSliceToValidityBitmap(nullMask []bool, length int) []byte {
	nbytes := (length + 7) / 8
	validity := make([]byte, nbytes)
	for i := range validity {
		validity[i] = 0xFF
	}
	for i := 0; i < length && i < len(nullMask); i++ {
		if nullMask[i] {
			// This index is null → clear the valid bit.
			validity[i/8] &= ^(1 << uint(i%8))
		}
	}
	// Clear trailing bits beyond length.
	if trail := length % 8; trail != 0 {
		validity[nbytes-1] &= (1 << uint(trail)) - 1
	}
	return validity
}
