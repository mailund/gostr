package gostr

import "fmt"

type AlphabetLookupError struct {
	char byte
}

func (err *AlphabetLookupError) Error() string {
	return fmt.Sprintf("byte %b is not in alphabet", err.char)
}
