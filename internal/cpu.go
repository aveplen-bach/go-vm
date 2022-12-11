package internal

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MemSize    int = 80
	StackLimit int = 16
	dumpWidth      = 8
)

type handler func()

type cpu struct {
	stack   []uint16
	program []uint16
	data    []uint16
	cnt     uint16
	sp      int
	ip      int
	hmap    map[int]handler
	running bool
}

func NewCpu() *cpu {
	return newcpu()
}

func newcpu() *cpu {
	ret := &cpu{}
	ret.init()
	return ret
}

func (c *cpu) init() {
	c.initstack()
	c.initmem()
	c.initdata()
	c.initsp()
	c.inithmap()
	c.initrunning()
}

func (c *cpu) initstack() {
	c.stack = make([]uint16, StackLimit)
}

func (c *cpu) initdata() {
	c.data = make([]uint16, MemSize)
}

func (c *cpu) initmem() {
	c.program = make([]uint16, MemSize)
}

func (c *cpu) initsp() {
	c.sp = -1
}

func (c *cpu) inithmap() {
	c.hmap = map[int]handler{
		NOP:    c.inop,
		ADD:    c.iadd,
		SUB:    c.isub,
		AND:    c.iand,
		OR:     c.ior,
		XOR:    c.ixor,
		NOT:    c.inot,
		IN:     c.iin,
		OUT:    c.iout,
		LOAD:   c.iload,
		STOR:   c.istor,
		JMP:    c.ijmp,
		JZ:     c.ijz,
		JNZ:    c.ijnz,
		PUSH:   c.ipush,
		DUP:    c.idup,
		SWAP:   c.iswap,
		ROL3:   c.irol3,
		DROP:   c.idrop,
		COMPL:  c.icomp,
		CINC:   c.icinc,
		CDEC:   c.icdec,
		CTS:    c.icts,
		STC:    c.istc,
		TERM:   c.iterm,
		OUTNUM: c.ioutnum,
		MUL:    c.imul,
	}
}

func (c *cpu) initrunning() {
	c.running = true
}

func WithMemProg(program []uint16, data []uint16) *cpu {
	ret := newcpu()
	copy(ret.program, program)
	copy(ret.data, data)
	return ret
}

func (c *cpu) MemDump() []uint16 {
	dump := make([]uint16, MemSize)
	copy(dump, c.program)
	return dump
}

func (c *cpu) DataDump() []uint16 {
	dump := make([]uint16, MemSize)
	copy(dump, c.data)
	return dump
}

func (c *cpu) StackDump() []uint16 {
	dump := make([]uint16, StackLimit)
	copy(dump, c.stack)
	return dump
}

func (c *cpu) GetSp() int {
	return c.sp
}

func (c *cpu) GetIp() int {
	return c.ip
}

func (c *cpu) Run() {
	for c.running {
		c.tick()
	}
}

func border() string {
	line := "-"
	for i := 0; i < dumpWidth; i++ {
		line += "---------"
	}
	for i := 0; i < dumpWidth-1; i++ {
		line += "-"
	}
	return fmt.Sprintf("+%s+\n", line)
}

func dataHeader() string {
	res := "|  data  |"
	for i := 0; i < dumpWidth; i++ {
		res += fmt.Sprintf("     +%d |", i)
	}
	res += "\r\n"
	return res
}

func memHeader() string {
	res := "|  instr |"
	for i := 0; i < dumpWidth; i++ {
		res += fmt.Sprintf("     +%d |", i)
	}
	res += "\r\n"
	return res
}

func stackHeader() string {
	res := "| stack  |"
	for i := 0; i < dumpWidth; i++ {
		res += fmt.Sprintf("     +%d |", i)
	}
	res += "\r\n"
	return res
}

func formatNumber(num uint16) string {
	prefix := "0x"
	snum := fmt.Sprintf("%x", num)
	for len(prefix)+len(snum) < len("0x0000") {
		prefix += "0"
	}

	if num == 0 {
		return prefix + snum
	}

	return prefix + snum
}

func dtable(dump []uint16) string {
	res := ""
	for i := 0; i < len(dump); i++ {
		if i%dumpWidth == 0 {
			res += fmt.Sprintf("| %#04x |", i/dumpWidth*dumpWidth)
		}

		res += fmt.Sprintf(" %s |", formatNumber(dump[i]))

		if i%dumpWidth == dumpWidth-1 {
			res += "\r\n"
		}
	}
	return res
}

func (c *cpu) Dump() string {
	md := c.MemDump()
	dd := c.DataDump()
	sd := c.StackDump()

	res := ""
	res += border()
	res += memHeader()
	res += dtable(md)
	res += border()

	res += "\n"

	res += border()
	res += dataHeader()
	res += dtable(dd)
	res += border()

	res += "\n"

	res += border()
	res += stackHeader()
	res += dtable(sd)
	res += border()

	res += "\n"

	res += fmt.Sprintf("counter register: %d\n", c.cnt)
	res += fmt.Sprintf("stack pointer: %d\n", c.sp)
	res += fmt.Sprintf("instruction pointer: %d\n", c.ip)

	return res
}

func (c *cpu) tick() {
	if !c.running {
		panic("attempt to tick when not running")
	}
	fetched := c.fetch()
	decoded := c.decode(fetched)
	c.execute(decoded)
}

func (c *cpu) Tick() {
	c.tick()
}

func (c *cpu) fetch() uint16 {
	cmd := c.program[c.ip]
	c.ip++
	return cmd
}

func (c *cpu) decode(opcode uint16) uint16 {
	return opcode
}

func (c *cpu) execute(cmd uint16) {
	h, ok := c.hmap[int(cmd)]
	if !ok {
		panic("unknown command")
	}
	h()
}

func (c *cpu) push(x uint16) {
	if c.sp == StackLimit-1 {
		panic("stack overflow")
	}
	c.sp++
	c.stack[c.sp] = x
}

func (c *cpu) pop() uint16 {
	if c.sp == -1 {
		panic("stack underflow")
	}
	ret := c.stack[c.sp]
	c.stack[c.sp] = 0
	c.sp--
	return ret
}

func (c *cpu) terminate() {
	c.running = false
}

// do nothing
func (c *cpu) inop() {}

// pop a, pop b, push a + b
func (c *cpu) iadd() {
	c.push(c.pop() + c.pop())
}

// pop a, pop b, push b - a
// !!! may not be desired behaviour !!!
func (c *cpu) isub() {
	t := c.pop()
	nt := c.pop()
	c.push(nt - t)
}

// pop a, pop b, push a & b
func (c *cpu) iand() {
	a := c.pop()
	b := c.pop()
	c.push(a & b)
}

// pop a, pop b, push a | b
func (c *cpu) ior() {
	a := c.pop()
	b := c.pop()
	c.push(a | b)
}

// pop a, pop b, push a ^ b
func (c *cpu) ixor() {
	a := c.pop()
	b := c.pop()
	c.push(a ^ b)
}

// pop a, push !a
func (c *cpu) inot() {
	c.push(^c.pop())
}

// read one byte from stdin and push to the stack
func (c *cpu) iin() {
	in := bufio.NewReader(os.Stdin)
	b, err := in.ReadByte()
	if err != nil {
		panic(err)
	}
	c.push(uint16(b))
}

// write top of the stack into stdout
func (c *cpu) iout() {
	out := bufio.NewWriter(os.Stdout)
	if err := out.WriteByte(byte(c.pop())); err != nil {
		panic(err)
	}
}

// pop a, push word read from memory[a]
func (c *cpu) iload() {
	c.push(c.data[c.pop()])
}

// pop a, pop b, write b to memory[a]
func (c *cpu) istor() {
	a := c.pop()
	c.data[a] = c.pop()
}

// pop a, goto a
func (c *cpu) ijmp() {
	c.ip = int(c.pop())
}

// pop a, pop b, if a == 0 goto b
func (c *cpu) ijz() {
	a := c.pop()
	b := c.pop()
	if a == 0 {
		c.ip = int(b)
	}
}

// push next word
func (c *cpu) ipush() {
	c.push(c.program[c.ip])
	c.ip++
}

// duplicate stack top
func (c *cpu) idup() {
	val := c.pop()
	c.push(val)
	c.push(val)
}

// swap two top values
func (c *cpu) iswap() {
	a := c.pop()
	b := c.pop()
	c.push(a)
	c.push(b)
}

// (a, b, c) -> (b, c, a)
func (c *cpu) irol3() {
	cc := c.pop()
	b := c.pop()
	a := c.pop()
	c.push(b)
	c.push(cc)
	c.push(a)
}

// write stack top into stdin as number
func (c *cpu) ioutnum() {
	a := c.pop()
	fmt.Printf("%d\n", a)
	c.push(a)
}

// pop a, pop b, if a != 0 goto b
func (c *cpu) ijnz() {
	a := c.pop()
	b := c.pop()
	if a != 0 {
		c.ip = int(b) + MemSize/2
	}
}

// pop stack top
func (c *cpu) idrop() {
	c.pop()
}

// push stack top complement
func (c *cpu) icomp() {
	c.push(-c.pop())
}

// increment counter
func (c *cpu) icinc() {
	c.cnt++
}

// decrement counter
func (c *cpu) icdec() {
	c.cnt--
}

// move value from counter to stack
func (c *cpu) icts() {
	c.push(c.cnt)
}

// move value from stack to counter
func (c *cpu) istc() {
	c.cnt = c.pop()
}

// terminate execution
func (c *cpu) iterm() {
	c.terminate()
}

// pop a, pop b, push a*b
func (c *cpu) imul() {
	c.push(c.pop() * c.pop())
}
