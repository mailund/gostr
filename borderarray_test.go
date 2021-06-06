package gostr

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr/test"
)

func Test_BorderarrayBasics(t *testing.T) {
	tests := []struct {
		name string
		x    string
		want []int
	}{
		{"(empty string)", "", []int{}},
		{"a", "a", []int{0}},
		{"aaa", "aaa", []int{0, 1, 2}},
		{"aaaba", "aaaba", []int{0, 1, 2, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Borderarray(tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Borderarray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StrictBorderarrayBasics(t *testing.T) {
	tests := []struct {
		name string
		x    string
		want []int
	}{
		{"(empty string)", "", []int{}},
		{"a", "a", []int{0}},
		{"aaa", "aaa", []int{0, 0, 2}},
		{"aaaba", "aaaba", []int{0, 0, 2, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrictBorderarray(tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrictBorderarray() = %v, want %v", got, tt.want)
			}
		})
	}
}

// FIXME: check if the border is the *longest* as well
func checkBorders(t *testing.T, x string, ba []int) bool {
	for i, b := range ba {
		if b != 0 && x[:b] != x[i-b+1:i+1] {
			t.Errorf(`x[:%d] == %q is not a border of %q`, b, x[:b], x[:i+1])
			t.Fatalf(`x = %q, ba = %v`, x, ba)
			return false
		}
	}
	return true
}

func Test_Borderarray(t *testing.T) {
	rng := test.NewRandomSeed(t)
	test.GenerateTestStrings(10, 20, rng,
		func(x string) {
			checkBorders(t, x, Borderarray(x))
		})
}

func checkStrict(t *testing.T, x string, ba []int) bool {
	for i, b := range ba[:len(ba)-1] {
		if b > 0 && x[b] == x[i+1] {
			t.Errorf(`x[:%d] == %q[%q] is not a strict border of %q[%q]`, b, x[:b], x[b], x[:i+1], x[i+1])
			t.Errorf(`x[%d] == %q == x[%d+1] (should be different)`, b, x[b], i)
			t.Fatalf(`x = %q, ba = %v`, x, ba)
			return false
		}
	}
	return true
}

func Test_StrictBorderarray(t *testing.T) {
	rng := test.NewRandomSeed(t)
	test.GenerateTestStrings(10, 20, rng,
		func(x string) {
			ba := StrictBorderarray(x)
			checkBorders(t, x, ba)
			checkStrict(t, x, ba)
		})
}
