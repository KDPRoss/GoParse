package main

import (
	"fmt"

	"github.com/kdpross/GoParse/pkg/parse"
)

func main() {
	fmt.Println("hello, world")

	if r := parse.Parse(parse.SeqLeft(ParseKind, parse.Eof()), "* -> * -> *"); r.SuccessQ() {
		k, _ := r.GetSuccess()
		fmt.Printf("parsed %s\n", ShowKind(k))
	}
}
