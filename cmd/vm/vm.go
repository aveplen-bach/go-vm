package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	"github.com/aveplen/sm/internal"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input   string `short:"i" long:"input" description:"Input file name"`
	Verbose bool   `short:"v" long:"verbose" description:"Dump machine state on every instruction"`
}

func main() {
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		return
	}

	// input file
	fin, err := os.Open(opts.Input)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fin.Close(); err != nil {
			panic(err)
		}
	}()

	stat, err := fin.Stat()
	if err != nil {
		panic(err)
	}

	// program read
	program := make([]uint16, stat.Size()/2)
	if err := binary.Read(fin, binary.LittleEndian, program); err != nil {
		panic(err)
	}

	data := make([]uint16, 0, 10)
	for _, v := range args[1:] {
		uintmem, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			panic(err)
		}
		data = append(data, uint16(uintmem))
	}

	cpu := internal.WithMemProg(program, data)

	cpu.Run()

	if opts.Verbose {
		fmt.Println(cpu.Dump())
	}
}
