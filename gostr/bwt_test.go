package gostr_test

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/mailund/gostr/gostr"
	"github.com/mailund/gostr/testutils"
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

	preproc := gostr.FMIndexExactPreprocess(xs)
	preproc(ps, func(i int) {
		testutils.CheckOccurrenceAt(t, xs, ps, i)
	})
}

func Test_MississippiBWTApprox0(t *testing.T) {
	xs := "mississippi"
	ps := "is"

	preproc := gostr.FMIndexApproxPreprocess(xs)
	preproc(ps, 0, func(i int, _ string) {
		testutils.CheckOccurrenceAt(t, xs, ps, i)
	})
}

func TestOTabEncoding(t *testing.T) {
	x := "aab"
	alpha := gostr.NewAlphabet(x)
	sa, _ := gostr.SaisWithAlphabet(x, alpha)
	xb, _ := alpha.MapToBytesWithSentinel(x)
	bwt := gostr.Bwt(xb, sa)

	otab1 := gostr.NewOTab(bwt, alpha.Size())
	otab2 := &gostr.OTab{}

	if reflect.DeepEqual(otab1, otab2) {
		t.Fatalf("The two otables should not be equal yet")
	}

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
		dec = gob.NewDecoder(&buf)
	)

	if err := enc.Encode(&otab1); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := dec.Decode(&otab2); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(otab1, otab2) {
		t.Errorf("These two otables should be equal now")
	}
}
