package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/parser"
)

func main() {
	from := flag.String("from", "", "start of time range (RFC3339, required)")
	to := flag.String("to", "", "end of time range (RFC3339, required)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice -from <ts> -to <ts> [file...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *from == "" || *to == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := parser.NewFilter(*from, *to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid time range: %v\n", err)
		os.Exit(1)
	}

	files := flag.Args()
	if len(files) == 0 {
		runReader(f, os.Stdin)
		return
	}

	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v\n", path, err)
			os.Exit(1)
		}
		runReader(f, file)
		file.Close()
	}
}

func runReader(f parser.Filter, r *os.File) {
	s := parser.NewScanner(r)
	lines := f.Slice(s)
	for _, l := range lines {
		fmt.Println(l)
	}
}
