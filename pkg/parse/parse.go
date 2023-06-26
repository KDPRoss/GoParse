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

package parse

// Based on 'Packrat Parsing' (B.Ford, 2002)

import (
	"github.com/gijsbers/go-pcre"

	"github.com/kdpross/GoParse/pkg/data"
)

type source struct {
	str  string
	undo []func()
}

type Result[A any] interface {
	SuccessQ() bool
	FailureQ() bool
	GetSuccess() (A, int)
}

type success[A any] struct {
	v  A
	ix int
}

var _ Result[int] = success[int]{1, 2}

func (success[A]) FailureQ() bool {
	return false
}

func (s success[A]) GetSuccess() (A, int) {
	return s.v, s.ix
}

func (success[A]) SuccessQ() bool {
	return true
}

type failure[A any] struct{}

var _ Result[int] = failure[int]{}

func (failure[A]) FailureQ() bool {
	return true
}

func (failure[A]) GetSuccess() (A, int) {
	panic("unimplemented")
}

func (failure[A]) SuccessQ() bool {
	return false
}

type Parser[A any] struct {
	core  func(src source) M[A]
	cache []data.Maybe[Result[A]]
}

// Artificial struct because Golang can't handle polymorphic
// aliases. :sad_panda:
type M[A any] struct {
	f func(int) Result[A]
}

// You're probably doing dodgy things if you need to
// (externally) refer to `Bind` and `Return`, but there are
// times when this is helpful. (Also, this lets us spell
// `Return` correctly. These dopey block-structured
// languages that require explicit `return`...)
func Return[A any](x A) M[A] {
	return M[A]{
		func(ix int) Result[A] {
			return success[A]{x, ix}
		},
	}
}

func Bind[A, B any](f M[A], g func(A) M[B]) M[B] {
	return M[B]{
		func(ix int) Result[B] {
			fSt := f.f(ix)
			if fSt.SuccessQ() {
				v, ix := fSt.GetSuccess()
				return g(v).f(ix)
			}

			return failure[B]{}
		},
	}
}

func getSt() M[int] {
	return M[int]{
		func(ix int) Result[int] {
			return success[int]{ix, ix}
		},
	}
}

func setSt(ix int) M[data.Unit] {
	return M[data.Unit]{
		func(int) Result[data.Unit] {
			return success[data.Unit]{data.Unit{}, ix}
		},
	}
}

func fail[A any]() M[A] {
	return M[A]{
		func(int) Result[A] {
			return failure[A]{}
		},
	}
}

func makeParser[A any](core func(src source) M[A]) Parser[A] {
	return Parser[A]{
		core:  core,
		cache: nil,
	}
}

func OneOf(p func(byte) bool) Parser[byte] {
	return makeParser(
		func(src source) M[byte] {
			return Bind(
				getSt(),
				func(ix int) M[byte] {
					if ix < len(src.str) && p(src.str[ix]) {
						return Bind(
							setSt(ix+1),
							func(data.Unit) M[byte] {
								return Return(src.str[ix])
							},
						)
					}

					return fail[byte]()
				},
			)
		},
	)
}

func Chr(c byte) Parser[byte] {
	return OneOf(func(cP byte) bool {
		return c == cP
	})
}

func NoneOf(p func(byte) bool) Parser[byte] {
	return OneOf(func(c byte) bool {
		return !p(c)
	})
}

func Txt(v string) Parser[string] {
	vLen := len(v)

	return makeParser(
		func(src source) M[string] {
			return Bind(
				getSt(),
				func(ix int) M[string] {
					var loop func(int) M[string]
					loop = func(i int) M[string] {
						if i == vLen {
							return Bind(
								setSt(ix+vLen),
								func(data.Unit) M[string] {
									return Return(v)
								},
							)
						}

						if i < vLen && ix+i < len(src.str) && v[i] == src.str[ix+i] {
							return loop(i + 1)
						}

						return fail[string]()
					}

					return loop(0)
				},
			)
		},
	)
}

// You better believe that we're using PCREs. It's the only
// redeeming bit of Perl.
func Regexp(reg string) Parser[string] {
	r := pcre.MustCompile("^"+reg, 0)
	return makeParser(func(src source) M[string] {
		return Bind(
			getSt(),
			func(ix int) M[string] {
				s := src.str[ix:]
				m := r.MatcherString(s, 0)
				if !m.Matches() {
					return fail[string]()
				}

				sP := m.GroupString(0)
				return Bind(
					setSt(ix+len(sP)),
					func(data.Unit) M[string] {
						return Return(sP)
					},
				)
			},
		)
	})
}

func Seq[A, B any](p1 Parser[A], p2 Parser[B]) Parser[data.Pair[A, B]] {
	return makeParser(
		func(src source) M[data.Pair[A, B]] {
			return Bind(
				p1.core(src),
				func(v1 A) M[data.Pair[A, B]] {
					return Bind(
						p2.core(src),
						func(v2 B) M[data.Pair[A, B]] {
							return Return(data.MkPair(v1, v2))
						},
					)
				},
			)
		},
	)
}

func Alt[A any](p1, p2 Parser[A]) Parser[A] {
	return makeParser(
		func(src source) M[A] {
			return M[A]{
				func(ix int) Result[A] {
					r := p1.core(src).f(ix)
					if r.SuccessQ() {
						return r
					}

					return p2.core(src).f(ix)
				},
			}
		},
	)
}

func Proc[A, B any](p Parser[A], f func(A) B) Parser[B] {
	return makeParser(
		func(src source) M[B] {
			return Bind(
				p.core(src),
				func(v A) M[B] {
					return Return(f(v))
				},
			)
		},
	)
}

func ParserJust[A any](v A) Parser[A] {
	return makeParser(
		func(source) M[A] {
			return Return(v)
		},
	)
}

func ParserFail[A any](string) Parser[A] {
	return makeParser(
		func(source) M[A] {
			return fail[A]()
		},
	)
}

// This is always so much nicer in a lazy language; need
// something to avoid eagerly evaluating, e.g., the RHS of
// an `Alt` to avoid divergence.
func Cache[A any](lz data.Lazy[Parser[A]]) Parser[A] {
	return makeParser(
		func(src source) M[A] {
			return M[A]{
				func(ix int) Result[A] {
					if ix > len(src.str) {
						return failure[A]{}
					}

					p := lz.Force()

					arr := p.cache

					if p.cache == nil {
						arr = make([]data.Maybe[Result[A]], len(src.str))
						for i := 0; i < len(src.str); i++ {
							arr[i] = data.Nothing[Result[A]]{}
						}

						p.cache = arr

						src.undo = append(
							src.undo,
							func() {
								p.cache = nil
							},
						)
					}

					var res Result[A]

					mFl := arr[ix]

					if mFl.JustQ() {
						res = mFl.GetJust()
					} else {
						res = p.core(src).f(ix)

						arr[ix] = data.MkJust(res)
					}

					if res.SuccessQ() {
						v, ix := res.GetSuccess()

						return success[A]{v, ix}
					}

					return failure[A]{}
				},
			}
		},
	)
}

func Rep[A any](p Parser[A]) Parser[[]A] {
	return makeParser(
		func(src source) M[[]A] {
			var loop func([]A) M[[]A]

			loop = func(vs []A) M[[]A] {
				return M[[]A]{
					func(ix int) Result[[]A] {
						r := p.core(src).f(ix)

						if r.SuccessQ() {
							v, ix := r.GetSuccess()

							return loop(append(vs, v)).f(ix)
						}

						return success[[]A]{vs, ix}
					},
				}
			}

			return loop([]A{})
		},
	)
}

func Eof() Parser[data.Unit] {
	return makeParser(
		func(src source) M[data.Unit] {
			return M[data.Unit]{
				func(ix int) Result[data.Unit] {
					if ix == len(src.str) {
						return success[data.Unit]{data.Unit{}, ix}
					}

					return failure[data.Unit]{}
				},
			}
		},
	)
}
