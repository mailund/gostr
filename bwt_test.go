package gostr_test

import (
	"reflect"
	"testing"

	"github.com/mailund/gostr"
	"github.com/mailund/gostr/test"
)

func Test_Ctab(t *testing.T) {
	x := "aab"
	alpha := gostr.NewAlphabet(x)
	xb, _ := alpha.MapToBytesWithSentinel(x)
	ctab := gostr.NewCTab(xb, alpha.Size())

	if len(ctab.CumSum) != alpha.Size() {
		t.Fatal("The ctable's cumsum has the wrong length")
	}

	if !reflect.DeepEqual(ctab.CumSum, []int{0, 1, 3}) {
		t.Fatal("We have the wrong cumsum")
	}
}

func Test_Otab(t *testing.T) {
	x := "aab"
	alpha := gostr.NewAlphabet(x)
	sa, _ := gostr.SaisWithAlphabet(x, alpha)

	xb, _ := alpha.MapToBytesWithSentinel(x)
	bwt := gostr.Bwt(xb, sa)
	otab := gostr.NewOTab(bwt, alpha.Size())

	expectedBwt := []byte{2, 0, 1, 1}
	if !reflect.DeepEqual(bwt, expectedBwt) {
		t.Fatalf("Expected bwt %v, got %v", expectedBwt, bwt)
	}

	expectedA := []int{0, 0, 0, 1, 2}
	expectedB := []int{0, 1, 1, 1, 1}

	var (
		a byte = 1
		b byte = 2
	)

	for i := range expectedA {
		if otab.Rank(a, i) != expectedA[i] {
			t.Errorf("Unexpected value at Rank(%b,%d) = %d\n", a, i, otab.Rank(a, i))
		}
	}

	for i := range expectedB {
		if otab.Rank(b, i) != expectedB[i] {
			t.Errorf("Unexpected value at Rank(%b,%d) = %d\n", b, i, otab.Rank(b, i))
		}
	}
}

func Test_BwtReverse(t *testing.T) {
	xs := "foobar"
	x, alpha := gostr.MapStringWithSentinel(xs)
	sa, _ := gostr.SaisWithAlphabet(xs, alpha)
	bwt := gostr.Bwt(x, sa)

	y := gostr.ReverseBwt(bwt)
	if !reflect.DeepEqual(x, y) {
		t.Fatalf("Expected %s == %s",
			alpha.RevmapBytes(x), alpha.RevmapBytes(y))
	}
}

func Test_MississippiBWT(t *testing.T) {
	xs := "mississippi"
	ps := "is"
	alpha := gostr.NewAlphabet(xs)
	x, _ := alpha.MapToBytesWithSentinel(xs)
	p, _ := alpha.MapToBytes(ps)

	sa, _ := gostr.SkewWithAlphabet(xs, alpha)
	bwt := gostr.Bwt(x, sa)
	ctab := gostr.NewCTab(bwt, alpha.Size())
	otab := gostr.NewOTab(bwt, alpha.Size())

	L, R := gostr.BwtSearch(x, p, ctab, otab)
	for i := L; i < R; i++ {
		test.CheckOccurrenceAt(t, xs, ps, int(sa[i]))
	}

	preproc := gostr.BwtPreprocess(xs)
	preproc(ps, func(i int) {
		test.CheckOccurrenceAt(t, xs, ps, i)
	})
}
