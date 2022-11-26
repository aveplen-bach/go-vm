package internal

import (
	"bufio"
	"unicode"
)

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
