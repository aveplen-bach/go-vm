package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/aveplen/sm/internal"
)

func main() {
	input := os.Args[1]
	fin, err := os.Open(input)
	defer fin.Close()
	if err != nil {
		panic(err)
	}

	finr := bufio.NewReader(fin)
	program := internal.Compile(*finr)
	fmt.Println(len(program))

	output := os.Args[2]
	fout, err := os.Create(output)
	defer fout.Close()
	if err != nil {
		panic(err)
	}

	foutw := bufio.NewWriter(fout)
	binary.Write(foutw, binary.LittleEndian, program)
	foutw.Flush()
}
