package gostr

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr/test"
)

func Test_Ctab(t *testing.T) {
	x := "aab"
	alpha := NewAlphabet(x)
	x_, _ := alpha.MapToBytesWithSentinel(x)
	ctab := NewCTab(x_, alpha.Size())
	if len(ctab.CumSum) != alpha.Size() {
		t.Fatal("The ctable's cumsum has the wrong length")
	}
	if !reflect.DeepEqual(ctab.CumSum, []int{0, 1, 3}) {
		t.Fatal("We have the wrong cumsum")
	}
}

func Test_Otab(t *testing.T) {
	x := "aab"
	alpha := NewAlphabet(x)
	sa, _ := SaisWithAlphabet(x, alpha)

	x_, _ := alpha.MapToBytesWithSentinel(x)
	bwt := Bwt(x_, sa)
	otab := NewOTab(bwt, alpha.Size())

	expected_bwt := []byte{2, 0, 1, 1}
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

func Test_BwtReverse(t *testing.T) {
	x_ := "foobar"
	x, alpha := MapStringWithSentinel(x_)
	sa, _ := SaisWithAlphabet(x_, alpha)
	bwt := Bwt(x, sa)

	y := ReverseBwt(bwt)
	if !reflect.DeepEqual(x, y) {
		t.Fatalf("Expected %s == %s",
			string(alpha.RevmapBytes(x)),
			string(alpha.RevmapBytes(y)))
	}
}

func Test_MississippiBWT(t *testing.T) {
	x_ := "mississippi"
	p_ := "is"
	alpha := NewAlphabet(x_)
	x, _ := alpha.MapToBytesWithSentinel(x_)
	p, _ := alpha.MapToBytes(p_)

	sa, _ := SkewWithAlphabet(x_, alpha)
	bwt := Bwt(x, sa)
	ctab := NewCTab(bwt, alpha.Size())
	otab := NewOTab(bwt, alpha.Size())

	L, R := BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		test.CheckOccurrenceAt(t, x_, p_, sa[i])
	}

	preproc := BwtPreprocess(x_)
	preproc(p_, func(i int) {
		test.CheckOccurrenceAt(t, x_, p_, i)
	})
}
