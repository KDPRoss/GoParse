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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kdpross/GoParse/pkg/data"
	"github.com/kdpross/GoParse/pkg/parse"
	parseext "github.com/kdpross/GoParse/pkg/parse_ext"
)

// A version of the old 'untyped arithmetic expressions'
// interpreter. In any ordinary PL context, we'd parse this
// to an AST and then interpret *that*, but I've woven the
// interpreter into the parser for brevity.
//
//      n ::= <number>
//      f ::= "(" a ")" | n
//      m ::= f "*" m | f
//      a ::= m "+" a | m

var parser = func() parse.Parser[int] {
	numP := parse.Proc(
		parseext.StringOf(parseext.DigitC),
		func(s string) int {
			i, _ := strconv.Atoi(s)

			return i
		})

	var factor parse.Parser[int]

	var mulP data.Lazy[parse.Parser[int]]

	mulP = data.MkLazy(func() parse.Parser[int] {
		mul := parse.Proc(
			parseext.SeqS(
				parseext.SeqLeftS(
					factor,
					parse.Txt("*"),
				),
				parse.Cache(mulP),
			),
			func(p data.Pair[int, int]) int {
				return p.First() * p.Second()
			},
		)

		return parse.Alt(mul, factor)
	})

	var addP data.Lazy[parse.Parser[int]]

	addP = data.MkLazy(func() parse.Parser[int] {
		add := parse.Proc(
			parseext.SeqS(
				parseext.SeqLeftS(
					parse.Cache(mulP),
					parse.Txt("+"),
				),
				parse.Cache(addP),
			),
			func(p data.Pair[int, int]) int {
				return p.First() + p.Second()
			},
		)

		return parse.Alt(add, parse.Cache(mulP))
	})

	factor = parse.Alt(parseext.Brackets(parse.Cache(addP)), numP)

	return parseext.SeqLeftS(
		parseext.SeqRightS(
			parse.ParserJust(data.Unit{}),
			parse.Cache(addP),
		),
		parse.Eof(),
	)
}()

func main() {
	for _, s := range []string{
		"Welcome to KDP's Wonderful World o' Parsing!",
		"",
		"Enter an expression or 'quit' to exit.",
		"",
	} {
		fmt.Println(s)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(":> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "quit" {
			fmt.Println("\nGoodbye!")

			return
		}

		if r := parse.Parse(parser, line); r.SuccessQ() {
			v, _ := r.GetSuccess()
			fmt.Printf("%d\n\n", v)
		} else {
			fmt.Print("Parse error!\n\n")
		}
	}
}
