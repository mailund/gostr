package gostr

import (
	"testing"

	"github.com/mailund/gostr/test"
)

func Test_MississippiBWT(t *testing.T) {
	x := "mississippi"
	sa := Skew(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	p := "is"
	L, R := BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		test.CheckOccurrenceAt(t, x, p, sa[i])
	}

	preproc := BwtPreprocess(x)
	preproc(p, func(i int) {
		test.CheckOccurrenceAt(t, x, p, i)
	})
}
