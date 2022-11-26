package internal

import "strings"

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

type lexemiter struct {
	runeiter  runeiterator
	opened    bool
	closed    bool
	exhausted bool
	buf       []rune
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
		return integer
	}

	if isinstruction(word) {
		return instruction
	}

	if islabel(word) {
		return label
	}

	if islabelref(word) {
		return labelreference
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
