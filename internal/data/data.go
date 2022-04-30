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

package data

// Some basic data structures that somebody 'forgot to
// implement'.

type Unit struct{}

type Pair[A, B any] struct {
	l A
	r B
}

func MkPair[A, B any](l A, r B) Pair[A, B] {
	return Pair[A, B]{l, r}
}

func (p Pair[A, B]) First() A {
	return p.l
}

func (p Pair[A, B]) Second() B {
	return p.r
}

type Maybe[A any] interface {
	JustQ() bool
	NothingQ() bool
	GetJust() A
}

type Just[A any] struct {
	x A
}

func MkJust[A any](x A) Just[A] {
	return Just[A]{x}
}

var _ Maybe[int] = Just[int]{3}

func (j Just[A]) GetJust() A {
	return j.x
}

func (Just[A]) JustQ() bool {
	return true
}

func (Just[A]) NothingQ() bool {
	return false
}

type Nothing[A any] struct{}

func MkNothing[A any]() Nothing[A] {
	return Nothing[A]{}
}

var _ Maybe[int] = Nothing[int]{}

func (Nothing[A]) GetJust() A {
	panic("unimplemented")
}

func (Nothing[A]) JustQ() bool {
	return false
}

func (Nothing[A]) NothingQ() bool {
	return true
}

type Lazy[A any] struct {
	avail bool
	v     A
	k     func() A
}

func (l *Lazy[A]) Force() A {
	if !l.avail {
		l.v = l.k()
		l.avail = true
	}

	return l.v
}

func MkLazy[A any](f func() A) Lazy[A] {
	return Lazy[A]{
		avail: false,
		k:     f,
	}
}
