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
// │ Copyright 2022-2025, K.D.P.Ross                             │
// └─────────────────────────────────────────────────────────────┘

package parse

import (
	"github.com/kdpross/GoParse/pkg/data"
)

func SeqRight[A, B any](p1 Parser[A], p2 Parser[B]) Parser[B] {
	return Proc(
		Seq(p1, p2),
		func(p data.Pair[A, B]) B {
			return p.Second()
		},
	)
}

func SeqLeft[A, B any](p1 Parser[A], p2 Parser[B]) Parser[A] {
	return Proc(
		Seq(p1, p2),
		func(p data.Pair[A, B]) A {
			return p.First()
		},
	)
}

// Tie everything together.
func Parse[A any](p Parser[A], s string) Result[A] {
	src := source{
		str:  s,
		undo: []func(){},
	}

	res := p.core(src).f(0)

	for i := len(src.undo) - 1; i >= 0; i-- {
		src.undo[i]()
	}

	return res
}
