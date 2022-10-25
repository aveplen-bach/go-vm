package internal

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MEM_SIZE    int = 200
	STACK_LIMIT int = 16
)

type handler func()

type cpu struct {
	stack   []uint32
	memory  []uint32
	cnt     uint32
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
	c.initsp()
	c.initip()
	c.inithmap()
	c.initrunning()
}

func (c *cpu) initstack() {
	c.stack = make([]uint32, STACK_LIMIT)
}

func (c *cpu) initmem() {
	c.memory = make([]uint32, MEM_SIZE)
}

func (c *cpu) initsp() {
	c.sp = 0
}

func (c *cpu) initip() {
	c.ip = MEM_SIZE / 2
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
		PUSH:   c.ipush,
		DUP:    c.idup,
		SWAP:   c.iswap,
		ROL3:   c.irol3,
		OUTNUM: c.ioutnum,
		JNZ:    c.ijnz,
		DROP:   c.idrop,
		COMPL:  c.icomp,
		CINC:   c.icinc,
		CDEC:   c.icdec,
		CTS:    c.icts,
		STC:    c.istc,
		TERM:   c.iterm,
	}
}

func (c *cpu) initrunning() {
	c.running = true
}

func WithMemProg(memory []uint32, program []uint32) *cpu {
	ret := newcpu()
	copy(ret.memory, memory)
	for i, v := range program {
		ret.memory[i+MEM_SIZE/2] = v
	}
	return ret
}

func (c *cpu) MemDump() []uint32 {
	dump := make([]uint32, MEM_SIZE)
	copy(dump, c.memory)
	return dump
}

func (c *cpu) StackDump() []uint32 {
	dump := make([]uint32, STACK_LIMIT)
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

func (c *cpu) fetch() uint32 {
	cmd := c.memory[c.ip]
	c.ip++
	return cmd
}

func (c *cpu) decode(opcode uint32) uint32 {
	return opcode
}

func (c *cpu) execute(cmd uint32) {
	h, ok := c.hmap[int(cmd)]
	if !ok {
		panic("unknown command")
	}
	h()
}

func (c *cpu) push(n uint32) {
	if c.sp == STACK_LIMIT {
		panic("stack overflow")
	}
	c.stack[c.sp] = n
	c.sp++
}

func (c *cpu) pop() uint32 {
	if c.sp == 0 {
		panic("stack underflow")
	}
	c.sp--
	return c.stack[c.sp]
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
	c.push(uint32(b))
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
	a := c.pop()
	c.push(c.memory[a])
}

// pop a, pop b, write b to memory[a]
func (c *cpu) istor() {
	a := c.pop()
	b := c.pop()
	c.memory[a] = b
}

// pop a, goto a
func (c *cpu) ijmp() {
	c.ip = int(c.pop()) + MEM_SIZE/2
}

// pop a, pop b, if a == 0 goto b
func (c *cpu) ijz() {
	a := c.pop()
	b := c.pop()
	if a == 0 {
		c.ip = int(b) + MEM_SIZE/2
	}
}

// push next word
func (c *cpu) ipush() {
	c.push(c.memory[c.ip])
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
	fmt.Printf("%d", c.pop())
}

// pop a, pop b, if a != 0 goto b
func (c *cpu) ijnz() {
	a := c.pop()
	b := c.pop()
	if a != 0 {
		c.ip = int(b) + MEM_SIZE/2
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
