package gostr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ApproxEdit is a type for edit operations
type ApproxEdit int

// EditOps is a type for a sequence of edit operations
type EditOps = []ApproxEdit

// Approximative matching edit operations.
const (
	M ApproxEdit = iota // Match/mismatch operations
	I            = iota // Insertion operations
	D            = iota // Deletion operations
)

var editsToString = map[ApproxEdit]string{
	M: "M", I: "I", D: "D",
}

var stringToEdits = map[string]ApproxEdit{
	"M": M, "I": I, "D": D,
}

// OpsToCigar turns a list of ops into a cigar.
func OpsToCigar(ops EditOps) string {
	var (
		res  = []string{}
		i, j int
	)

	for ; i < len(ops); i = j {
		for j = i + 1; j < len(ops) && ops[i] == ops[j]; j++ {
		}

		res = append(res, fmt.Sprintf("%d%s", j-i, editsToString[ops[i]]))
	}

	return strings.Join(res, "")
}

// CigarToOps turns a cigar string into the list of edit ops.
func CigarToOps(cigar string) (EditOps, error) {
	r := regexp.MustCompile(`\d+[MID]`)
	ops := EditOps{}

	// This check is really inefficient, but I don't have time to
	// implement a better parser of cigars right now. A scan from
	// the beginning that chumps of digits would be a lot faster.
	// I'll leave that for later...
	for _, s := range r.Split(cigar, -1) {
		if s != "" {
			return ops, &InvalidCigar{x: cigar}
		}
	}

	for _, op := range r.FindAllString(cigar, -1) {
		rep, _ := strconv.Atoi(op[:len(op)-1])
		opcode := stringToEdits[string(op[len(op)-1])]

		for i := 0; i < rep; i++ {
			ops = append(ops, opcode)
		}
	}

	return ops, nil
}

// ExtractAlignment extracts a pairwise alignment from the reference, x,
// the read, p, the position and the edits cigar.
func ExtractAlignment(x, p string, pos int, cigar string) (subx, subp string, err error) {
	i, j := pos, 0

	ops, err := CigarToOps(cigar)
	if err != nil {
		return "", "", err
	}

	for _, op := range ops {
		switch op {
		case M:
			subx += string(x[i])
			subp += string(p[j])
			i++
			j++

		case I:
			subx += "-"
			subp += string(p[j])
			j++

		case D:
			subx += string(x[i])
			subp += "-"
			i++
		}
	}

	return subx, subp, nil
}

// CountEdits counts the number of edits in the local alignment between x and p
// specified by pos and cigar
func CountEdits(x, p string, pos int, cigar string) (int, error) {
	edits := 0

	subx, subp, err := ExtractAlignment(x, p, pos, cigar)
	if err != nil {
		return 0, err
	}

	for i := range subx {
		if subx[i] != subp[i] {
			edits++
		}
	}

	return edits, nil
}
