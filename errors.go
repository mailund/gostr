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
