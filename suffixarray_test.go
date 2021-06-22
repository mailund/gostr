package gostr_test // black box testing...

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

type SAAlgo = func(x string) []int32

var sa_algorithms = map[string]SAAlgo{
	"Skew":       gostr.Skew,
	"Sais":       gostr.Sais,
	"SuffixTree": gostr.StSaConstruction,
}

func runBasicTest(algo SAAlgo) func(*testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name   string
			x      string
			wantSA []int32
		}{
			{`We handle empty strings`, "", []int32{0}},
			{`Unique characters "a"`, "a", []int32{1, 0}},
			{`Unique characters "ab"`, "ab", []int32{2, 0, 1}},
			{`Unique characters "ba"`, "ba", []int32{2, 1, 0}},
			{`Unique characters "abc"`, "abc", []int32{3, 0, 1, 2}},
			{`Unique characters "bca"`, "bca", []int32{3, 2, 0, 1}},
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

func Test_AlphabetErrors(t *testing.T) {
	alpha := gostr.NewAlphabet("foo")
	x := "bar" // wrong alphabet
	if _, err := gostr.SaisWithAlphabet(x, alpha); err == nil {
		t.Error("Expected an error making Sais SA")
	}
	if _, err := gostr.SkewWithAlphabet(x, alpha); err == nil {
		t.Error("Expected an error making Skew SA")
	}
}
