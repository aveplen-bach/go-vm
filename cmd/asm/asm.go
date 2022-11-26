package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"

	"github.com/aveplen/sm/internal"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input  string `short:"i" long:"input" description:"Input file"`
	Output string `short:"o" long:"output" description:"Output file"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		log.Fatal(err)
	}

	// input file reader
	log.Printf("opening input file %s\n", opts.Input)
	fin, err := os.Open(opts.Input)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Println("closing input file")
		if err := fin.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	finr := bufio.NewReader(fin)

	// output file reader
	log.Printf("opening output file %s\n", opts.Output)
	fout, err := os.Create(opts.Output)
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Println("closing output file")
		if err := fout.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	foutw := bufio.NewWriter(fout)
	defer func() {
		log.Println("flushing output file")
		if err := foutw.Flush(); err != nil {
			log.Fatal(err)
		}
	}()

	// main compiler call
	log.Println("compiling program")
	program, err := internal.Compile(*finr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("writing compiled program")
	if err := binary.Write(foutw, binary.LittleEndian, program); err != nil {
		log.Fatal(err)
	}
}
