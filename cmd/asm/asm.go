package main

import (
	"bufio"
	"encoding/binary"
	"os"

	"github.com/aveplen/sm/internal"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input   string `short:"i" long:"input" description:"Input file"`
	Output  string `short:"o" long:"output" description:"Output file"`
	Verbose bool   `short:"v" long:"verbose" description:"Print list of tokens after compilation"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		return
	}

	// input file reader
	fin, err := os.Open(opts.Input)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fin.Close(); err != nil {
			panic(err)
		}
	}()

	finr := bufio.NewReader(fin)

	// output file reader
	fout, err := os.Create(opts.Output)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fout.Close(); err != nil {
			panic(err)
		}
	}()

	foutw := bufio.NewWriter(fout)
	defer func() {
		if err := foutw.Flush(); err != nil {
			panic(err)
		}
	}()

	// main compiler call
	program, err := internal.Compile(*finr, opts.Verbose)
	if err != nil {
		panic(err)
	}

	if err := binary.Write(foutw, binary.LittleEndian, program); err != nil {
		panic(err)
	}
}
