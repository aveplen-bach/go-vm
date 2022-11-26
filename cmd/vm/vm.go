package main

import (
	"encoding/binary"
	"log"
	"os"
	"strconv"

	"github.com/aveplen/sm/internal"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input   string `short:"i" long:"input" description:"Input file name"`
	Verbose bool   `short:"v" long:"verbose" description:"Dump machine state on every instruction"`
	Pause   int    `short:"p" long:"pause" description:"Lenght of the pause after command execution (ms)"`
}

func main() {
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// input file
	log.Printf("opening input file: %s\n", opts.Input)
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

	log.Println("reading program length as file metadata")
	stat, err := fin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// program read
	log.Println("reading program from file")
	program := make([]uint16, stat.Size()/2)
	if err := binary.Read(fin, binary.LittleEndian, program); err != nil {
		log.Fatal(err)
	}

	log.Println("creating memory preseeded with passed array")
	meminit := make([]uint16, 0, 10)
	for _, v := range args[1:] {
		uintmem, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		meminit = append(meminit, uint16(uintmem))
	}

	log.Println("creating cpu instance with memory and program")
	cpu := internal.WithMemProg(meminit, program)

	log.Println("creating run options array")
	runOpts := make([]internal.RunOpt, 0)

	if opts.Pause > 0 {
		runOpts = append(runOpts, internal.WithPause(opts.Pause))
	}

	if opts.Verbose {
		runOpts = append(runOpts, internal.WithVerbose())
	}

	log.Println("starting execution")
	cpu.Run(runOpts...)
}
