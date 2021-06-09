package gostr_test // black box testing...

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

type SAAlgoString = func(x *gostr.String) []int
type SAAlgo = func(x string) []int

func wrap(algo SAAlgoString) SAAlgo {
	return func(x_ string) []int {
		x, _ := gostr.NewString(x_, nil)
		return algo(x)
	}
}

var sa_algorithms = map[string]SAAlgo{
	"Skew": wrap(gostr.Skew),
	"Sais": wrap(gostr.Sais),
}

func runBasicTest(algo SAAlgo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name   string
			x      string
			wantSA []int
		}{
			{`We handle empty strings`, "", []int{0}},
			{`Unique characters "a"`, "a", []int{1, 0}},
			{`Unique characters "ab"`, "ab", []int{2, 0, 1}},
			{`Unique characters "ba"`, "ba", []int{2, 1, 0}},
			{`Unique characters "abc"`, "abc", []int{3, 0, 1, 2}},
			{`Unique characters "bca"`, "bca", []int{3, 2, 0, 1}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotSA := algo(tt.x); !reflect.DeepEqual(gotSA, tt.wantSA) {
					t.Errorf("Got = %v, want %v", gotSA, tt.wantSA)
				}
			})
		}
	}
}

func Test_SuffixArraysBasic(t *testing.T) {
	for name, algo := range sa_algorithms {
		t.Run(name, runBasicTest(algo))
	}
}

func runConsistencyTest(algo SAAlgo) func(*testing.T) {
	return func(t *testing.T) {
		rng := test.NewRandomSeed(t)
		test.GenerateTestStrings(50, 150, rng,
			func(x string) {
				test.CheckSuffixArray(t, x, algo(x))
			})
	}
}

func Test_SuffixArraysConsistency(t *testing.T) {
	for name, algo := range sa_algorithms {
		t.Run(name, runConsistencyTest(algo))
	}
}
