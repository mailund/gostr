package gostr_test

import (
	"errors"
	"testing"

	"github.com/mailund/gostr"
)

func TestAlphabetLookupError(t *testing.T) {
	x := "foobar"
	alpha := gostr.NewAlphabet(x)

	if _, err := alpha.MapToBytes("qux"); err == nil {
		t.Fatal("Expected an error here")
	} else if _, ok := err.(*gostr.AlphabetLookupError); !ok {
		t.Errorf("Unexpected error type: %q", err)
	} else if err.Error() != "byte q is not in alphabet" {
		t.Errorf("Unexpected error message: %s", err)
	}
}

func TestInvalidCigar_Is(t *testing.T) {
	cigarErr := gostr.NewInvalidCigar("foo")
	if cigarErr.Error() != "invalid cigar: foo" {
		t.Errorf("Unexpected error message: %s", cigarErr)
	}

	otherCigarErr := gostr.NewInvalidCigar("foo")
	otherDifferentCigarErr := gostr.NewInvalidCigar("bar")

	if !errors.Is(cigarErr, otherCigarErr) {
		t.Error("these errors should be considered the same")
	}

	if errors.Is(cigarErr, otherDifferentCigarErr) {
		t.Error("these errors should be considered different")
	}

	otherErr := errors.New("some other error") //nolint:goerr113 // ignore new error for testing
	if errors.Is(cigarErr, otherErr) {
		t.Error("these errors should be considered different")
	}
}
