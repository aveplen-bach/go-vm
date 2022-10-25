package internal

import "fmt"

const (
	NOP = iota
	ADD
	SUB
	AND
	OR
	XOR
	NOT
	IN
	OUT
	LOAD
	STOR
	JMP
	JZ
	PUSH
	DUP
	SWAP
	ROL3
	OUTNUM
	JNZ
	DROP
	COMPL
	CINC
	CDEC
	CTS
	STC
	TERM
)

func StoiSafe(name string) (int, error) {
	return stoi(name)
}

func Stoi(name string) int {
	i, err := stoi(name)
	if err != nil {
		panic(err)
	}
	return i
}

func Sinst(name string) bool {
	return sinst(name)
}

func ItosSafe(val int) (string, error) {
	return itos(val)
}

func Itos(val int) string {
	i, err := itos(val)
	if err != nil {
		panic(err)
	}
	return i
}

func Iinst(val int) bool {
	return iinst(val)
}

var mapping = map[string]int{
	"nop":    NOP,
	"add":    ADD,
	"sub":    SUB,
	"and":    AND,
	"or":     OR,
	"xor":    XOR,
	"not":    NOT,
	"in":     IN,
	"out":    OUT,
	"load":   LOAD,
	"stor":   STOR,
	"jmp":    JMP,
	"jz":     JZ,
	"push":   PUSH,
	"dup":    DUP,
	"swap":   SWAP,
	"rol3":   ROL3,
	"outnum": OUTNUM,
	"jnz":    JNZ,
	"drop":   DROP,
	"compl":  COMPL,
	"cinc":   CINC,
	"cdec":   CDEC,
	"cts":    CTS,
	"stc":    STC,
	"term":   TERM,
}

func stoi(name string) (int, error) {
	i, ok := mapping[name]
	if !ok {
		return 0, fmt.Errorf("could not decode instruction '%s'", name)
	}
	return i, nil
}

func sinst(name string) bool {
	_, ok := mapping[name]
	return ok
}

var rmapping = map[int]string{
	NOP:    "nop",
	ADD:    "add",
	SUB:    "sub",
	AND:    "and",
	OR:     "or",
	XOR:    "xor",
	NOT:    "not",
	IN:     "in",
	OUT:    "out",
	LOAD:   "load",
	STOR:   "stor",
	JMP:    "jmp",
	JZ:     "jz",
	PUSH:   "push",
	DUP:    "dup",
	SWAP:   "swap",
	ROL3:   "rol3",
	OUTNUM: "outnum",
	JNZ:    "jnz",
	DROP:   "drop",
	COMPL:  "compl",
	CINC:   "cinc",
	CDEC:   "cdec",
	CTS:    "cts",
	STC:    "stc",
	TERM:   "term",
}

func itos(val int) (string, error) {
	i, ok := rmapping[val]
	if !ok {
		return "", fmt.Errorf("could not decode instruction %#04x", val)
	}

	return i, nil
}

func iinst(val int) bool {
	_, ok := rmapping[val]
	return ok
}
