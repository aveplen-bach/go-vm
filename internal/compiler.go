package internal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	_INSTRUCTION = iota + 1
	_INTEGER
	_LABEL
	_LABELREF
	_COMMENT
)

var (
	breakpoints = []rune{'\n', ' '}
)

func runein(r rune, coll []rune) bool {
	for _, x := range coll {
		if x == r {
			return true
		}
	}
	return false
}

type runeiter struct {
	in     bufio.Reader
	opened bool
	closed bool
	n      rune
}

func newruneiter(in bufio.Reader) runeiter {
	ret := runeiter{in: in}
	ret.walk()
	ret.opened = true
	return ret
}

func (ri *runeiter) walk() {
	if ri.closed {
		return
	}

	n, _, err := ri.in.ReadRune()
	if err != nil {
		ri.n = n
		ri.closed = true
		return
	}

	if ri.n == unicode.ReplacementChar {
		panic("could not read rune: invalid code")
	}

	ri.n = n
}

func (ri *runeiter) hasnext() bool {
	return ri.opened && !ri.closed
}

func (ri *runeiter) next() rune {
	ret := ri.n
	ri.walk()
	return ret
}

type lexem struct {
	val string
	typ int
}

type runeiterator interface {
	next() rune
	hasnext() bool
}

type lexemiter struct {
	runeiter  runeiterator
	opened    bool
	closed    bool
	exhausted bool
	buf       []rune
}

func newlexemiter(runeiter runeiterator) lexemiter {
	ret := lexemiter{runeiter: runeiter}
	ret.walk()
	ret.opened = true
	return ret
}

func (li *lexemiter) walk() {
	if li.closed {
		return
	}

	if li.exhausted {
		li.closed = true
		return
	}

	li.buf = []rune{}
	for li.runeiter.hasnext() {
		nextrune := li.runeiter.next()

		if runein(nextrune, breakpoints) {
			return
		}

		li.buf = append(li.buf, nextrune)
	}

	li.exhausted = true
}

func (li *lexemiter) hasnext() bool {
	return li.opened && !li.closed
}

func isinteger(word string) bool {
	for _, l := range word {
		if !(('0' <= l) && (l <= '9')) {
			return false
		}
	}
	return true
}

func isinstruction(word string) bool {
	return Sinst(strings.ToLower(word))
}

func islabel(word string) bool {
	for _, l := range word {
		if l == '_' {
			continue
		}

		if l == ':' {
			continue
		}

		if l < 'a' {
			return false
		}

		if l > 'z' {
			return false
		}
	}

	lensatisf := len(word) >= 2

	runes := []rune(word)
	colonend := runes[len(runes)-1] == ':'

	return lensatisf && colonend
}

func islabelref(word string) bool {
	for _, l := range word {
		if l == '_' {
			continue
		}

		if l == '&' {
			continue
		}

		if l < 'a' {
			return false
		}

		if l > 'z' {
			return false
		}
	}

	lensatisf := len(word) >= 2
	amperstart := []rune(word)[0] == '&'

	return lensatisf && amperstart
}

func decode(word string) int {
	if isinteger(word) {
		return _INTEGER
	}

	if isinstruction(word) {
		return _INSTRUCTION
	}

	if islabel(word) {
		return _LABEL
	}

	if islabelref(word) {
		return _LABELREF
	}

	panic("could not decode lexem: unknown format word")
}

func (li *lexemiter) next() lexem {
	val := string(li.buf)
	ret := lexem{
		val: val,
		typ: decode(val),
	}
	li.walk()
	return ret
}

type lexemiterator interface {
	next() lexem
	hasnext() bool
}

type labelref struct {
	at   int
	name string
}

type compiler struct {
	lexit  lexemiterator
	labels map[string]uint32
	lrefq  []labelref
	ino    int
}

func newcompiler(lexit lexemiterator) *compiler {
	return &compiler{
		lexit:  lexit,
		labels: make(map[string]uint32),
		lrefq:  make([]labelref, 0),
		ino:    0,
	}
}

func (c *compiler) compile() []uint32 {
	buf := make([]uint32, 0)

	for c.lexit.hasnext() {
		lexem := c.lexit.next()

		var apnd uint32
		switch lexem.typ {

		case _INSTRUCTION:
			apnd = c.compileinstr(lexem.val)

		case _INTEGER:
			apnd = c.compileint(lexem.val)

		case _LABEL:
			apnd = c.compilelabel(lexem.val)

		case _LABELREF:
			apnd = c.compilelabelref(lexem.val)

		case _COMMENT:
			continue

		default:
			panic("unknown lexem type")
		}

		buf = append(buf, apnd)
		c.ino++
	}

	program := c.resolvelabelrefs(buf)
	return program
}

func (c *compiler) compileinstr(value string) uint32 {
	return uint32(Stoi(strings.ToLower(value)))
}

func (c *compiler) compileint(value string) uint32 {
	asint64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(fmt.Errorf("could not parse int: %w", err))
	}
	return uint32(asint64)
}

func (c *compiler) compilelabel(value string) uint32 {
	raw := value[:len(value)-1]
	if _, ok := c.labels[raw]; ok {
		panic(fmt.Errorf("attempt to overwrite '%s' label", raw))
	}

	c.labels[raw] = uint32(c.ino)
	return NOP
}

func (c *compiler) compilelabelref(value string) uint32 {
	raw := value[1:]
	labref, ok := c.labels[raw]
	if !ok {
		c.lrefq = append(c.lrefq, labelref{
			at:   c.ino,
			name: raw,
		})
	}
	return labref
}

func (c *compiler) resolvelabelrefs(prog []uint32) []uint32 {
	processed := make([]uint32, len(prog))
	copy(processed, prog)

	for _, labelref := range c.lrefq {
		ref, ok := c.labels[labelref.name]
		if !ok {
			panic(fmt.Errorf("could not backref label '%s'", labelref.name))
		}

		processed[labelref.at] = ref
	}

	return processed
}

func Compile(in bufio.Reader) []uint32 {
	rit := newruneiter(in)
	lexit := newfsmlex(&rit)
	comp := newcompiler(lexit)
	return comp.compile()
}
