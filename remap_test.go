package gostr

import (
	"reflect"
	"testing"
)

func Test_Remap(t *testing.T) {
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
			gotX, gotAsize := Remap(tt.args.s)
			if !reflect.DeepEqual(gotX, tt.wantX) {
				t.Errorf("remap() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotAsize != tt.wantAsize {
				t.Errorf("remap() gotAsize = %v, want %v", gotAsize, tt.wantAsize)
			}
		})
	}
}
