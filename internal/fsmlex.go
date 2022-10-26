package internal

import "fmt"

const (
	_FSM_INITIAL = iota + 1
	_FSM_NUMBER
	_FSM_HEX_OR_BIN_NUMBER
	_FSM_HEX_NUMBER
	_FSM_BIN_NUMBER
	_FSM_LABELREF
	_FSM_INSTR
	_FSM_LABEL
	_FSM_COMMENT
	_FSM_COMMENT_ML
	_FSM_COMMENT_ML_CLOSING
	_FSM_COMMENT_SL
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
		_FSM_INITIAL:            f.initial,
		_FSM_NUMBER:             f.number,
		_FSM_HEX_OR_BIN_NUMBER:  f.hbnumber,
		_FSM_HEX_NUMBER:         f.hnumber,
		_FSM_BIN_NUMBER:         f.bnumber,
		_FSM_LABELREF:           f.labelref,
		_FSM_INSTR:              f.instr,
		_FSM_LABEL:              f.label,
		_FSM_COMMENT:            f.comment,
		_FSM_COMMENT_ML:         f.commentml,
		_FSM_COMMENT_ML_CLOSING: f.commentmlclosing,
		_FSM_COMMENT_SL:         f.commentsl,
	}
}

func (f *fsmlex) initial(next rune) int {
	if next == '0' {
		f.buf = append(f.buf, next)
		return _FSM_HEX_OR_BIN_NUMBER
	}
	if digit(next) {
		f.buf = append(f.buf, next)
		return _FSM_NUMBER
	}

	if next == '&' {
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}

	if labstch(next) {
		f.buf = append(f.buf, next)
		return _FSM_INSTR
	}

	if next == '/' {
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}

	if whch(next) {
		return _FSM_INITIAL
	}

	panic(fmt.Errorf("could not determine next state from '%v' at INITIAL", next))
}

func (f *fsmlex) comment(next rune) int {
	if next == '/' {
		f.buf = append(f.buf, next)
		return _FSM_COMMENT_SL
	}
	if next == '*' {
		f.buf = append(f.buf, next)
		return _FSM_COMMENT_ML
	}
	panic(fmt.Errorf("could not determine next state from '%v' at COMMENT", next))
}

func (f *fsmlex) commentml(next rune) int {
	if next == '*' {
		f.buf = append(f.buf, next)
		return _FSM_COMMENT_ML_CLOSING
	}
	f.buf = append(f.buf, next)
	return _FSM_COMMENT_ML
}

func (f *fsmlex) commentmlclosing(next rune) int {
	if next == '/' {
		f.buf = append(f.buf, next)
		f.yield()
		return _FSM_INITIAL
	}
	if next == '*' {
		f.buf = append(f.buf, next)
		return _FSM_COMMENT_ML_CLOSING
	}
	f.buf = append(f.buf, next)
	return _FSM_COMMENT_ML
}

func (f *fsmlex) commentsl(next rune) int {
	if next == '\n' || next == '\r' {
		f.yield()
		return _FSM_INITIAL
	}
	f.buf = append(f.buf, next)
	return _FSM_COMMENT_SL
}

func (f *fsmlex) number(next rune) int {
	if digit(next) || next == 'b' || next == 'x' {
		f.buf = append(f.buf, next)
		return _FSM_NUMBER
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}
	panic(fmt.Errorf("illegal character inside number: '%s'", string(next)))
}

func (f *fsmlex) hbnumber(next rune) int {
	if next == 'b' {
		f.buf = append(f.buf, next)
		return _FSM_BIN_NUMBER
	}
	if next == 'x' {
		f.buf = append(f.buf, next)
		return _FSM_HEX_NUMBER
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}

	panic(fmt.Errorf("illegal number format: '%s'", string(next)))
}

func (f *fsmlex) hnumber(next rune) int {
	if hexdig(next) {
		f.buf = append(f.buf, next)
		return _FSM_HEX_NUMBER
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}
	panic(fmt.Errorf("illegal hex digit: '%s'", string(next)))
}

func (f *fsmlex) bnumber(next rune) int {
	if bindig(next) {
		f.buf = append(f.buf, next)
		return _FSM_BIN_NUMBER
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}
	panic(fmt.Errorf("illegal binary digit: '%s'", string(next)))
}

func (f *fsmlex) labelref(next rune) int {
	if labstch(next) || digit(next) {
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
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
	_FSM_NUMBER:             _INTEGER,
	_FSM_HEX_OR_BIN_NUMBER:  _INTEGER,
	_FSM_HEX_NUMBER:         _INTEGER,
	_FSM_BIN_NUMBER:         _INTEGER,
	_FSM_LABELREF:           _LABELREF,
	_FSM_INSTR:              _INSTRUCTION,
	_FSM_LABEL:              _LABEL,
	_FSM_COMMENT_ML_CLOSING: _COMMENT,
	_FSM_COMMENT_SL:         _COMMENT,
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
		return _FSM_INSTR
	}
	if next == ':' {
		f.buf = append(f.buf, next)
		f.state = _FSM_LABEL
		f.yield()
		return _FSM_INITIAL
	}
	if next == '&' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_LABELREF
	}
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}
	if next == '/' {
		f.yield()
		f.buf = append(f.buf, next)
		return _FSM_COMMENT
	}
	panic(fmt.Errorf("illegal character for instruction or label: '%s'", string(next)))
}

func (f *fsmlex) label(next rune) int {
	if whch(next) {
		f.yield()
		return _FSM_INITIAL
	}
	panic(fmt.Errorf("illegal character for label: '%v'", next))
}

func newfsmlex(runeiter runeiterator) *fsmlex {
	ret := &fsmlex{
		runeiter: runeiter,
		state:    _FSM_INITIAL,
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
