package gostr

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr/test"
)

func Test_Ctab(t *testing.T) {
	x, _ := NewString("aab", nil)
	ctab := Ctab(x)
	if x.Alpha != ctab.Alpha {
		t.Fatal("x and ctab should have the same alphabet")
	}
	if len(ctab.CumSum) != ctab.Alpha.Size() {
		t.Fatal("The ctable's cumsum has the wrong length")
	}
	if !reflect.DeepEqual(ctab.CumSum, []int{0, 1, 3}) {
		t.Fatal("We have the wrong cumsum")
	}
}

func Test_Otab(t *testing.T) {
	x, _ := NewString("aab", nil)
	sa := Sais(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)

	bwt := BwtString(x, sa)
	expected_bwt := string([]byte{2, 0, 1, 1})
	if !reflect.DeepEqual(bwt, expected_bwt) {
		t.Fatalf("Expected bwt %v, got %v", expected_bwt, bwt)
	}
	expected := []int{
		0, 0, 1, 2, // a row
		1, 1, 1, 1} // b row
	if !reflect.DeepEqual(otab.table, expected) {
		t.Fatalf("Unexpected otable: %v", otab.table)
	}
}

func Test_MississippiBWT(t *testing.T) {
	x_ := "mississippi"
	p_ := "is"
	x, _ := NewString(x_, nil)
	p, _ := NewString(p_, x.Alpha)
	sa := Skew(x)
	ctab := Ctab(x)
	otab := Otab(x, sa, ctab)
	L, R := BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		test.CheckOccurrenceAt(t, x_, p_, sa[i])
	}

	preproc := BwtPreprocess(x)
	preproc(p_, func(i int) {
		test.CheckOccurrenceAt(t, x_, p_, i)
	})
}
