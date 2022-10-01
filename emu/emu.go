package emu

import (
	"fmt"
	"os"
)

const (
	STACK_LIMIT int = 4000
	MEM_SIZE    int = 4000
)

type Cpu struct {
	stack  [STACK_LIMIT]int32
	sp     int
	memory [MEM_SIZE]uint32
	ip     int
}

func New() Cpu {
	instance := Cpu{}
	return instance
}

func WithProgram(program [MEM_SIZE]uint32) Cpu {
	return Cpu{
		memory: program,
	}
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

// pushes given value to a stack
// sp is incremented after new value is set
// will panic if stack is full
func (c *Cpu) push(n int32) {
	if c.sp == STACK_LIMIT {
		panic("stack overflow")
	}
	c.stack[c.sp] = n
	c.sp++
}

// pops top value from stack and returns it
// sp is decremented before value is popped
// will panic if stack is empty
func (c *Cpu) pop() int32 {
	if c.sp == 0 {
		panic("stack underflow")
	}
	c.sp--
	return c.stack[c.sp]
}

// does nothing
func (c *Cpu) inop() {
	return
}

// adds two top elements
// +------    +------
// | 1 2   -> | 3
// +------    +------
func (c *Cpu) iadd() {
	c.push(c.pop() + c.pop())
}

// substracts top value from next top
// !!! may not be desired behaviour !!!
// +------    +------
// | 2 1   -> | 1
// +------    +------
func (c *Cpu) isub() {
	t := c.pop()
	nt := c.pop()
	c.push(nt - t)
}

// bitwise and
// +------    +------
// | 7 5   -> | 5
// +------    +------
func (c *Cpu) iand() {
	c.push(c.pop() & c.pop())
}

// bitwise or
// +------    +------
// | 7 5   -> | 7
// +------    +------
func (c *Cpu) ior() {
	c.push(c.pop() | c.pop())
}

// bitwise xor
// +------    +------
// | 7 5   -> | 2
// +------    +------
func (c *Cpu) ixor() {
	c.push(c.pop() ^ c.pop())
}

// bitwise not
// +------    +------
// | 7     -> | something, idk
// +------    +------
func (c *Cpu) inot() {
	c.push(^c.pop())
}

// read one byte from stdin and push to the stack
func (c *Cpu) iin() {
	buffer := make([]byte, 1)
	if n, err := os.Stdin.Read(buffer); err != nil || n != 1 {
		panic(err)
	}
	c.push(int32(buffer[0]))
}

// write top of the stack into stdout
func (c *Cpu) iout() {
	fmt.Printf("%d\n", c.pop())
}

func (c *Cpu) iload() {
	addr := c.pop()
	c.push(int32(c.memory[addr]))
}

func (c *Cpu) istor() {
	addr := c.pop()
	c.memory[addr] = uint32(c.pop())
}

func (c *Cpu) ijmp() {
	c.ip = int(c.pop())
}

func (c *Cpu) ijz() {
	a := c.pop()
	b := c.pop()
	if a == 0 {
		c.ip = int(b)
	}
}

func (c *Cpu) ipush() {
	c.push(int32(c.memory[c.ip]))
	c.ip++
}

func (c *Cpu) idup() {
	val := c.pop()
	c.push(val)
	c.push(val)
}

func (c *Cpu) iswap() {
	a := c.pop()
	b := c.pop()
	c.push(a)
	c.push(b)
}

func (c *Cpu) irol3() {
	a := c.pop()
	b := c.pop()
	cc := c.pop()
	c.push(a)
	c.push(cc)
	c.push(b)
}

func (c *Cpu) ioutnum() {
	fmt.Printf("%d", c.pop())
}

func (c *Cpu) ijnz() {
	a := c.pop()
	b := c.pop()
	if a != 0 {
		c.ip = int(b)
	}
}

func (c *Cpu) idrop() {
	c.pop()
}

func (c *Cpu) icomp() {
	c.push(-c.pop())
}
