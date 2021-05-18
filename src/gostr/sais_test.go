package gostr

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func Test_remap(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantX     []int
		wantAsize int
	}{
		{`Handles empty strings`, args{""}, []int{0}, 1},
		{`Handles "abc"`, args{"abc"}, []int{1, 2, 3, 0}, 4},
		{`Handles "abab"`, args{"abab"}, []int{1, 2, 1, 2, 0}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotAsize := remap(tt.args.s)
			if !reflect.DeepEqual(gotX, tt.wantX) {
				t.Errorf("remap() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotAsize != tt.wantAsize {
				t.Errorf("remap() gotAsize = %v, want %v", gotAsize, tt.wantAsize)
			}
		})
	}
}

func Test_newBitArray(t *testing.T) {
	tests := []struct {
		name string
		bits []bool
	}{
		{"Empty bit vector", []bool{}},
		{"0", []bool{false}},
		{"1", []bool{true}},
		{"010", []bool{false, true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newBitArray(len(tt.bits), tt.bits...)
			if got.length != len(tt.bits) {
				t.Errorf("newBitArray(%d,%v) got the wrong size. Expected %d, got %d",
					len(tt.bits), tt.bits, len(tt.bits), got.length)
			}
			for i, b := range tt.bits {
				if got.get(i) != b {
					t.Errorf("newBitArray(%d,%v), want %v at index %d but got %v",
						len(tt.bits), tt.bits, b, i, got.get(i))
				}
			}
		})
	}
}

// These are tested through the newBitArray test for now...
func Test_bitArray_get(t *testing.T) {}
func Test_bitArray_set(t *testing.T) {}

func Test_classifyST(t *testing.T) {
	type args struct {
		x []int
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{`String "$"`, args{[]int{0}}, []bool{true}},
		{`String "a$"`, args{[]int{1, 0}}, []bool{false, true}},
		{`String "ab$"`, args{[]int{1, 2, 0}}, []bool{true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isS := newBitArray(len(tt.args.x))
			classifyS(isS, tt.args.x)
			if isS.length != len(tt.want) {
				t.Errorf("classifyS() = %v has the wrong length (want %v)", isS, len(tt.want))
			}
			for i, b := range tt.want {
				if isS.get(i) != b {
					t.Errorf("classifyS() = %v, bit %d should be %v but is %v",
						isS, i, b, isS.get(i))
				}
			}
		})
	}
}

func Test_countBuckets(t *testing.T) {
	type args struct {
		x     []int
		asize int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{`Sentinel string`, args{[]int{0}, 1}, []int{1}},
		{`"abc$"`, args{[]int{1, 2, 3, 0}, 4}, []int{1, 1, 1, 1}},
		{`"aba$"`, args{[]int{1, 2, 1, 0}, 3}, []int{1, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countBuckets(tt.args.x, tt.args.asize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countBuckets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bucketsFronts(t *testing.T) {
	type args struct {
		buckets []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{`Singleton`, args{[]int{1}}, []int{0}},
		{`[1, 2, 3]`, args{[]int{1, 2, 3}}, []int{0, 1, 3}},
	}
	for _, tt := range tests {
		buf := make([]int, len(tt.args.buckets))
		t.Run(tt.name, func(t *testing.T) {
			bucketsFronts(buf, tt.args.buckets)
			if !reflect.DeepEqual(buf, tt.want) {
				t.Errorf("bucketsBegin() = %v, want %v", buf, tt.want)
			}
		})
	}
}

func Test_bucketsEnd(t *testing.T) {
	type args struct {
		buckets []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{`Singleton`, args{buckets: []int{1}}, []int{1}},
		{`[1, 2, 3]`, args{buckets: []int{1, 2, 3}}, []int{1, 3, 6}},
	}
	for _, tt := range tests {
		buf := make([]int, len(tt.args.buckets))
		t.Run(tt.name, func(t *testing.T) {
			bucketsEnd(buf, tt.args.buckets)
			if !reflect.DeepEqual(buf, tt.want) {
				t.Errorf("bucketsEnd() = %v, want %v", buf, tt.want)
			}
		})
	}
}

func Test_isLMS(t *testing.T) {
	type args struct {
		isS *bitArray
		i   int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLMS(tt.args.isS, tt.args.i); got != tt.want {
				t.Errorf("isLMS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recSAIS(t *testing.T) {
	type args struct {
		x     []int
		SA    []int
		asize int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isS := newBitArray(len(tt.args.x))
			recSAIS(tt.args.x, tt.args.SA, tt.args.asize, isS)
		})
	}
}

func TestSAIS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		args   args
		wantSA []int
	}{
		{`We handle empty strings`, args{""}, []int{}},
		{`Unique characters "a"`, args{"a"}, []int{0}},
		{`Unique characters "ab"`, args{"ab"}, []int{0, 1}},
		{`Unique characters "ba"`, args{"ba"}, []int{1, 0}},
		{`Unique characters "abc"`, args{"abc"}, []int{0, 1, 2}},
		{`Unique characters "bca"`, args{"bca"}, []int{2, 0, 1}},
		{`mississippi`, args{"mississippi"},
			[]int{10, 7, 4, 1, 0, 9, 8, 6, 3, 5, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSA := SAIS(tt.args.s); !reflect.DeepEqual(gotSA, tt.wantSA) {
				t.Errorf("SAIS() = %v, want %v", gotSA, tt.wantSA)
			}
		})
	}
}

func TestMississippiSAIS(t *testing.T) {
	x := "mississippi"
	testSASorted(x, SAIS(x), t)
}

func TestRandomStringsSAIS(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	n := 30       // testing 30 random strings, enough to hit all mod 3 lengths
	maxlen := 100 // max length 100 (so we can still inspect them)
	for i := 0; i < n; i++ {
		slen := rng.Intn(maxlen)
		x := randomString(slen, "acgt", rng)
		t.Logf(`Testing string "%s".`, x)
		testSASorted(x, SAIS(x), t)
	}
}

func benchmarkSAconstruction(
	constr func(string) []int,
	n int,
	b *testing.B) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < b.N; i++ {
		constr(randomString(n, "abcdefg", rng))
	}
}

func BenchmarkSkew10000(b *testing.B) {
	benchmarkSAconstruction(Skew, 10000, b)
}

func BenchmarkSkew100000(b *testing.B) {
	benchmarkSAconstruction(Skew, 100000, b)
}

func BenchmarkSkew1000000(b *testing.B) {
	benchmarkSAconstruction(Skew, 1000000, b)
}

func BenchmarkSAIS10000(b *testing.B) {
	benchmarkSAconstruction(SAIS, 10000, b)
}
func BenchmarkSAIS100000(b *testing.B) {
	benchmarkSAconstruction(SAIS, 100000, b)
}

func BenchmarkSAIS1000000(b *testing.B) {
	benchmarkSAconstruction(SAIS, 1000000, b)
}
