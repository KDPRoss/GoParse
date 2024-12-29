package main

import (
	"fmt"
	"testing"

	"github.com/kdpross/GoParse/pkg/parse"
	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	p := parse.SeqLeft(ParseType, parse.Eof())

	for _, c := range []struct{ lab, s1, s2 string }{
		{"var", "x", "x"},
		{"var multichar", "xyz", "xyz"},
		{"tabs", "(a : *) => a", "((a : *) => a)"},
		{"tcval", "Foo a b c", "(Foo a b c)"},
		{"tcval no args", "Foo", "Foo"},
		{"tcval nested", "Foo ((a : *) => a) (Bar b) c", "(Foo ((a : *) => a) (Bar b) c)"},
		{"tarr", "(Foo x) -> (Bar y)", "((Foo x) -> (Bar y))"},
		{"tarr nested", "(Foo x) -> (Bar y) -> (Moo z)", "((Foo x) -> ((Bar y) -> (Moo z)))"},
		{"app", "m a", "(m a)"},
		{"app nested", "m a (b -> c) (Foo)", "(((m a) (b -> c)) Foo)"},
		{"bracketed", "(((((x)))))", "x"},
		{"tuple", "(a, b, c)", "(a, b, c)"},
		{"tuple nested", "(a, b -> c, Foo bar baz, (d, e, f))", "(a, (b -> c), (Foo bar baz), (d, e, f))"},
		{"tarr nested simpl", "Foo x -> Bar y -> Moo z", "((Foo x) -> ((Bar y) -> (Moo z)))"},
		{"tabs nested", "(a : *) => (b : *) => (p -> q) -> a -> Foo b c -> (d, e, f)", "((a : *) => ((b : *) => ((p -> q) -> (a -> ((Foo b c) -> (d, e, f))))))"},
	} {
		t.Run(c.lab, func(t *testing.T) {
			v := parse.Parse(p, c.s1)
			assert.True(t, v.SuccessQ())
			if v.SuccessQ() {
				typ, _ := v.GetSuccess()
				fmt.Printf("%q\n", typ)
				assert.Equal(t, c.s2, ShowType(typ))
			}
		})
	}
}
