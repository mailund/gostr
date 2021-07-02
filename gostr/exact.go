package main

import (
	"fmt"
	"os"

	"github.com/mailund/biof"
	"github.com/mailund/cli"
	"github.com/mailund/gostr"
)

var exactAlgos = map[string]func(x, p string, fn func(int)){
	"naive":  gostr.Naive,
	"border": gostr.BorderSearch,
	"kmp":    gostr.Kmp,
	"bmh":    gostr.Bmh,
}

type ExactArgs struct {
	Algo   cli.Choice  `flag:"algo" short:"a" descr:"choice of algorithm to use"`
	Out    cli.OutFile `flag:"out" short:"o" descr:"output file"`
	Genome cli.InFile  `pos:"genome" descr:"FASTA file containing the genome"`
	Reads  cli.InFile  `pos:"reads" descr:"FASTQ file containing the reads"`
}

func InitArgs() interface{} {
	exactAlgKeys := []string{}
	for k := range exactAlgos {
		exactAlgKeys = append(exactAlgKeys, k)
	}

	return &ExactArgs{
		Out:  cli.OutFile{Writer: os.Stdout},
		Algo: cli.Choice{Choice: "border", Options: exactAlgKeys},
	}
}

func GetFastaRecords(f cli.InFile) map[string]string {
	defer f.Close()

	recs, err := biof.ReadFasta(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading fasta file: %s\n", err.Error())
	}

	return recs
}

func MapFastq(f cli.InFile, fn func(*biof.FastqRecord)) error {
	defer f.Close()

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
