package internal

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func Test_runeiter_walk(t *testing.T) {
	tests := []struct {
		name string
		ri   *runeiter
	}{
		{
			name: "should set the next rune",
			ri: &runeiter{
				in:     *bufio.NewReader(bytes.NewBuffer([]byte("b"))),
				opened: true,
				closed: false,
				n:      rune('a'),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ri.walk()

			if tt.ri.n != rune('b') {
				t.Errorf("did not set the next rune: '%v', expected: '%v'", tt.ri.n, rune('b'))
			}
		})
	}
}

func Test_runeiter_next(t *testing.T) {
	tests := []struct {
		name string
		ri   *runeiter
		want rune
	}{
		{
			name: "should return current and goto next",
			ri: &runeiter{
				in:     *bufio.NewReader(bytes.NewBuffer([]byte("b"))),
				opened: true,
				closed: false,
				n:      rune('a'),
			},
			want: rune('a'),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ri.next(); got != tt.want {
				t.Errorf("runeiter.next() = %v, want %v", got, tt.want)
			}
			if tt.ri.n != rune('b') {
				t.Errorf("did not walk to next rune: '%v', expected: '%v'", tt.ri.n, rune('b'))
			}
		})
	}
}

func Test_lexemiter_walk(t *testing.T) {
	tests := []struct {
		name string
		li   *lexemiter
		want []rune
	}{
		{
			name: "should goto next lexem",
			li: &lexemiter{
				runeiter: &runeiter{
					in:     *bufio.NewReader(bytes.NewBuffer([]byte("bc ddd"))),
					opened: true,
					closed: false,
					n:      rune('a'),
				},
				opened: true,
				closed: false,
				buf:    []rune{},
			},
			want: []rune("abc"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.li.walk()

			if !reflect.DeepEqual(tt.want, tt.li.buf) {
				t.Errorf("did not walk to next lexem: '%v', expected: '%v'", tt.li.buf, tt.want)
			}
		})
	}
}

func Test_lexemiter_next(t *testing.T) {
	tests := []struct {
		name string
		li   *lexemiter
		want lexem
	}{
		{
			name: "should goto next lexem",
			li: &lexemiter{
				buf: []rune("nop"),
				runeiter: &runeiter{
					in:     *bufio.NewReader(bytes.NewBuffer([]byte("bc ddd"))),
					opened: true,
					closed: false,
					n:      rune('a'),
				},
			},
			want: lexem{
				val: "nop",
				typ: instruction,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.li.next(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexemiter.next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fsmlexFromString(src string) *fsmlex {
	return newfsmlex(&striter{buf: []rune(src)})
}

func Test_compiler_compile(t *testing.T) {
	tests := []struct {
		name    string
		c       *compiler
		want    []uint16
		wlabels map[string]uint16
	}{
		{
			name: "should compile stream of instructions",
			c: &compiler{
				lexit:  fsmlexFromString("add In jMp NOP"),
				labels: make(map[string]uint16),
			},
			want: []uint16{ADD, IN, JMP, NOP},
		},
		{
			name: "should compile stream of numbers",
			c: &compiler{
				lexit:  fsmlexFromString("1 2 3 123 456"),
				labels: make(map[string]uint16),
			},
			want: []uint16{1, 2, 3, 123, 456},
		},
		{
			name: "should compile stream of labels terminated by instruction",
			c: &compiler{
				lexit:  fsmlexFromString("a: b: c: d: add"),
				labels: make(map[string]uint16),
			},
			want: []uint16{NOP, NOP, NOP, NOP, ADD},
			wlabels: map[string]uint16{
				"a": 0,
				"b": 1,
				"c": 2,
				"d": 3,
			},
		},
		{
			name: "should compile stream of labels not terminated by instruction",
			c: &compiler{
				lexit:  fsmlexFromString("a: b: c: d:"),
				labels: make(map[string]uint16),
			},
			want: []uint16{NOP, NOP, NOP, NOP},
			wlabels: map[string]uint16{
				"a": 0,
				"b": 1,
				"c": 2,
				"d": 3,
			},
		},
		{
			name: "should compile stream of label refs",
			c: &compiler{
				lexit: fsmlexFromString("&a &b &c &d"),
				labels: map[string]uint16{
					"a": 123,
					"b": 456,
					"c": 789,
					"d": 1011,
				},
			},
			want: []uint16{123, 456, 789, 1011},
		},
		{
			name: "should reference label crated before",
			c: &compiler{
				lexit:  fsmlexFromString("add nop a: load &a jmp"),
				labels: make(map[string]uint16),
			},
			want: []uint16{ADD, NOP, NOP, LOAD, 2, JMP},
		},
		{
			name: "should be able to reference stacked lablels",
			c: &compiler{
				lexit:  fsmlexFromString("add nop a: b: load &a &b jmp"),
				labels: make(map[string]uint16),
			},
			want: []uint16{ADD, NOP, NOP, NOP, LOAD, 2, 3, JMP},
		},
		{
			name: "should be able to resolve references before labels",
			c: &compiler{
				lexit:  fsmlexFromString("add nop &a &b load a: b: jmp"),
				labels: make(map[string]uint16),
			},
			want: []uint16{ADD, NOP, 5, 6, LOAD, NOP, NOP, JMP},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.c.compile(false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compiler.compile() = %v, want %v", got, tt.want)
			}

			if tt.wlabels != nil {
				if !reflect.DeepEqual(tt.wlabels, tt.c.labels) {
					t.Errorf("wlabels != c.labels")
				}
			}
		})
	}
}
