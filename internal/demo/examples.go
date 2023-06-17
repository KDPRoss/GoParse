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
	"fmt"
	"strconv"

	"github.com/kdpross/GoParse/pkg/data"
	"github.com/kdpross/GoParse/pkg/parse"
	"github.com/kdpross/GoParse/pkg/parse_ext"
)

func main() {
	fmt.Println("GoParse Demonstration")

	fooParse := parse.Txt("foo")
	if r := parse.Parse(fooParse, "foo"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %q\n", v) // prints `v = "foo"`
	}
	if r := parse.Parse(fooParse, "bar"); r.FailureQ() {
		fmt.Println("Ruh roh! Parse failure.")
	}

	baaarParse := parse.Regexp("ba+r")
	if r := parse.Parse(baaarParse, "baaaaaaaaaar"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %q\n", v) // prints `v = "baaaaaaaaaar"`
	}

	fooBarParse := parse.Seq(fooParse, baaarParse)
	if r := parse.Parse(fooBarParse, "foobaaaaaaaaar"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = (%q, %v)\n", v.First(), v.Second()) // prints `v = ("foo", baaaaaaaaar)`
	}

	num := parse.Proc(
		parse.Regexp("[0-9]+"),
		func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		},
	)
	if r := parse.Parse(num, "1234"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %d\n", v) // prints `v = 1234`
	}

	numList := parseext.RepSep(num, parse.Txt(","))
	if r := parse.Parse(numList, "1,2,12,57"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %v\n", v) // prints `v = [1 2 12 57]`
	}

	numListSpaces := parseext.RepSep1(
		num,
		parse.Seq(
			parse.Txt(","),
			parseext.Spaces,
		),
	)
	if r := parse.Parse(numListSpaces, "1,2,  3,   4"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %v\n", v) // prints `v = [1 2 3 4]`
	}

	var numListCustom data.Lazy[parse.Parser[[]int]]
	numListCustom = data.MkLazy(func() parse.Parser[[]int] {
		cons := func(p data.Pair[int, []int]) []int {
			return append([]int{p.First()}, p.Second()...)
		}
		emptyList := parse.ParserJust([]int{})
		nonemptyTail := parse.Proc(
			parse.Seq(
				parse.SeqRight(parse.Txt(","), num),
				parse.Cache(numListCustom),
			),
			cons,
		)
		nonemptyList := parse.Proc(parse.Seq(num, parse.Alt(nonemptyTail, emptyList)), cons)

		return parse.Alt(nonemptyList, emptyList)
	})
	if r := parse.Parse(parse.Cache(numListCustom), "1,23,456"); r.SuccessQ() {
		v, _ := r.GetSuccess()
		fmt.Printf("v = %v\n", v) // prints `v = [1 23 456]`
	}
}
