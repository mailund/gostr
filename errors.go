package gostr

import "fmt"

// AlphabetLookupError are errors that occur if you look up a character
// that is not in the alphabet.
type AlphabetLookupError struct {
	char byte
}

// Error implements the interface for errors.
func (err *AlphabetLookupError) Error() string {
	return fmt.Sprintf("byte %b is not in alphabet", err.char)
}

// InvalidCigar are errors when you use a cigar that isn't in the right format
type InvalidCigar struct {
	x string
}

// NewInvalidCigar creates an InvalidCigar error
func NewInvalidCigar(x string) *InvalidCigar {
	return &InvalidCigar{x: x}
}

// Error implements the interface for errors.
func (err *InvalidCigar) Error() string {
	return fmt.Sprintf("invalid cigar: %s", err.x)
}

// Error implements the Is interface for errors.
func (err *InvalidCigar) Is(other error) bool {
	if ic, ok := other.(*InvalidCigar); ok {
		return ic.x == err.x
	}

	return false
}
