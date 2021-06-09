package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func Test_MississippiBWT(t *testing.T) {
	x, _ := NewString("mississippi", nil)
	x_ := x.ToGoString() // FIXME
	sa := Skew(x)
	ctab := Ctab(x_)
	otab := Otab(x_, sa, ctab)
	p := "is"
	L, R := BwtSearch(x_, p, ctab, otab)
	for i := L; i < R; i++ {
		test.CheckOccurrenceAt(t, x_, p, sa[i])
	}

	preproc := BwtPreprocess(x)
	preproc(p, func(i int) {
		test.CheckOccurrenceAt(t, x_, p, i)
	})
}
