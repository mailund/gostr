package gostr

import "testing"

func basicExactTests(
	algoname string,
	exact func(x, p string, cb func(int)),
	t *testing.T) {
	res := []int{}
	append_res := func(i int) {
		res = append(res, i)
	}

	type args struct {
		x        string
		p        string
		callback func(i int)
	}
	tests := []struct {
		name     string
		args     args
		expected []int
	}{
		{"aaa/",
			args{"aaa", "", append_res},
			[]int{0, 1, 2, 3},
		},
		{"aaa/a",
			args{"aaa", "a", append_res},
			[]int{0, 1, 2},
		},
		{"aaa/b",
			args{"aaa", "b", append_res},
			[]int{},
		},
		{"aaa/aa",
			args{"aaa", "aa", append_res},
			[]int{0, 1},
		},
		{"mississippi/ssi",
			args{"mississippi", "ssi", append_res},
			[]int{2, 5},
		},
		{"mississippi/ppi",
			args{"mississippi", "ppi", append_res},
			[]int{8},
		},
	}
	for _, tt := range tests {
		t.Run(algoname+":"+tt.name, func(t *testing.T) {
			exact(tt.args.x, tt.args.p, tt.args.callback)
			if !equal_arrays(tt.expected, res) {
				t.Errorf("Searching for %s in %s and found %v (expected %v)\n",
					tt.args.p, tt.args.x, res, tt.expected)
			}
			res = []int{} // Reset
		})
	}
}

func TestNaive(t *testing.T) { basicExactTests("Naive", Naive, t) }
