package main

import (
	"encoding/binary"
	"os"
	"strconv"

	"github.com/aveplen/sm/internal"
)

func main() {
	input := os.Args[1]
	fin, err := os.Open(input)
	defer fin.Close()
	if err != nil {
		panic(err)
	}

	stat, err := fin.Stat()
	if err != nil {
		panic(err)
	}

	program := make([]uint32, stat.Size()/4)
	if err := binary.Read(fin, binary.LittleEndian, program); err != nil {
		panic(err)
	}

	strmeminit := os.Args[2:]
	meminit := make([]uint32, 0, 10)
	for _, v := range strmeminit {
		uintmem, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			panic(err)
		}

		meminit = append(meminit, uint32(uintmem))
	}

	cpu := internal.WithMemProg(meminit, program)
	cpu.Run()
}
