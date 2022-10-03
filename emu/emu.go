package emu

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MEM_SIZE    int = 200
	STACK_LIMIT int = 20
)

type Cpu struct {
	stack  []uint32
	memory []uint32
	sp     int
	ip     int
}

func WithMem(memory []uint32, program []uint32) Cpu {
	instance := Cpu{
		stack: make([]uint32, STACK_LIMIT),
		ip:    MEM_SIZE / 2,
	}

	instance.memory = make([]uint32, MEM_SIZE)
	copy(instance.memory, memory)
	for i, v := range program {
		instance.memory[i+MEM_SIZE/2] = v
	}

	return instance
}

func (c *Cpu) MemDump() []uint32 {
	dump := make([]uint32, MEM_SIZE)
	copy(dump, c.memory)
	return dump
}

func (c *Cpu) Tick() {
	fetched := c.fetch()
	decoded := c.decode(fetched)
	c.execute(decoded)
}

func (c *Cpu) fetch() uint32 {
	cmd := c.memory[c.ip]
	c.ip++
	return cmd
}

func (c *Cpu) decode(opcode uint32) uint32 {
	return opcode
}

func (c *Cpu) execute(cmd uint32) {
	switch cmd {
	case NOP:
		c.inop()
	case ADD:
		c.iadd()
	case SUB:
		c.isub()
	case AND:
		c.iand()
	case OR:
		c.ior()
	case XOR:
		c.ixor()
	case NOT:
		c.inot()
	case IN:
		c.iin()
	case OUT:
		c.iout()
	case LOAD:
		c.iload()
	case STOR:
		c.istor()
	case JMP:
		c.ijmp()
	case JZ:
		c.ijz()
	case PUSH:
		c.ipush()
	case DUP:
		c.idup()
	case SWAP:
		c.iswap()
	case ROL3:
		c.irol3()
	case OUTNUM:
		c.ioutnum()
	case JNZ:
		c.ijnz()
	case DROP:
		c.idrop()
	case COMPL:
		c.icomp()
	default:
		panic("unknown command")
	}
}

func (c *Cpu) push(n uint32) {
	if c.sp == STACK_LIMIT {
		panic("stack overflow")
	}
	c.stack[c.sp] = n
	c.sp++
}

func (c *Cpu) pop() uint32 {
	if c.sp == 0 {
		panic("stack underflow")
	}
	c.sp--
	return c.stack[c.sp]
}

// do nothing
func (c *Cpu) inop() {}

// pop a, pop b, push a + b
func (c *Cpu) iadd() {
	c.push(c.pop() + c.pop())
}

// pop a, pop b, push b - a
// !!! may not be desired behaviour !!!
func (c *Cpu) isub() {
	t := c.pop()
	nt := c.pop()
	c.push(nt - t)
}

// pop a, pop b, push a & b
func (c *Cpu) iand() {
	a := c.pop()
	b := c.pop()
	c.push(a & b)
}

// pop a, pop b, push a | b
func (c *Cpu) ior() {
	a := c.pop()
	b := c.pop()
	c.push(a | b)
}

// pop a, pop b, push a ^ b
func (c *Cpu) ixor() {
	a := c.pop()
	b := c.pop()
	c.push(a ^ b)
}

// pop a, push !a
func (c *Cpu) inot() {
	c.push(^c.pop())
}

// read one byte from stdin and push to the stack
func (c *Cpu) iin() {
	in := bufio.NewReader(os.Stdin)
	b, err := in.ReadByte()
	if err != nil {
		panic(err)
	}
	c.push(uint32(b))
}

// write top of the stack into stdout
func (c *Cpu) iout() {
	out := bufio.NewWriter(os.Stdout)
	if err := out.WriteByte(byte(c.pop())); err != nil {
		panic(err)
	}
}

// pop a, push word read from memory[a]
func (c *Cpu) iload() {
	a := c.pop()
	c.push(c.memory[a])
}

// pop a, pop b, write b to memory[a]
func (c *Cpu) istor() {
	a := c.pop()
	b := c.pop()
	c.memory[a] = b
}

// pop a, goto a
func (c *Cpu) ijmp() {
	c.ip = int(c.pop()) + MEM_SIZE/2
}

// pop a, pop b, if a == 0 goto b
func (c *Cpu) ijz() {
	a := c.pop()
	b := c.pop()
	if a == 0 {
		c.ip = int(b) + MEM_SIZE/2
	}
}

// push next word
func (c *Cpu) ipush() {
	c.push(c.memory[c.ip])
	c.ip++
}

// duplicate stack top
func (c *Cpu) idup() {
	val := c.pop()
	c.push(val)
	c.push(val)
}

// swap two top values
func (c *Cpu) iswap() {
	a := c.pop()
	b := c.pop()
	c.push(a)
	c.push(b)
}

// (a, b, c) -> (b, c, a)
func (c *Cpu) irol3() {
	cc := c.pop()
	b := c.pop()
	a := c.pop()
	c.push(b)
	c.push(cc)
	c.push(a)
}

// write stack top into stdin as number
func (c *Cpu) ioutnum() {
	fmt.Printf("%d", c.pop())
}

// pop a, pop b, if a != 0 goto b
func (c *Cpu) ijnz() {
	a := c.pop()
	b := c.pop()
	if a != 0 {
		c.ip = int(b) + MEM_SIZE/2
	}
}

// pop stack top
func (c *Cpu) idrop() {
	c.pop()
}

// push stack top complement
func (c *Cpu) icomp() {
	c.push(-c.pop())
}
