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
				if got.get(int32(i)) != b {
					t.Errorf("newBitArray(%d,%v), want %v at index %d but got %v",
						len(tt.bits), tt.bits, b, i, got.get(int32(i)))
				}
			}
		})
	}

	for i := 0; i < 67; i++ {
		bv := newBitArray(i)

		if 8*len(bv.bytes) < i {
			t.Errorf("There are not enough bytes (%d) in the bit-array to hold %d bits.\n",
				len(bv.bytes), i)
		}

		if i <= 8*(len(bv.bytes)-1) {
			t.Errorf("There are too many bytes (%d) in the bit-array to hold %d bits.\n",
				len(bv.bytes), i)
		}
	}
}

// These are tested through the newBitArray test for now...
/* Not currently used
func Test_bitArray_get(t *testing.T) {}
func Test_bitArray_set(t *testing.T) {}
*/

func Test_classifyST(t *testing.T) {
	type args struct {
		x []int32
	}

	tests := []struct {
		name string
		args args
		want []bool
	}{
		{`String "$"`, args{[]int32{0}}, []bool{true}},
		{`String "a$"`, args{[]int32{1, 0}}, []bool{false, true}},
		{`String "ab$"`, args{[]int32{1, 2, 0}}, []bool{true, false, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isS := newBitArray(len(tt.args.x))
			classifyS(isS, tt.args.x)
			if isS.length != len(tt.want) {
				t.Errorf("classifyS() = %v has the wrong length (want %v)", isS, len(tt.want))
			}
			for i, b := range tt.want {
				if isS.get(int32(i)) != b {
					t.Errorf("classifyS() = %v, bit %d should be %v but is %v",
						isS, i, b, isS.get(int32(i)))
				}
			}
		})
	}
}

func Test_countBuckets(t *testing.T) {
	type args struct {
		x     []int32
		asize int
	}

	tests := []struct {
		name string
		args args
		want []int32
	}{
		{`Sentinel string`, args{[]int32{0}, 1}, []int32{1}},
		{`"abc$"`, args{[]int32{1, 2, 3, 0}, 4}, []int32{1, 1, 1, 1}},
		{`"aba$"`, args{[]int32{1, 2, 1, 0}, 3}, []int32{1, 2, 1}},
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
		buckets []int32
	}

	tests := []struct {
		name string
		args args
		want []int32
	}{
		{`Singleton`, args{[]int32{1}}, []int32{0}},
		{`[1, 2, 3]`, args{[]int32{1, 2, 3}}, []int32{0, 1, 3}},
	}

	for _, tt := range tests {
		buf := make([]int32, len(tt.args.buckets))
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
		buckets []int32
	}

	tests := []struct {
		name string
		args args
		want []int32
	}{
		{`Singleton`, args{buckets: []int32{1}}, []int32{1}},
		{`[1, 2, 3]`, args{buckets: []int32{1, 2, 3}}, []int32{1, 3, 6}},
	}

	for _, tt := range tests {
		buf := make([]int32, len(tt.args.buckets))
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
		i   int32
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
		x     []int32
		SA    []int32
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
			t.Parallel()
			
			isS := newBitArray(len(tt.args.x))
			recSais(tt.args.x, tt.args.SA, tt.args.asize, isS)
		})
	}
}

func Test_equalLMS(t *testing.T) {
	type args struct {
		x []int32
		i int32
		j int32
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"both end of string",
			args{[]int32{2, 1, 1, 0}, 4, 4},
			true},
		{"first end of string",
			args{[]int32{2, 1, 1, 0}, 4, 0},
			false},
		{"second end of string",
			args{[]int32{2, 1, 1, 0}, 0, 4},
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isS := newBitArray(len(tt.args.x))
			classifyS(isS, tt.args.x)
			if got := equalLMS(tt.args.x, isS, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("equalLMS() = %v, want %v", got, tt.want)
			}
		})
	}
}
