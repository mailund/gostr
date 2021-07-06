package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/mailund/biof"
	"github.com/mailund/cli"
	"github.com/mailund/gostr"
)

func preprocFile(genomeFileName string) string {
	return genomeFileName + ".gostr_bwt"
}

type bwtPreprocArgs struct {
	GenomeName string `pos:"genome" descr:"FASTA file to preprocess."`
}

func initBwtPreproc() interface{} {
	return &bwtPreprocArgs{}
}

func bwtPreproc(i interface{}) {
	args, ok := i.(*bwtPreprocArgs)
	if !ok {
		panic("Unexpected arguments to bwtPreproc")
	}

	// We need the name of the file to make the preprocessed
	// file name, so we can't use an InFile here. We must
	// explicitly open and read the genome.
	genomeFile, err := os.Open(args.GenomeName)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s", args.GenomeName, err)
	}

	recs := getFastaRecords(genomeFile)
	genomeFile.Close()

	outFile, err := os.Create(preprocFile(args.GenomeName))
	if err != nil {
		log.Fatalf("Couldn't open %s: %s", preprocFile(args.GenomeName), err)
	}

	processed := map[string]*gostr.FMIndexTables{}
	for chrom, seq := range recs {
		// We always preprocess the full set, even if it is
		// slower if we intent to do exact matching
		processed[chrom] = gostr.BuildFMIndexApproxTables(seq)
	}

	enc := gob.NewEncoder(outFile)
	if err := enc.Encode(processed); err != nil {
		log.Fatal("encode error:", err)
	}

	if err := outFile.Close(); err != nil {
		log.Fatalf("Error closing file %s: %s", preprocFile(args.GenomeName), err)
	}
}

var bwtPreprocess = cli.NewCommand(cli.CommandSpec{
	Name:   "preproc",
	Short:  "preprocess bwt/fm-index",
	Long:   "Build tables for fast mapping.",
	Init:   initBwtPreproc,
	Action: bwtPreproc,
})

type bwtExactArgs struct {
	Out        cli.OutFile `flag:"out" short:"o" descr:"output file"`
	GenomeName string      `pos:"genome" descr:"FASTA file to preprocess."`
	Reads      cli.InFile  `pos:"reads" descr:"FASTQ file containing the reads"`
}

func initBwtExact() interface{} {
	return &bwtExactArgs{Out: cli.OutFile{Writer: os.Stdout}}
}

func readPreprocTables(fname string) map[string]*gostr.FMIndexTables {
	infile, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Couldn't open file: %s, did you remember to preprocess?",
			fname)
		os.Exit(1)
	}

	recs := map[string]*gostr.FMIndexTables{}
	dec := gob.NewDecoder(infile)

	if err := dec.Decode(&recs); err != nil {
		log.Fatalf("Error decoding preprocessing file %s: %s",
			fname, err)
	}

	infile.Close()

	return recs
}

func bwtExact(i interface{}) {
	args, ok := i.(*bwtExactArgs)
	if !ok {
		panic("Unexpected arguments to bwtExact")
	}

	recs := readPreprocTables(preprocFile(args.GenomeName))

	// Wrap the preprocessed tables in functions...
	searchFuncs := map[string]func(p string, cb func(i int)){}
	for rname, tbls := range recs {
		searchFuncs[rname] = gostr.FMIndexExactFromTables(tbls)
	}

	// FINALLY! we get to mapping. First the function for the mapping...
	mapper := func(rec *biof.FastqRecord) {
		for rname, search := range searchFuncs {
			search(rec.Read, func(pos int) {
				cigar := fmt.Sprintf("%dM", len(rec.Read))
				if err := biof.PrintSam(args.Out, rec.Name, rname, pos, cigar, rec.Read, rec.Qual); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to SAM file: %s\n", err)
				}
			})
		}
	}

	// ...then the actual scan
	if err := mapFastq(args.Reads, mapper); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning reads: %s\n", err)
		os.Exit(1)
	}

	// Finally, some clean up
	args.Reads.Close()

	if err := args.Out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing output file: %s", err)
		os.Exit(1)
	}
}

var bwtExactMap = cli.NewCommand(cli.CommandSpec{
	Name:   "exact",
	Short:  "exact mapping",
	Long:   "Search for exact matches in genome.",
	Init:   initBwtExact,
	Action: bwtExact,
})

type bwtApproxArgs struct {
	Out        cli.OutFile `flag:"out" short:"o" descr:"output file"`
	Edits      int         `flag:"edits" short:"d" descr:"maximun number of edits to consider"`
	GenomeName string      `pos:"genome" descr:"FASTA file to preprocess."`
	Reads      cli.InFile  `pos:"reads" descr:"FASTQ file containing the reads"`
}

func initBwtApprox() interface{} {
	return &bwtApproxArgs{
		Out:   cli.OutFile{Writer: os.Stdout},
		Edits: 1,
	}
}

func bwtApprox(i interface{}) {
	args, ok := i.(*bwtApproxArgs)
	if !ok {
		panic("Unexpected arguments to bwtApprox")
	}

	recs := readPreprocTables(preprocFile(args.GenomeName))

	// Wrap the preprocessed tables in functions...
	searchFuncs := map[string]func(p string, edits int, cb func(i int, cigar string)){}
	for rname, tbls := range recs {
		searchFuncs[rname] = gostr.FMIndexApproxFromTables(tbls)
	}

	// FINALLY! we get to mapping. First the function for the mapping...
	mapper := func(rec *biof.FastqRecord) {
		for rname, search := range searchFuncs {
			search(rec.Read, args.Edits, func(pos int, cigar string) {
				if err := biof.PrintSam(args.Out, rec.Name, rname, pos, cigar, rec.Read, rec.Qual); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to SAM file: %s\n", err)
				}
			})
		}
	}

	// ...then the actual scan
	if err := mapFastq(args.Reads, mapper); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning reads: %s\n", err)
		os.Exit(1)
	}

	// Finally, some clean up
	args.Reads.Close()

	if err := args.Out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing output file: %s", err)
		os.Exit(1)
	}
}

var bwtApproxMap = cli.NewCommand(cli.CommandSpec{
	Name:   "approx",
	Short:  "approximative mapping",
	Long:   "Search for approximative matches in genome.",
	Init:   initBwtApprox,
	Action: bwtApprox,
})

var bwt = cli.NewMenu(
	"bwt",
	"pattern matching using bwt/fm-index",
	"Search for exact and approximative matches of reads in a genome.",
	bwtPreprocess, bwtExactMap, bwtApproxMap)
