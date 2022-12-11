package internal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type lexem struct {
	val string
	typ int
}

type runeiterator interface {
	next() rune
	hasnext() bool
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
	labels map[string]uint16
	lrefq  []labelref
	ino    int
}

func NewCompiler(lexit lexemiterator) *compiler {
	return &compiler{
		lexit:  lexit,
		labels: make(map[string]uint16),
		lrefq:  make([]labelref, 0),
		ino:    0,
	}
}

func (c *compiler) compile(verbose bool) ([]uint16, error) {
	buf := make([]uint16, 0)

	for c.lexit.hasnext() {
		lexem := c.lexit.next()

		var apnd uint16
		switch lexem.typ {

		case instruction:
			if verbose {
				fmt.Printf("{INSTRUCTION %s}\n", lexem.val)
			}
			apnd = c.compileinstr(lexem.val)

		case integer:
			if verbose {
				fmt.Printf("{INTEGER %s}\n", lexem.val)
			}
			apnd = c.compileint(lexem.val)

		case label:
			if verbose {
				fmt.Printf("{LABEL %s}\n", lexem.val)
			}
			apnd = c.compilelabel(lexem.val)

		case labelreference:
			if verbose {
				fmt.Printf("{LABELREF %s}\n", lexem.val)
			}
			apnd = c.compilelabelref(lexem.val)

		case comment:
			if verbose {
				fmt.Printf("{COMMENT %s}\n", lexem.val)
			}
			continue

		default:
			return nil, fmt.Errorf("unknown lexem type")
		}

		buf = append(buf, apnd)
		c.ino++
	}

	program := c.resolvelabelrefs(buf)
	return program, nil
}

func (c *compiler) compileinstr(value string) uint16 {
	return uint16(Stoi(strings.ToLower(value)))
}

func (c *compiler) compileint(value string) uint16 {
	asint64, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		panic(fmt.Errorf("could not parse int: %w", err))
	}
	return uint16(asint64)
}

func (c *compiler) compilelabel(value string) uint16 {
	raw := value[:len(value)-1]
	if _, ok := c.labels[raw]; ok {
		panic(fmt.Errorf("attempt to overwrite '%s' label", raw))
	}

	c.labels[raw] = uint16(c.ino)
	return NOP
}

func (c *compiler) compilelabelref(value string) uint16 {
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

func (c *compiler) resolvelabelrefs(prog []uint16) []uint16 {
	processed := make([]uint16, len(prog))
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

func Compile(in bufio.Reader, verbose bool) ([]uint16, error) {
	rit := newruneiter(in)
	lexit := newfsmlex(&rit)
	comp := NewCompiler(lexit)

	prog, err := comp.compile(verbose)
	if err != nil {
		return nil, fmt.Errorf("compiler.compile(): %w", err)
	}
	return prog, nil
}
