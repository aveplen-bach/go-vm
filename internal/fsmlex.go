package internal

import "fmt"

const (
	fsmInitial = iota + 1
	fsmNumber
	fsmHexOrBinNumber
	fsmHexNumber
	fsmBinNumber
	fsmLabelref
	fsmInstruction
	fsmLabel
	fsmComment
	fsmCommentML
	fsmCommentMLClosing
	fsmCommentSL
)

type statehandle func(rune) int

type fsmlex struct {
	runeiter  runeiterator
	state     int
	outbox    lexem
	buf       []rune
	hmap      map[int]statehandle
	ready     bool
	exhausted bool
	closed    bool
}

func (f *fsmlex) init() {
	f.inithmap()
	f.walk()
}

func (f *fsmlex) inithmap() {
	f.hmap = map[int]statehandle{
		fsmInitial:          f.initial,
		fsmNumber:           f.number,
		fsmHexOrBinNumber:   f.hbnumber,
		fsmHexNumber:        f.hnumber,
		fsmBinNumber:        f.bnumber,
		fsmLabelref:         f.labelref,
		fsmInstruction:      f.instr,
		fsmLabel:            f.label,
		fsmComment:          f.comment,
		fsmCommentML:        f.commentml,
		fsmCommentMLClosing: f.commentmlclosing,
		fsmCommentSL:        f.commentsl,
	}
}

func (f *fsmlex) initial(next rune) int {
	if next == '0' {
		f.buf = append(f.buf, next)
		return fsmHexOrBinNumber
	}
	if digit(next) {
		f.buf = append(f.buf, next)
		return fsmNumber
	}

	if next == '&' {
		f.buf = append(f.buf, next)
		return fsmLabelref
	}

	if labstch(next) {
		f.buf = append(f.buf, next)
		return fsmInstruction
	}

	if next == '/' {
		f.buf = append(f.buf, next)
		return fsmComment
	}

	if whch(next) {
		return fsmInitial
	}

	panic(fmt.Errorf("could not determine next state from '%v' at INITIAL", next))
}

func (f *fsmlex) comment(next rune) int {
	if next == '/' {
		f.buf = append(f.buf, next)
		return fsmCommentSL
	}
	if next == '*' {
		f.buf = append(f.buf, next)
		return fsmCommentML
	}
	panic(fmt.Errorf("could not determine next state from '%v' at COMMENT", next))
}

func (f *fsmlex) commentml(next rune) int {
	if next == '*' {
		f.buf = append(f.buf, next)
		return fsmCommentMLClosing
	}
	f.buf = append(f.buf, next)
	return fsmCommentML
}

func (f *fsmlex) commentmlclosing(next rune) int {
	if next == '/' {
		f.buf = append(f.buf, next)
		f.yield()
		return fsmInitial
	}
	if next == '*' {
		f.buf = append(f.buf, next)
		return fsmCommentMLClosing
	}
	f.buf = append(f.buf, next)
	return fsmCommentML
}

func (f *fsmlex) commentsl(next rune) int {
	if next == '\n' || next == '\r' {
		f.yield()
		return fsmInitial
	}
	f.buf = append(f.buf, next)
	return fsmCommentSL
}

func (f *fsmlex) number(next rune) int {
	if digit(next) || next == 'b' || next == 'x' {
		f.buf = append(f.buf, next)
		return fsmNumber
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	panic(fmt.Errorf("illegal character inside number: '%s'", string(next)))
}

func (f *fsmlex) hbnumber(next rune) int {
	if next == 'b' {
		f.buf = append(f.buf, next)
		return fsmBinNumber
	}
	if next == 'x' {
		f.buf = append(f.buf, next)
		return fsmHexNumber
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}

	panic(fmt.Errorf("illegal number format: '%s'", string(next)))
}

func (f *fsmlex) hnumber(next rune) int {
	if hexdig(next) {
		f.buf = append(f.buf, next)
		return fsmHexNumber
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	panic(fmt.Errorf("illegal hex digit: '%s'", string(next)))
}

func (f *fsmlex) bnumber(next rune) int {
	if bindig(next) {
		f.buf = append(f.buf, next)
		return fsmBinNumber
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	panic(fmt.Errorf("illegal binary digit: '%s'", string(next)))
}

func (f *fsmlex) labelref(next rune) int {
	if labstch(next) || digit(next) {
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	panic(fmt.Errorf("illegal character for labelref: '%v'", next))
}

func whch(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func digit(r rune) bool {
	return '0' <= r && r <= '9'
}

func hexdig(r rune) bool {
	return 'a' <= r && r <= 'f' ||
		'A' <= r && r <= 'F' ||
		digit(r)
}

func bindig(r rune) bool {
	return r == '0' || r == '1'
}

func lowch(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func uppch(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func labstch(r rune) bool {
	return r == '_' || lowch(r) || uppch(r)
}

func labch(r rune) bool {
	return labstch(r) || digit(r)
}

var statetyp = map[int]int{
	fsmNumber:           integer,
	fsmHexOrBinNumber:   integer,
	fsmHexNumber:        integer,
	fsmBinNumber:        integer,
	fsmLabelref:         labelreference,
	fsmInstruction:      instruction,
	fsmLabel:            label,
	fsmCommentMLClosing: comment,
	fsmCommentSL:        comment,
}

func (f *fsmlex) yield() {
	lexem1 := string(f.buf)

	typ, ok := statetyp[f.state]
	if !ok {
		panic(fmt.Errorf("could not decide on token type: '%s'", lexem1))
	}

	f.outbox = lexem{
		val: string(f.buf),
		typ: typ,
	}

	f.buf = []rune{}
	f.ready = true
}

func (f *fsmlex) instr(next rune) int {
	if labch(next) {
		f.buf = append(f.buf, next)
		return fsmInstruction
	}
	if next == ':' {
		f.buf = append(f.buf, next)
		f.state = fsmLabel
		f.yield()
		return fsmInitial
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmLabelref
	}
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return fsmComment
	}
	panic(fmt.Errorf("illegal character for instruction or label: '%s'", string(next)))
}

func (f *fsmlex) label(next rune) int {
	if whch(next) {
		f.yield()
		return fsmInitial
	}
	panic(fmt.Errorf("illegal character for label: '%v'", next))
}

func newfsmlex(runeiter runeiterator) *fsmlex {
	ret := &fsmlex{
		runeiter: runeiter,
		state:    fsmInitial,
	}
	ret.init()
	return ret
}

func (f *fsmlex) walk() {
	if f.exhausted {
		f.closed = true
	}

	for f.runeiter.hasnext() {
		next := f.runeiter.next()

		h, ok := f.hmap[f.state]
		if !ok {
			panic(fmt.Errorf("no handler for current state: %d", f.state))
		}

		f.state = h(next)

		if f.ready {
			return
		}
	}

	// Yield and exhaust to collect remaining
	// unfinished lexem on future next() call.
	if len(f.buf) != 0 {
		f.yield()
		f.exhausted = true
		return
	}

	f.closed = true
}

func (f *fsmlex) hasnext() bool {
	return !f.closed
}

func (f *fsmlex) next() lexem {
	current := f.outbox
	f.ready = false
	f.walk()
	return current
}
