package gostr

import (
	"log"
	"os"
	"testing"
)

func TestNaiveConstruction(t *testing.T) {
	st := NaiveST("mississippi")

	f, err := os.Create("dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()

	t.Error("testing")
}
