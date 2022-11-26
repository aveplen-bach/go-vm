package internal

import (
	"fmt"
)

const (
	MEMORY_DUMP_WIDTH = 8
	STACK_DUMP_WIDTH  = 8
)

func PrintDump(cpu *cpu) {
	// memory header
	fmt.Print("+")
	for i := 0; i < 117; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	fmt.Print("|   memory   ||")
	for i := 0; i < MEMORY_DUMP_WIDTH; i++ {
		fmt.Printf("         +%d |", i)
	}
	fmt.Print("\n")

	fmt.Print("+")
	for i := 0; i < 117; i++ {
		fmt.Print("=")
	}
	fmt.Print("+\n")

	// memory dump
	memdump := cpu.MemDump()
	for i := 0; i < MemSize; i++ {
		if i%MEMORY_DUMP_WIDTH == 0 {
			fmt.Printf("| %#08x || ", i)
		}

		fmt.Printf("%#08x", memdump[i])

		if i%MEMORY_DUMP_WIDTH != MEMORY_DUMP_WIDTH-1 {
			fmt.Print(" | ")
			continue
		}

		fmt.Print(" |\n")
	}

	// memory border bottom
	fmt.Print("+")
	for i := 0; i < 117; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	fmt.Println("stack:")
	stackdump := cpu.StackDump()
	for i := 0; i < StackLimit; i++ {
		fmt.Printf("%#08x", stackdump[i])

		if i%STACK_DUMP_WIDTH != STACK_DUMP_WIDTH-1 {
			fmt.Print(" | ")
			continue
		}

		fmt.Print("\n")
	}
	fmt.Println()

	fmt.Printf("sp: %d\n", cpu.GetSp())

	fmt.Printf("ip: %d\n", cpu.GetIp())
}
