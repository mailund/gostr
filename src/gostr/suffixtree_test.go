package gostr

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNaiveConstruction(t *testing.T) {
	st := NaiveST("mississippi")

	f, err := os.Create("naive-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()

	labels := st.Root.LeafLabels()
	for _, i := range labels {
		fmt.Println(i)
	}
}

func TestMcCreightConstruction(t *testing.T) {
	st := McCreight("mississippi")

	f, err := os.Create("mccreight-dot.dot")
	if err != nil {
		log.Fatal(err)
	}
	st.ToDot(f)
	f.Close()

	labels := st.Root.LeafLabels()
	for _, i := range labels {
		fmt.Println(i)
	}

	fmt.Println("Everything went fine! Don't worry.")
	t.Error("testing")
}
