package gostr

import (
	"reflect"
	"testing"
)

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
			recSais(tt.args.x, tt.args.SA, tt.args.asize, isS)
		})
	}
}

func Test_SaisBasic(t *testing.T) {
	type args struct {
		s   string
		inc bool
	}
	tests := []struct {
		name   string
		args   args
		wantSA []int
	}{
		{`We handle empty strings`, args{"", false}, []int{}},
		{`Unique characters "a"`, args{"a", false}, []int{0}},
		{`Unique characters "a"`, args{"a", true}, []int{1, 0}},
		{`Unique characters "ab"`, args{"ab", false}, []int{0, 1}},
		{`Unique characters "ab"`, args{"ab", true}, []int{2, 0, 1}},
		{`Unique characters "ba"`, args{"ba", false}, []int{1, 0}},
		{`Unique characters "ba"`, args{"ba", true}, []int{2, 1, 0}},
		{`Unique characters "abc"`, args{"abc", false}, []int{0, 1, 2}},
		{`Unique characters "abc"`, args{"abc", true}, []int{3, 0, 1, 2}},
		{`Unique characters "bca"`, args{"bca", false}, []int{2, 0, 1}},
		{`Unique characters "bca"`, args{"bca", true}, []int{3, 2, 0, 1}},
		{`mississippi`, args{"mississippi", false},
			[]int{10, 7, 4, 1, 0, 9, 8, 6, 3, 5, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSA := Sais(tt.args.s, tt.args.inc); !reflect.DeepEqual(gotSA, tt.wantSA) {
				t.Errorf("SAIS() = %v, want %v", gotSA, tt.wantSA)
			}
		})
	}
}

func Test_SaisMississippi(t *testing.T) {
	x := "mississippi"
	testSASorted(x, Sais(x, false), t)
}

func Test_SaisRandomStrings(t *testing.T) {
	testRandomSASorted(Sais, t)
}
