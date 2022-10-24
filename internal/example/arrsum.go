package example

import (
	"github.com/aveplen/sm/emu"
)

var program = []uint32{
	// sum up here
	/* 00 */ emu.PUSH,
	/* 01 */ 0,

	// load array length
	/* 02 */ emu.PUSH,
	/* 03 */ 0,
	/* 04 */ emu.LOAD,

	// while (top > 0) {
	/* 		05 */ emu.DUP,
	/* 		06 */ emu.PUSH,
	/* 		07 */ 23,
	/* 		08 */ emu.SWAP,
	/* 		09 */ emu.JZ,

	// 		decrement iterator
	/* 		10 */ emu.DUP,
	/* 		11 */ emu.LOAD,
	/* 		12 */ emu.SWAP,
	/* 		13 */ emu.PUSH,
	/* 		14 */ 1,
	/* 		15 */ emu.SUB,

	// 		add element to result
	/* 		16 */ emu.SWAP,
	/* 		17 */ emu.ROL3,
	/* 		18 */ emu.ADD,
	/* 		19 */ emu.SWAP,

	/* 		20 */ emu.PUSH,
	/* 		21 */ 5,
	/* 		22 */ emu.JMP,
	// }

	// save result into memory[0]
	/* 23 */ emu.STOR,

	// infinite loop
	/* 24 */ emu.PUSH,
	/* 25 */ 24,
	/* 26 */ emu.JMP,
}

func CpuWithArraySum(arr []int) emu.Cpu {
	memory := make([]uint32, 0, len(arr)+1)
	memory = append(memory, uint32(len(arr)))
	for _, v := range arr {
		memory = append(memory, uint32(v))
	}
	return emu.WithMem(memory, program)
}

// required memory layout:
//	      5 1 2 3 4 5 ...
//	      ^
//	      |
//	length/result
func ArraySum(arr []int) int {
	cpu := CpuWithArraySum(arr)
	for i := 0; i < 1000; i++ {
		cpu.Tick()
	}
	dump := cpu.MemDump()
	return int(dump[0])
}
