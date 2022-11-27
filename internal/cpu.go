package internal

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	MemSize    int = 200
	StackLimit int = 16
)

type handler func()

type cpu struct {
	stack   []uint16
	memory  []uint16
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
	c.initsp()
	c.initip()
	c.inithmap()
	c.initrunning()
}

func (c *cpu) initstack() {
	c.stack = make([]uint16, StackLimit)
}

func (c *cpu) initmem() {
	c.memory = make([]uint16, MemSize)
}

func (c *cpu) initsp() {
	c.sp = 0
}

func (c *cpu) initip() {
	c.ip = MemSize / 2
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
		MUL:    c.imul,
	}
}

func (c *cpu) initrunning() {
	c.running = true
}

func WithMemProg(memory []uint16, program []uint16) *cpu {
	ret := newcpu()
	copy(ret.memory, memory)
	for i, v := range program {
		ret.memory[i+MemSize/2] = v
	}
	return ret
}

func (c *cpu) MemDump() []uint16 {
	dump := make([]uint16, MemSize)
	copy(dump, c.memory)
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

type RunConfig struct {
	Pause   int
	Verbose bool
}

type RunOpt func(*RunConfig)

func WithPause(pause int) RunOpt {
	return func(rc *RunConfig) {
		rc.Pause = pause
	}
}

func WithVerbose() RunOpt {
	return func(rc *RunConfig) {
		rc.Verbose = true
	}
}

func (c *cpu) Run(opts ...RunOpt) {
	config := &RunConfig{}

	for _, o := range opts {
		o(config)
	}

	for c.running {
		c.tick()

		if config.Pause > 0 {
			time.Sleep(time.Duration(int(time.Millisecond) * config.Pause))
		}

		if config.Verbose {
			PrintDump(c)
		}
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

func (c *cpu) fetch() uint16 {
	cmd := c.memory[c.ip]
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

func (c *cpu) push(n uint16) {
	if c.sp == StackLimit {
		panic("stack overflow")
	}
	c.stack[c.sp] = n
	c.sp++
}

func (c *cpu) pop() uint16 {
	if c.sp == 0 {
		panic("stack underflow")
	}
	c.stack[c.sp] = 0
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
	c.ip = int(c.pop()) + MemSize/2
}

// pop a, pop b, if a == 0 goto b
func (c *cpu) ijz() {
	a := c.pop()
	b := c.pop()
	if a == 0 {
		c.ip = int(b) + MemSize/2
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
	fmt.Printf("%d\n", c.pop())
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
