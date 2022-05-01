// ┌─────────────────────────────────────────────────────────────┐
// │ GoParse: A Golang parser-combinator library.                │
// │                                                             │
// │ This codebase is licensed for the following purposes only:  │
// │                                                             │
// │ - study of the code                                         │
// │                                                             │
// │ - compiling / running an unaltered copy of the code for     │
// │   noncommercial educational and entertainment purposes only │
// │                                                             │
// │ - gratis redistribution of the code in entirety and in      │
// │   unaltered form for any aforementioned purpose             │
// │                                                             │
// │ Copyright 2022, K.D.P.Ross                                  │
// └─────────────────────────────────────────────────────────────┘

package parseext

import (
	"github.com/kdpross/go-parse/pkg/data"
	"github.com/kdpross/go-parse/pkg/parse"
)

func cons[A any](p data.Pair[A, []A]) []A {
	return append([]A{p.First()}, p.Second()...)
}

func Rep1[A any](p parse.Parser[A]) parse.Parser[[]A] {
	return parse.Proc(parse.Seq(p, parse.Rep(p)), cons[A])
}

func RepSep1[A, B any](p parse.Parser[A], s parse.Parser[B]) parse.Parser[[]A] {
	return parse.Proc(parse.Seq(p, parse.Rep(parse.SeqRight(s, p))), cons[A])
}

func RepSep[A, B any](p parse.Parser[A], s parse.Parser[B]) parse.Parser[[]A] {
	return parse.Alt(RepSep1(p, s), parse.ParserJust([]A{}))
}

// This seems to come up as a common pattern: Identifiers,
// etc. have some rule for the first character and different
// rules for subsequent ones.
func IdentOf(fst, rst func(byte) bool) parse.Parser[string] {
	return parse.Proc(
		parse.Seq(parse.OneOf(fst), parse.Rep(parse.OneOf(rst))),
		func(p data.Pair[byte, []byte]) string {
			return string(append([]byte{p.First()}, p.Second()...))
		},
	)
}

func StringOf(p func(byte) bool) parse.Parser[string] {
	return IdentOf(p, p)
}

func spaceQ(b byte) bool {
	return b == ' '
}

var Spaces = parse.Proc(
	parse.Rep(parse.OneOf(spaceQ)),
	func(bs []byte) int {
		return len(bs)
	},
)

var Spaces1 = parse.Proc(
	Rep1(parse.OneOf(spaceQ)),
	func(bs []byte) int {
		return len(bs)
	},
)

// Allow spaces.
func SeqS[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[data.Pair[A, B]] {
	return parse.Seq(parse.SeqLeft(p1, Spaces), p2)
}

func SeqLeftS[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[A] {
	return parse.SeqLeft(parse.SeqLeft(p1, Spaces), p2)
}

func SeqRightS[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[B] {
	return parse.SeqRight(p1, parse.SeqRight(Spaces, p2))
}

// Require spaces.
func SeqS1[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[data.Pair[A, B]] {
	return parse.Seq(parse.SeqLeft(p1, Spaces1), p2)
}

func SeqLeftS1[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[A] {
	return parse.SeqLeft(parse.SeqLeft(p1, Spaces1), p2)
}

func SeqRightS1[A, B any](p1 parse.Parser[A], p2 parse.Parser[B]) parse.Parser[B] {
	return parse.SeqRight(p1, parse.SeqRight(Spaces1, p2))
}

func Maybe[A any](p parse.Parser[A]) parse.Parser[data.Maybe[A]] {
	return parse.Alt(
		parse.Proc(
			p,
			func(x A) data.Maybe[A] {
				return data.MkJust(x)
			},
		),
		parse.ParserJust[data.Maybe[A]](data.MkNothing[A]()),
	)
}

func RangeC(l, h byte) func(byte) bool {
	return func(c byte) bool {
		return c >= l && c <= h
	}
}

func OneOfC(s string) func(byte) bool {
	m := map[byte]bool{}

	for _, c := range s {
		m[byte(c)] = true
	}

	return func(c byte) bool {
		return m[c]
	}
}

func Brackets[A any](p parse.Parser[A]) parse.Parser[A] {
	return parse.SeqLeft(parse.SeqRight(parse.Txt("("), p), parse.Txt(")"))
}

var UpperC = RangeC('A', 'Z')
var LowerC = RangeC('a', 'z')
var DigitC = RangeC('0', '9')
var VarP = StringOf(LowerC)
