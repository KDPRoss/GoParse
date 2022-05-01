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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kdpross/go-parse/internal/data"
	"github.com/kdpross/go-parse/internal/parse"
	"github.com/kdpross/go-parse/internal/parse_ext"
)

func main() {
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

	fmt.Println("Welcome to KDP's Wonderful World o' Parsing!")
	fmt.Println()
	fmt.Println("Enter an expression or 'quit' to exit.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(":> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "quit" {
			fmt.Println("\nGoodbye!")

			return
		}

		r := parse.Parse(
			parseext.SeqLeftS(
				parseext.SeqRightS(
					parse.ParserJust(data.Unit{}),
					parse.Cache(addP),
				),
				parse.Eof(),
			),
			line,
		)

		if r.SuccessQ() {
			v, _ := r.GetSuccess()
			fmt.Printf("%d\n\n", v)
		} else {
			fmt.Print("Parse error!\n\n")
		}
	}
}
