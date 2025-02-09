package testutils_test

import (
	"fmt"
	"testing"

	test "github.com/mailund/gostr/testutils"
)

func TestFibonacciString(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "a"},
		{1, "b"},
		{2, "ab"},
		{3, "bab"},
		{4, "abbab"},
		{5, "bababbab"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Fib(%d)", tt.n), func(t *testing.T) {
			if got := test.FibonacciString(tt.n); got != tt.want {
				t.Errorf("FibonacciString() = %v, want %v", got, tt.want)
			}
		})
	}
}
