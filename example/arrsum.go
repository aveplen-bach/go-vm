package example

import (
	"github.com/aveplen/sm/emu"
)

// required initial memory layout:
//
// array length
//
//	|
//	v
//	5 1 2 3 4 5 ...
//	  ^
//	  |
//	result
func ArraySum(arr []uint32) uint32 {

	program := []uint32{
		// sum up here
		/* 00 */ emu.PUSH,
		/* 01 */ 0,

		// load array length
		/* 02 */ emu.PUSH,
		/* 03 */ 0,
		/* 04 */ emu.LOAD,

		/* 05 */ emu.DUP,

		// push finish addr
		/* 06 */ emu.PUSH,
		/* 07 */ 34,
		/* 08 */ emu.PUSH,
		/* 09 */ 1,
		/* 10 */ emu.ADD,
		/* 11 */ emu.PUSH,
		/* 12 */ 0,
		/* 13 */ emu.LOAD,
		/* 14 */ emu.ADD,
		/* 15 */ emu.SWAP,

		// if len(arr) == 0 skip iteration
		/* 16 */ emu.JZ,

		// decrement iterator
		/* 17 */ emu.DUP,
		/* 18 */ emu.LOAD,
		/* 19 */ emu.SWAP,
		/* 20 */ emu.PUSH,
		/* 21 */ 1,
		/* 22 */ emu.SUB,

		// add element to result
		/* 23 */ emu.SWAP,
		/* 24 */ emu.ROL3,
		/* 25 */ emu.ADD,
		/* 26 */ emu.SWAP,

		// loop
		/* 27 */ emu.PUSH,
		/* 28 */ 6,
		/* 29 */ emu.PUSH,
		/* 30 */ 0,
		/* 31 */ emu.LOAD,
		/* 32 */ emu.ADD,
		// goto 6
		/* 33 */ emu.JMP,

		// save result into memory[1]
		/* 34 */ emu.PUSH,
		/* 35 */ 1,
		/* 36 */ emu.ADD,
		/* 37 */ emu.STOR,

		// infinite loop
		/* 38 */ emu.PUSH,
		/* 39 */ 38,
		/* 40 */ emu.PUSH,
		/* 41 */ 1,
		/* 42 */ emu.ADD,
		/* 43 */ emu.PUSH,
		/* 44 */ 0,
		/* 45 */ emu.LOAD,
		/* 46 */ emu.ADD,
		/* 47 */ emu.JMP,
	}

	memory := append([]uint32{uint32(len(arr))}, arr...)
	start := len(memory)

	memory = append(memory, program...)

	cpu := emu.WithMem(memory, start)

	for i := 0; i < 1000; i++ {
		cpu.Tick()
	}

	dump := cpu.MemDump()
	return dump[1]
}
