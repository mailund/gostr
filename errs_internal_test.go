package gostr

import (
	"testing"
)

func checkThenError() (err error) {
	defer catchError(&err)

	// This should catch the error and make it the return value for
	// the function.
	checkError(NewInvalidCigar("foo"))

	return nil
}

func TestCheckCatch(t *testing.T) {
	if err := checkThenError(); err == nil {
		t.Fatal("We expected an error")
	} else if err.Error() != "invalid cigar: foo" {
		t.Errorf("Unexpected error message: %s", err)
	}
}
