package gostr_test

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

func TestOpsToCigar(t *testing.T) {
	type args struct {
		ops gostr.EditOps
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Single M",
			args{ops: []gostr.ApproxEdit{gostr.M}},
			"1M",
		},
		{
			"Single D",
			args{ops: []gostr.ApproxEdit{gostr.D}},
			"1D",
		},
		{
			"Single I",
			args{ops: []gostr.ApproxEdit{gostr.I}},
			"1I",
		},
		{
			"IIMMMDDI",
			args{ops: []gostr.ApproxEdit{
				gostr.I, gostr.I, gostr.M, gostr.M, gostr.M, gostr.D, gostr.D, gostr.I}},
			"2I3M2D1I",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gostr.OpsToCigar(tt.args.ops); got != tt.want {
				t.Errorf("OpsToCigar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCigarToOps(t *testing.T) {
	tests := []struct {
		cigar string
		want  gostr.EditOps
	}{
		{
			"1M",
			gostr.EditOps{gostr.M},
		},
		{
			"10M",
			gostr.EditOps{gostr.M, gostr.M, gostr.M, gostr.M, gostr.M, gostr.M, gostr.M, gostr.M, gostr.M, gostr.M},
		},
		{
			"1I",
			gostr.EditOps{gostr.I},
		},
		{
			"1D",
			gostr.EditOps{gostr.D},
		},
		{
			"1D2M3I",
			gostr.EditOps{gostr.D, gostr.M, gostr.M, gostr.I, gostr.I, gostr.I},
		},
	}

	for _, tt := range tests {
		t.Run(tt.cigar, func(t *testing.T) {
			if got := gostr.CigarToOps(tt.cigar); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CigarToOps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractAlignment(t *testing.T) {
	type args struct {
		x     string
		p     string
		pos   int
		cigar string
	}

	tests := []struct {
		name     string
		args     args
		wantSubx string
		wantSubp string
	}{
		{
			"Just matches",
			args{"acgtacgt", "gtac", 2, "4M"},
			"gtac", "gtac",
		},
		{
			"Deletion",
			args{"acgtacgt", "gtc", 2, "2M1D1M"},
			"gtac", "gt-c",
		},
		{
			"Insertion",
			args{"acgtacgt", "gtaac", 2, "2M1I2M"},
			"gt-ac", "gtaac",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubx, gotSubp := gostr.ExtractAlignment(tt.args.x, tt.args.p, tt.args.pos, tt.args.cigar)
			if gotSubx != tt.wantSubx {
				t.Errorf("ExtractAlignment() gotSubx = %v, want %v", gotSubx, tt.wantSubx)
			}
			if gotSubp != tt.wantSubp {
				t.Errorf("ExtractAlignment() gotSubp = %v, want %v", gotSubp, tt.wantSubp)
			}
		})
	}
}

func TestCountEdits(t *testing.T) {
	type args struct {
		x     string
		p     string
		pos   int
		cigar string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Just matches",
			args{"acgtacgt", "gtac", 2, "4M"},
			0, // "gtac" vs "gtac",
		},
		{
			"Deletion",
			args{"acgtacgt", "gtc", 2, "2M1D1M"},
			1, // "gtac", "gt-c",
		},
		{
			"Insertion",
			args{"acgtacgt", "gtaac", 2, "2M1I2M"},
			1, // "gt-ac", "gtaac",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gostr.CountEdits(tt.args.x, tt.args.p, tt.args.pos, tt.args.cigar)
			if got != tt.want {
				t.Errorf("CountEdits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type approxAlgo = func(string) func(string, int, func(int, string))

var approxAlgorithms = map[string]approxAlgo{
	"BWA": gostr.FMIndexApproxPreprocess,
}

func runRandomApproxOccurencesTests(algo approxAlgo) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		rng := test.NewRandomSeed(t)
		test.GenerateTestStringsAndPatterns(10, 50, rng,
			func(x, p string) {
				search := algo(x)
				for edits := 1; edits < 3; edits++ {
					search(p, edits, func(pos int, cigar string) {
						count := gostr.CountEdits(x, p, pos, cigar)
						if count > edits {
							t.Errorf("Match at pos %d needs too many edits, %d vs %d",
								pos, count, edits)
						}
					})
				}
			})
	}
}

func TestRandomApproxOccurences(t *testing.T) {
	t.Helper()

	for name, algo := range approxAlgorithms {
		t.Run(name, runRandomApproxOccurencesTests(algo))
	}
}
