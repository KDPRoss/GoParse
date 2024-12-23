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
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kdpross/GoParse/pkg/data"
)

func TestResultSuccessCoherent(t *testing.T) {
	v := 5
	r := success[int]{
		v: v,
	}

	assert.False(t, r.FailureQ())
	require.True(t, r.SuccessQ())
	vP, _ := r.GetSuccess()
	assert.Equal(t, v, vP)
}

func TestResultsFailureCoherent(t *testing.T) {
	r := failure[int]{}

	assert.False(t, r.SuccessQ())
	assert.True(t, r.FailureQ())
}

// return a >>= h ≡ h a
func TestMonadLeftIdentity(t *testing.T) {
	a := 5
	h := func(n int) M[int] {
		return Return(n + 5)
	}

	assert.Equal(t, Bind(Return(a), h).f(0), h(a).f(0))
}

// m >>= return ≡ m
func TestMonadRightIdentity(t *testing.T) {
	m := Return(5)

	assert.Equal(t, Bind(m, Return[int]).f(0), m.f(0))
}

// (m >>= g) >>= h ≡ m >>= (λx . g x >>= h)
func TestMonadAssociativity(t *testing.T) {
	m := Return(5)
	g := func(n int) M[string] {
		return Return(strconv.Itoa(n))
	}
	h := func(s string) M[int] {
		i, _ := strconv.Atoi(s)

		return Return(i)
	}

	assert.Equal(t, Bind(Bind(m, g), h).f(0), Bind(m, func(x int) M[int] { return Bind(g(x), h) }).f(0))
}

func TestSetStGetStCoherent(t *testing.T) {
	n := 5
	m := Bind(
		setSt(n),
		func(data.Unit) M[int] {
			return getSt()
		},
	)

	assert.Equal(t, m.f(0), Return(n).f(n))
}

func TestReturnSucceeds(t *testing.T) {
	n := 5
	r := Return(n).f(0)

	require.True(t, r.SuccessQ())
	nP, _ := r.GetSuccess()
	assert.Equal(t, n, nP)
}

func TestFailFails(t *testing.T) {
	assert.True(t, fail[int]().f(0).FailureQ())
}

func TestOneOf(t *testing.T) {
	baddies := []string{
		"A",
		"Foozle",
		"5",
		"__",
		"",
	}

	for _, r := range []struct {
		lab string
		f   func(byte) bool
		ss  []string
	}{
		{"reject everything", func(byte) bool { return false }, []string{}},
		{"lowercase", func(b byte) bool { return b >= 'a' && b <= 'z' }, []string{}},
		{"Q", func(b byte) bool { return b == 'Q' }, []string{}},
	} {
		t.Run(fmt.Sprintf("test 'OneOf' '%s'", r.lab), func(t *testing.T) {
			p := OneOf(r.f)

			for _, s := range r.ss {
				r := Parse(OneOf(r.f), s)

				require.True(t, r.SuccessQ())
				assert.Equal(t, success[byte]{s[0], 1}, r)
			}

			for _, s := range baddies {
				assert.True(t, Parse(p, s).FailureQ())
			}
		})
	}
}

func TestChr(t *testing.T) {
	baddies := []string{
		"A",
		"Foozle",
		"5",
		"__",
		"",
	}

	for _, c := range []byte("aQz*.") {
		t.Run(fmt.Sprintf("test 'Chr' '%v'", c), func(t *testing.T) {
			p := Chr(c)

			for _, s := range baddies {
				sP := string([]byte{c}) + s
				r := Parse(p, sP)

				require.True(t, r.SuccessQ())
				assert.Equal(t, success[byte]{c, 1}, r)
			}

			for _, s := range baddies {
				assert.True(t, Parse(p, s).FailureQ())
			}
		})
	}
}

func TestNoneOf(t *testing.T) {
	goodies := []string{
		"A",
		"Foozle",
		"5",
		"__",
	}

	for _, c := range []byte("aQz*.") {
		t.Run(fmt.Sprintf("test 'NoneOf' '%v'", c), func(t *testing.T) {
			p := NoneOf(func(cP byte) bool {
				return cP == c
			})

			for _, s := range goodies {
				sP := string([]byte{c}) + s

				require.True(t, Parse(p, sP).FailureQ())
			}

			for _, s := range goodies {
				r := Parse(p, s)

				require.True(t, r.SuccessQ())
				assert.Equal(t, success[byte]{s[0], 1}, r)
			}

			assert.True(t, Parse(p, "").FailureQ())
		})
	}
}

func TestTxt(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	for _, s := range ss {
		t.Run(fmt.Sprintf("test 'Txt' '%s'", s), func(t *testing.T) {
			p := Txt(s)

			for _, sP := range ss {
				r := Parse(p, sP)

				if strings.HasPrefix(sP, s) {
					require.True(t, r.SuccessQ())
					assert.Equal(t, success[string]{s, len(s)}, r)
				} else {
					assert.True(t, r.FailureQ())
				}
			}
		})
	}
}

func TestRegexp(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	t.Run("wildcard matches everything", func(t *testing.T) {
		p := Regexp(".*")

		for _, s := range ss {
			r := Parse(p, s)

			require.True(t, r.SuccessQ())

			sP, _ := r.GetSuccess()

			assert.Equal(t, s, sP)
		}
	})

	t.Run("5-7 chars matches correct num chars", func(t *testing.T) {
		p := Regexp(".{5,7}")

		for _, s := range ss {
			r := Parse(p, s)

			if len(s) >= 5 {
				require.True(t, r.SuccessQ())

				sP, _ := r.GetSuccess()

				fmt.Printf("reached true branch with s=%q;sP=%q;lenprop=%v\n", s, sP, len(sP) >= 5 && len(sP) <= 7)

				assert.True(t, len(sP) >= 5 && len(sP) <= 7)
				assert.True(t, strings.HasPrefix(s, sP))
			} else {
				assert.True(t, r.FailureQ())
			}
		}
	})

	t.Run("bop time", func(t *testing.T) {
		p := Regexp(".*bop")

		for _, s := range ss {
			r := Parse(p, s)

			if strings.Contains(s, "bop") {
				require.True(t, r.SuccessQ())

				sP, _ := r.GetSuccess()

				assert.True(t, strings.HasSuffix(sP, "bop"))
			} else {
				require.True(t, r.FailureQ())
			}
		}
	})

	t.Run("'oo' in-context match", func(t *testing.T) {
		p := Regexp(".*oo.*")

		for _, s := range ss {
			r := Parse(p, s)

			if strings.Contains(s, "oo") {
				require.True(t, r.SuccessQ())

				sP, _ := r.GetSuccess()

				assert.Equal(t, s, sP)
			} else {
				assert.True(t, r.FailureQ())
			}
		}
	})

	t.Run("doesn't match digits", func(t *testing.T) {
		p := Regexp(".*[0-9].*")

		for _, s := range ss {
			assert.True(t, Parse(p, s).FailureQ())
		}
	})
}

func TestRegexpRegression(t *testing.T) {
	p := Regexp("[ab]+")

	r1 := Parse(p, "fooabab")
	assert.True(t, r1.FailureQ())

	r2 := Parse(p, "aaabbbfooba")
	require.True(t, r2.SuccessQ())
	s, _ := r2.GetSuccess()
	assert.Equal(t, "aaabbb", s)
}

func TestSeq(t *testing.T) {
	s := "foozlebopper"

	for i := 0; i < len(s); i++ {
		s1 := s[0:i]
		s2 := s[i:]

		p := Proc(
			Seq(Txt(s1), Txt(s2)),
			func(p data.Pair[string, string]) string {
				return p.First() + p.Second()
			},
		)

		r := Parse(p, s)

		require.True(t, r.SuccessQ())

		sP, _ := r.GetSuccess()

		assert.Equal(t, s, sP)
	}
}

func TestAlt(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	t.Run("left", func(t *testing.T) {
		p := Txt(ss[0])
		for _, s := range ss[1:] {
			p = Alt(p, Txt(s))
		}

		for _, s := range ss {
			r := Parse(p, s)

			require.True(t, r.SuccessQ())

			exp := func() string {
				for _, sP := range ss {
					if strings.HasPrefix(s, sP) {
						return sP
					}
				}

				return s
			}()

			sP, _ := r.GetSuccess()

			assert.Equal(t, exp, sP)
		}
	})

	t.Run("right", func(t *testing.T) {
		p := Txt(ss[len(ss)-1])
		for i := len(ss) - 2; i >= 0; i-- {
			p = Alt(p, Txt(ss[i]))

		}

		for _, s := range ss {
			r := Parse(p, s)

			require.True(t, r.SuccessQ())

			exp := func() string {
				for i := len(ss) - 1; i >= 0; i-- {
					sP := ss[i]

					if strings.HasPrefix(s, sP) {
						return sP
					}
				}

				return s
			}()

			sP, _ := r.GetSuccess()

			assert.Equal(t, exp, sP)
		}
	})
}

func TestParserJust(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	n := 5
	p := ParserJust(n)

	for _, s := range ss {
		r := Parse(p, s)

		require.True(t, r.SuccessQ())

		nP, ix := r.GetSuccess()
		assert.Zero(t, ix)
		assert.Equal(t, n, nP)
	}
}

func TestParserFail(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	p := ParserFail[int]("whatever")

	for _, s := range ss {
		assert.True(t, Parse(p, s).FailureQ())
	}
}

// Also tests `Eof`.
func TestCache(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	// Probably the most-complicated way of consuming the
	// entire string. Mostly, we just care that this
	// terminates.
	var pSilly data.Lazy[Parser[string]]
	pSilly = data.MkLazy(func() Parser[string] {
		return Proc(
			Seq(
				Regexp("."),
				Alt(
					Proc(
						Eof(),
						func(data.Unit) string {
							return ""
						},
					),
					Cache(pSilly),
				),
			),
			func(p data.Pair[string, string]) string {
				return p.First() + p.Second()
			},
		)
	})

	for _, s := range ss {
		r := Parse(Cache(pSilly), s)

		require.True(t, r.SuccessQ())

		sP, _ := r.GetSuccess()

		assert.Equal(t, sP, s)
	}
}

func TestRep(t *testing.T) {
	ss := []string{
		"foo",
		"foozle",
		"foozlebopper",
		"foozleboppersnoozer",
		"bar",
		"baz",
		"bazzle",
	}

	p := Proc(
		Rep(Regexp(".")),
		func(ss []string) string {
			return strings.Join(ss, "")
		},
	)

	for _, s := range ss {
		r := Parse(p, s)

		require.True(t, r.SuccessQ())

		sP, _ := r.GetSuccess()

		assert.Equal(t, s, sP)
	}
}

func TestSeqLeft(t *testing.T) {
	s := "foozlebopper"

	for i := 0; i < len(s); i++ {
		s1 := s[0:i]
		s2 := s[i:]

		p := Parse(SeqLeft(Txt(s1), Txt(s2)), s)

		require.True(t, p.SuccessQ())

		sP, _ := p.GetSuccess()

		assert.Equal(t, s1, sP)
	}
}

func TestSeqRight(t *testing.T) {
	s := "foozlebopper"

	for i := 0; i < len(s); i++ {
		s1 := s[0:i]
		s2 := s[i:]

		p := Parse(SeqRight(Txt(s1), Txt(s2)), s)

		require.True(t, p.SuccessQ())

		sP, _ := p.GetSuccess()

		assert.Equal(t, s2, sP)
	}
}

func TestGuard(t *testing.T) {
	s := "foozleschnozzler"

	p1 := Parse(Guard(Txt(s), func(_ string) bool { return true }), s)
	require.True(t, p1.SuccessQ())
	sP, _ := p1.GetSuccess()
	assert.Equal(t, s, sP)

	p2 := Parse(Guard(Txt(s), func(_ string) bool { return false }), s)
	assert.False(t, p2.SuccessQ())
}

func TestWord(t *testing.T) {
	s := "foo bar"

	p1 := Parse(SeqLeft(Txt("foo"), Eow()), s)
	require.True(t, p1.SuccessQ())

	p2 := Parse(SeqLeft(Txt("fo"), Eow()), s)
	require.False(t, p2.SuccessQ())

	s2 := "foo"
	p3 := Parse(SeqLeft(Txt("foo"), Eow()), s2)
	require.True(t, p3.SuccessQ())
}
