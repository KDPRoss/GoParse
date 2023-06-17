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
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kdpross/GoParse/pkg/parse"
)

func TestRep1(t *testing.T) {
	p := Rep1(parse.Txt("X"))

	for _, s := range []string{
		"",
		"X",
		"XX",
		"X foo",
		"XX foo",
		"ZZZ",
	} {
		r := parse.Parse(p, s)

		if strings.HasPrefix(s, "X") {
			assert.True(t, r.SuccessQ())
		} else {
			assert.True(t, r.FailureQ())
		}
	}
}

func TestRepSep1(t *testing.T) {
	p := RepSep1(parse.Regexp("x+"), parse.Txt(","))

	for _, s := range []string{
		"",
		"x",
		"xxx",
		"x,x,x",
		"x,xx,xxxfoo",
		"moo",
	} {
		r := parse.Parse(p, s)

		if strings.HasPrefix(s, "x") {
			require.True(t, r.SuccessQ())

			ss, _ := r.GetSuccess()

			assert.Equal(t, len(strings.Split(s, ",")), len(ss))
		} else {
			assert.True(t, r.FailureQ())
		}
	}
}

func TestRepSep(t *testing.T) {
	p := RepSep(parse.Regexp("x+"), parse.Txt(","))

	for _, s := range []string{
		"",
		"x",
		"xxx",
		"x,x,x",
		"x,xx,xxxfoo",
		"moo",
	} {
		r := parse.Parse(p, s)

		require.True(t, r.SuccessQ())

		ss, _ := r.GetSuccess()

		exp := func() int {
			if strings.HasPrefix(s, "x") {
				return len(strings.Split(s, ","))
			}

			return 0
		}()

		assert.Equal(t, exp, len(ss))
	}
}

func TestIdentOf(t *testing.T) {
	p1 := IdentOf(
		LowerC,
		func(c byte) bool {
			return LowerC(c) || UpperC(c)
		},
	)

	p2 := parse.Regexp("[a-z][a-zA-Z]*")

	for _, s := range []string{
		"",
		"x",
		"fZ",
		"foo",
		"barZ",
		"camelCase",
		"Abc",
	} {
		r := parse.Parse(parse.SeqLeft(p1, parse.Eof()), s)

		if parse.Parse(parse.Seq(p2, parse.Eof()), s).SuccessQ() {
			require.True(t, r.SuccessQ())

			sP, _ := r.GetSuccess()

			assert.Equal(t, s, sP)
		} else {
			require.True(t, r.FailureQ())
		}
	}
}

func TestStringOf(t *testing.T) {
	p := StringOf(DigitC)

	for _, s := range []string{
		"",
		"abc",
		"123",
		"0",
	} {
		r := parse.Parse(parse.SeqLeft(p, parse.Eof()), s)

		if _, err := strconv.Atoi(s); err == nil {
			require.True(t, r.SuccessQ())

			sP, _ := r.GetSuccess()

			assert.Equal(t, s, sP)
		} else {
			require.True(t, r.FailureQ())
		}
	}
}

func TestSeqSAndSeqS1(t *testing.T) {
	pL := StringOf(LowerC)
	pR := StringOf(DigitC)
	p := SeqS(pL, pR)
	p1 := SeqS1(pL, pR)

	for _, s := range []string{
		"",
		"f1",
		"f 1",
		"f  1",
		"foo123",
		"foo 123",
		"foo  123",
	} {
		r := parse.Parse(p, s)

		if len(s) > 0 && LowerC(s[0]) && DigitC(s[len(s)-1]) {
			require.True(t, r.SuccessQ())
		} else {
			require.True(t, r.FailureQ())
		}

		r1 := parse.Parse(p1, s)

		if strings.Contains(s, " ") && len(s) > 0 && LowerC(s[0]) && DigitC(s[len(s)-1]) {
			require.True(t, r1.SuccessQ())
		} else {
			require.True(t, r1.FailureQ())
		}
	}
}

func TestMaybe(t *testing.T) {
	p := Maybe(StringOf(LowerC))

	for _, s := range []string{
		"",
		"a",
		"abc",
		" abc",
		"5",
	} {
		r := parse.Parse(p, s)

		require.True(t, r.SuccessQ())

		m, _ := r.GetSuccess()

		if len(s) > 0 && LowerC(s[0]) {

			require.True(t, m.JustQ())

			rr := parse.Parse(parse.Regexp("[a-z]+"), s)

			require.True(t, rr.SuccessQ())

			sP, _ := rr.GetSuccess()

			assert.Equal(t, m.GetJust(), sP)
		} else {
			require.True(t, m.NothingQ())
		}
	}
}

func TestBrackets(t *testing.T) {
	pi := StringOf(LowerC)
	p := Brackets(pi)

	for _, s := range []string{
		"foo",
		"(foo)",
		"(foo)",
		"(foo",
		"foo)",
	} {
		r := parse.Parse(p, s)

		if s[0] == '(' && s[len(s)-1] == ')' {
			require.True(t, r.SuccessQ())

			ri := parse.Parse(pi, s[1:len(s)-1])

			require.True(t, ri.SuccessQ())

			sP, _ := r.GetSuccess()
			si, _ := ri.GetSuccess()

			require.Equal(t, sP, si)
		} else {
			assert.True(t, r.FailureQ())
		}
	}

	assert.True(t, parse.Parse(p, "(foo())").FailureQ())
}
