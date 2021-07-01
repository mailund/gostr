package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailund/biof"
	"github.com/mailund/cli"
	"github.com/mailund/cli/interfaces"
	"github.com/mailund/gostr"
)

// This can go in cli one day...
type Choice struct {
	Choice  string
	Options []string
}

func (c *Choice) Set(x string) error {
	for _, v := range c.Options {
		if v == x {
			c.Choice = v
			return nil
		}
	}

	return interfaces.ParseErrorf("%s is not a valid choice", x)
}

func (c *Choice) String() string {
	return c.Choice + ", choices: {" + strings.Join(c.Options, ",") + "}"
}

// FIXME: I *must* ensure that there is a default value for an output file,
// because nil writers will crash the program. Need a protocol in cli for that.
// But that is a problem for a later day.
type OutFile struct {
	io.Writer
	fname string
}

func (o *OutFile) Set(x string) error {
	f, err := os.Create(x)
	if err != nil {
		return interfaces.ParseErrorf("couldn't open file %s: %s", x, err)
	}

	o.Writer = f

	return nil
}

func (o *OutFile) String() string {
	switch o.Writer {
	case os.Stdout:
		return "stdout"
	case os.Stderr:
		return "stderr"
	default:
		return o.fname
	}
}

var exactAlgos = map[string]func(x, p string, fn func(int)){
	"naive":  gostr.Naive,
	"border": gostr.BorderSearch,
	"kmp":    gostr.Kmp,
	"bmh":    gostr.Bmh,
}

type ExactArgs struct {
	Algo   Choice  `flag:"algo" descr:"choice of algorithm to use"`
	Out    OutFile `flag:"o" descr:"output file"`
	Genome string  `pos:"genome" descr:"FASTA file containing the genome"`
	Reads  string  `pos:"reads" descr:"FASTQ file containing the reads"`
}

func InitArgs() interface{} {
	exactAlgKeys := []string{}
	for k := range exactAlgos {
		exactAlgKeys = append(exactAlgKeys, k)
	}

	return &ExactArgs{
		Out:  OutFile{Writer: os.Stdout},
		Algo: Choice{Choice: "border", Options: exactAlgKeys},
	}
}

func GetFastaRecords(fname string) map[string]string {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open fasta file: %s\n", err.Error())
		os.Exit(1)
	} else {
		defer f.Close()
	}

	recs, err := biof.ReadFasta(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading fasta file: %s\n", err.Error())
	}

	return recs
}

func MapFastq(fname string, fn func(*biof.FastqRecord)) error {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open fastq file: %s\n", err)
	} else {
		defer f.Close()
	}

	return biof.ScanFastq(f, fn)
}

func ExactMapping(i interface{}) {
	args, ok := i.(*ExactArgs)
	if !ok {
		panic("Unexpected arguments to ExactMapping")
	}

	algo := exactAlgos[args.Algo.Choice]
	recs := GetFastaRecords(args.Genome)

	mapper := func(rec *biof.FastqRecord) {
		for rname, seq := range recs {
			algo(seq, rec.Read, func(pos int) {
				cigar := fmt.Sprintf("%dM", len(rec.Read))
				if err := biof.PrintSam(args.Out, rec.Name, rname, pos, cigar, rec.Read, rec.Qual); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to SAM file: %s\n", err)
				}
			})
		}
	}

	if err := MapFastq(args.Reads, mapper); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning reads: %s\n", err)
	}
}

var exact = cli.NewCommand(cli.CommandSpec{
	Name:   "exact",
	Short:  "exact pattern matching",
	Long:   "Search for exact matches of reads in a genome.",
	Init:   InitArgs,
	Action: ExactMapping,
})
