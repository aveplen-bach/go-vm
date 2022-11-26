package internal

import (
	"reflect"
	"testing"
)

type striter struct {
	buf []rune
	ptr int
}

func (s *striter) next() rune {
	ret := s.buf[s.ptr]
	s.ptr++
	return ret
}

func (s *striter) hasnext() bool {
	return s.ptr < len(s.buf)
}

func Test_fsmlex_next(t *testing.T) {
	tests := []struct {
		name string
		li   *fsmlex
		want lexem
	}{
		// single line comments
		{
			name: "stop single line comment on line break of \\n",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("//abc\nabc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "//abc",
				typ: comment,
			},
		},
		{
			name: "stop single line comment on line break of \\r",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("//abc\r abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "//abc",
				typ: comment,
			},
		},
		{
			name: "stop single line comment on line break of \\r\\n",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("//abc\r\n abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "//abc",
				typ: comment,
			},
		},

		// multi line comments
		{
			name: "consume chunk as multiline comment",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("/* abc \r a \n b \r\n abc */ c")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "/* abc \r a \n b \r\n abc */",
				typ: comment,
			},
		},
		{
			name: "ignore asterisks until slash",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("/* abc \r a \n b \r\n abc * * **** c */")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "/* abc \r a \n b \r\n abc * * **** c */",
				typ: comment,
			},
		},

		// numbers
		{
			name: "consume short integer",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("1 abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "1",
				typ: integer,
			},
		},
		{
			name: "consume long integer",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("1234567890 abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "1234567890",
				typ: integer,
			},
		},
		{
			name: "consume hex integer",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("0x123FB abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "0x123FB",
				typ: integer,
			},
		},
		{
			name: "consume binary integer",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("0b010101011 abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "0b010101011",
				typ: integer,
			},
		},

		// labels
		{
			name: "consume label starting with lowercase",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("abc: abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "abc:",
				typ: label,
			},
		},
		{
			name: "consume label starting with uppercase",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("Abc: abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "Abc:",
				typ: label,
			},
		},
		{
			name: "consume label starting with underscore",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("_abc: abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "_abc:",
				typ: label,
			},
		},

		// label references
		{
			name: "consume label reference starting with lowercase",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("&abc abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "&abc",
				typ: labelreference,
			},
		},
		{
			name: "consume label reference starting with uppercase",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("&Abc abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "&Abc",
				typ: labelreference,
			},
		},
		{
			name: "consume label reference starting with underscore",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("&_abc abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "&_abc",
				typ: labelreference,
			},
		},

		// instructions
		{
			name: "consume instruction",
			li: &fsmlex{
				runeiter: &striter{buf: []rune("abc abc")},
				state:    fsmInitial,
				buf:      []rune{},
			},
			want: lexem{
				val: "abc",
				typ: instruction,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.li.init()
			got := tt.li.next()

			if tt.want != got {
				t.Errorf("unexpected outbox value: '%v', expected : '%v'", tt.li.outbox, tt.want)
			}
		})
	}
}

func lex(src string) []lexem {
	res := make([]lexem, 0)
	iter := newfsmlex(&striter{buf: []rune(src)})
	for iter.hasnext() {
		res = append(res, iter.next())
	}
	return res
}

func mklabel(val string) lexem {
	return lexem{
		val: val,
		typ: label,
	}
}

func mklabelref(val string) lexem {
	return lexem{
		val: val,
		typ: labelreference,
	}
}

func mkinteger(val string) lexem {
	return lexem{
		val: val,
		typ: integer,
	}
}

func mkinstr(val string) lexem {
	return lexem{
		val: val,
		typ: instruction,
	}
}

func mkcomment(val string) lexem {
	return lexem{
		val: val,
		typ: comment,
	}
}

func Test_lex(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want []lexem
	}{
		// instructions
		{
			name: "instructions",
			args: args{
				src: "nop or xor and",
			},
			want: []lexem{
				mkinstr("nop"),
				mkinstr("or"),
				mkinstr("xor"),
				mkinstr("and"),
			},
		},
		{
			name: "instructions with whitespace",
			args: args{
				src: "nop  \t  or \r \n xor\r\nand",
			},
			want: []lexem{
				mkinstr("nop"),
				mkinstr("or"),
				mkinstr("xor"),
				mkinstr("and"),
			},
		},
		{
			name: "instructions with comments",
			args: args{
				src: "nop //or\r\nxor/*com\r\nment*/and",
			},
			want: []lexem{
				mkinstr("nop"),
				mkcomment("//or"),
				mkinstr("xor"),
				mkcomment("/*com\r\nment*/"),
				mkinstr("and"),
			},
		},

		// numbers
		{
			name: "numbers with whitespace",
			args: args{
				src: "123  \t  0b01011\n1 \r\n0xFFFFFF",
			},
			want: []lexem{
				mkinteger("123"),
				mkinteger("0b01011"),
				mkinteger("1"),
				mkinteger("0xFFFFFF"),
			},
		},
		{
			name: "numbers with comments",
			args: args{
				src: "123//comment\n0b01011/*comment*/1//\r\n0xFFFFFF",
			},
			want: []lexem{
				mkinteger("123"),
				mkcomment("//comment"),
				mkinteger("0b01011"),
				mkcomment("/*comment*/"),
				mkinteger("1"),
				mkcomment("//"),
				mkinteger("0xFFFFFF"),
			},
		},

		// labels
		{
			name: "labels no whitespace",
			args: args{
				src: "a:b:c:d:e:",
			},
			want: []lexem{
				mklabel("a:"),
				mklabel("b:"),
				mklabel("c:"),
				mklabel("d:"),
				mklabel("e:"),
			},
		},
		{
			name: "labels with whitespace",
			args: args{
				src: "a: \t b: \r\n c: \n\n\n d: \t\t e:",
			},
			want: []lexem{
				mklabel("a:"),
				mklabel("b:"),
				mklabel("c:"),
				mklabel("d:"),
				mklabel("e:"),
			},
		},
		{
			name: "labels with comments",
			args: args{
				src: "a://\nb:/*comment*/c://a:b:c:\nd:/*com:*/\ne:",
			},
			want: []lexem{
				mklabel("a:"),
				mkcomment("//"),
				mklabel("b:"),
				mkcomment("/*comment*/"),
				mklabel("c:"),
				mkcomment("//a:b:c:"),
				mklabel("d:"),
				mkcomment("/*com:*/"),
				mklabel("e:"),
			},
		},

		// label references
		{
			name: "labels refenreces no whitespace",
			args: args{
				src: "&a&b&c&d&e",
			},
			want: []lexem{
				mklabelref("&a"),
				mklabelref("&b"),
				mklabelref("&c"),
				mklabelref("&d"),
				mklabelref("&e"),
			},
		},
		{
			name: "labels refenreces with whitespace",
			args: args{
				src: "&a \t &b \r\n &c \n\n\n &d \t\t &e",
			},
			want: []lexem{
				mklabelref("&a"),
				mklabelref("&b"),
				mklabelref("&c"),
				mklabelref("&d"),
				mklabelref("&e"),
			},
		},
		{
			name: "labels refenreces with comments",
			args: args{
				src: "&a//\n&b/*comment*/&c//&a&b&c\n&d/*&com*/&e",
			},
			want: []lexem{
				mklabelref("&a"),
				mkcomment("//"),
				mklabelref("&b"),
				mkcomment("/*comment*/"),
				mklabelref("&c"),
				mkcomment("//&a&b&c"),
				mklabelref("&d"),
				mkcomment("/*&com*/"),
				mklabelref("&e"),
			},
		},

		// all together
		{
			name: "all together no whitespace",
			args: args{
				src: "a:0b0101 0xfD&b",
			},
			want: []lexem{
				mklabel("a:"),
				mkinteger("0b0101"),
				mkinteger("0xfD"),
				mklabelref("&b"),
			},
		},
		{
			name: "all together with whitespace",
			args: args{
				src: "a:\t\t0b0101\r\n\r\n0xfD\t\r\n&b",
			},
			want: []lexem{
				mklabel("a:"),
				mkinteger("0b0101"),
				mkinteger("0xfD"),
				mklabelref("&b"),
			},
		},
		{
			name: "all together with comments",
			args: args{
				src: "//\na:/*b:*/0b0101//123\r\n0xfD&b//123",
			},
			want: []lexem{
				mkcomment("//"),
				mklabel("a:"),
				mkcomment("/*b:*/"),
				mkinteger("0b0101"),
				mkcomment("//123"),
				mkinteger("0xfD"),
				mklabelref("&b"),
				mkcomment("//123"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lex(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
