package gostr_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mailund/gostr/gostr"
	"github.com/mailund/gostr/testutils"
)

func TestOpsToCigar(t *testing.T) {
	t.Parallel()

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
			args{ops: []gostr.ApproxEdit{gostr.Match}},
			"1M",
		},
		{
			"Single D",
			args{ops: []gostr.ApproxEdit{gostr.Delete}},
			"1D",
		},
		{
			"Single I",
			args{ops: []gostr.ApproxEdit{gostr.Insert}},
			"1I",
		},
		{
			"IIMMMDDI",
			args{ops: []gostr.ApproxEdit{
				gostr.Insert, gostr.Insert, gostr.Match, gostr.Match, gostr.Match, gostr.Delete, gostr.Delete, gostr.Insert}},
			"2I3M2D1I",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := gostr.OpsToCigar(tt.args.ops); got != tt.want {
				t.Errorf("OpsToCigar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCigarToOps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cigar   string
		want    gostr.EditOps
		wantErr error
	}{
		{
			"1M",
			gostr.EditOps{gostr.Match},
			nil,
		},
		{
			"10M",
			gostr.EditOps{gostr.Match, gostr.Match, gostr.Match, 
				gostr.Match, gostr.Match, gostr.Match, gostr.Match, 
				gostr.Match, gostr.Match, gostr.Match},
			nil,
		},
		{
			"1I",
			gostr.EditOps{gostr.Insert},
			nil,
		},
		{
			"1D",
			gostr.EditOps{gostr.Delete},
			nil,
		},
		{
			"1D2M3I",
			gostr.EditOps{gostr.Delete, gostr.Match, gostr.Match, gostr.Insert, 
				gostr.Insert, gostr.Insert},
			nil,
		},
		{
			"invalid",
			gostr.EditOps{},
			gostr.NewInvalidCigar("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.cigar, func(t *testing.T) {
			t.Parallel()

			got, gotErr := gostr.CigarToOps(tt.cigar)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Unexpected error, %q", gotErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CigarToOps() = %v, want %v", got, tt.want)
			}

			if gotErr != nil && gotErr.Error() != "invalid cigar: "+tt.cigar {
				t.Errorf("Unexpected error message: %s", gotErr)
			}
		})
	}
}

func TestExtractAlignment(t *testing.T) {
	t.Parallel()

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
		wantErr  error
	}{
		{
			"Just matches",
			args{"acgtacgt", "gtac", 2, "4M"},
			"gtac", "gtac",
			nil,
		},
		{
			"Deletion",
			args{"acgtacgt", "gtc", 2, "2M1D1M"},
			"gtac", "gt-c",
			nil,
		},
		{
			"Insertion",
			args{"acgtacgt", "gtaac", 2, "2M1I2M"},
			"gt-ac", "gtaac",
			nil,
		},
		{
			"Invalid",
			args{"acgtacgt", "gtaac", 2, "invalid"},
			"", "",
			gostr.NewInvalidCigar("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotSubx, gotSubp, gotErr := gostr.ExtractAlignment(tt.args.x, tt.args.p, tt.args.pos, tt.args.cigar)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf("ExtractAlignment() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
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
	t.Parallel()

	type args struct {
		x     string
		p     string
		pos   int
		cigar string
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr error
	}{
		{
			"Just matches",
			args{"acgtacgt", "gtac", 2, "4M"},
			0, // "gtac" vs "gtac",
			nil,
		},
		{
			"Deletion",
			args{"acgtacgt", "gtc", 2, "2M1D1M"},
			1, // "gtac", "gt-c",
			nil,
		},
		{
			"Insertion",
			args{"acgtacgt", "gtaac", 2, "2M1I2M"},
			1, // "gt-ac", "gtaac",
			nil,
		},
		{
			"Invalid",
			args{"acgtacgt", "gtaac", 2, "invalid"},
			0, // error...
			gostr.NewInvalidCigar("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := gostr.CountEdits(tt.args.x, tt.args.p, tt.args.pos, tt.args.cigar)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf("Unexpected error %v", gotErr)
			}
			if got != tt.want {
				t.Errorf("CountEdits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type approxAlgo = func(string) func(string, int, func(int, string))

var approxAlgorithms = map[string]approxAlgo{ //nolint:gochecknoglobals // I'm fine with a global here
	"BWA": gostr.FMIndexApproxPreprocess,
}

func runRandomApproxOccurencesTests(algo approxAlgo) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		rng := testutils.NewRandomSeed(t)
		testutils.GenerateTestStringsAndPatterns(10, 20, rng,
			func(x, p string) {
				search := algo(x)
				for edits := 1; edits < 3; edits++ {
					search(p, edits, func(pos int, cigar string) {
						count, _ := gostr.CountEdits(x, p, pos, cigar)
						if count > edits {
							t.Log(pos, cigar)
							ax, ap, _ := gostr.ExtractAlignment(x, p, pos, cigar)
							t.Logf("%s\n%s\n\n", ax, ap)

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
	t.Parallel()

	for name, algo := range approxAlgorithms {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runRandomApproxOccurencesTests(algo)(t)
		})
		
		t.Run(name, runRandomApproxOccurencesTests(algo))
	}
}
