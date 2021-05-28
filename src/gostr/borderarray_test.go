package gostr

import (
	"reflect"
	"testing"
)

func Test_Borderarray(t *testing.T) {
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

func Test_StrictBorderarray(t *testing.T) {
	tests := []struct {
		name string
		x string
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
