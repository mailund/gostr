package gostr_test // black box testing...

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr/gostr"
	"github.com/mailund/gostr/testutils"
)

type SAAlgo = func(x string) []int32

var saAlgorithms = map[string]SAAlgo{
	"Skew":       gostr.Skew,
	"Sais":       gostr.Sais,
	"SuffixTree": gostr.StSaConstruction,
}

func runBasicTest(algo SAAlgo) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

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
	t.Helper()

	for name, algo := range saAlgorithms {
		t.Run(name, runBasicTest(algo))
	}
}

func runConsistencyTest(algo SAAlgo) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		rng := testutils.NewRandomSeed(t)
		testutils.GenerateTestStrings(50, 150, rng,
			func(x string) {
				testutils.CheckSuffixArray(t, x, algo(x))
			})
	}
}

func Test_SuffixArraysConsistency(t *testing.T) {
	for name, algo := range saAlgorithms {
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
