package internal

var program = []uint16{
	// load len(arr)
	/* 00 */ PUSH,
	/* 01 */ 0,
	/* 02 */ LOAD,
	/* 03 */ DUP,

	// if len(arr) == 0 { terminate }
	/* 04 */ PUSH,
	/* 05 */ 38, // sum routine addr
	/* 06 */ SWAP,
	/* 07 */ JZ,

	/* 08 */ STC,

	// load arr onto the stack
	/* 09 */ CTS,
	/* 10 */ LOAD,
	/* 11 */ CDEC,
	/* 12 */ CTS,
	/* 13 */ PUSH,
	/* 14 */ 20, // sum routine addr
	/* 15 */ SWAP,
	/* 16 */ JZ,
	/* 17 */ PUSH,
	/* 18 */ 9, // repeat
	/* 19 */ JMP,

	// load len(arr)
	/* 20 */ PUSH,
	/* 21 */ 0,
	/* 22 */ LOAD,
	/* 23 */ STC,
	/* 24 */ CDEC,

	// sum(stack)
	/* 25 */ CTS,
	/* 26 */ PUSH,
	/* 27 */ 35, // end routine addr
	/* 28 */ SWAP,
	/* 29 */ JZ,
	/* 30 */ ADD,
	/* 31 */ CDEC,
	/* 32 */ PUSH,
	/* 33 */ 25, // repeat
	/* 34 */ JMP,

	// move top to mem[0]
	/* 35 */ PUSH,
	/* 36 */ 0,
	/* 37 */ STOR,

	// terminate
	/* 38 */ TERM,
}

// required memory layout:
//
//	      5 1 2 3 4 5 ...
//	      ^
//	      |
//	length/result
func ArraySum(arr []int) int {
	memory := make([]uint16, 0, len(arr)+1)
	memory = append(memory, uint16(len(arr)))
	for _, v := range arr {
		memory = append(memory, uint16(v))
	}

	cpu := WithMemProg(memory, program)
	cpu.Run()

	return int(cpu.MemDump()[0])
}
