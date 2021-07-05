package test

import (
	"reflect"
	"testing"
)

// IntArraysEqual tests if arrays a and b are equal.
func IntArraysEqual(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}

// IsPrefix tests if string x is a prefix of y
func IsPrefix(x, y string) bool {
	if len(y) < len(x) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}

// OccurrenceAt returns if string p occurrs at index i in string x.
func OccurrenceAt(x, p string, i int) bool {
	return IsPrefix(p, x[i:])
}

// CheckOccurrenceAt tests if string p occurrs at index i in string x and reports
// an error to t otherwise
func CheckOccurrenceAt(t *testing.T, x, p string, i int) bool {
	t.Helper()

	result := OccurrenceAt(x, p, i)
	if !result {
		t.Errorf(`We have an incorrect match: "%s" doesn't match "%s"`,
			p, x[i:])
	}

	return result
}

// CheckAllOccurrences Tests if string p occurrs at all indices
// in occ. Reports an error to t otherwise
func CheckAllOccurrences(t *testing.T, x, p string, occ []int) bool {
	t.Helper()

	result := true

	for _, i := range occ {
		if !OccurrenceAt(x, p, i) {
			t.Errorf(`We have an incorrect match: "%s" doesn't match "%s"`,
				p, x[i:])

			result = false
		}
	}

	return result
}
