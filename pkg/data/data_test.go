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

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPairCoherent(t *testing.T) {
	n := 5
	s := "foozle"
	p := MkPair(n, s)

	assert.Equal(t, n, p.First())
	assert.Equal(t, s, p.Second())
}

func TestMaybeJustCoherent(t *testing.T) {
	n := 5
	m := MkJust(n)

	assert.False(t, m.NothingQ())
	require.True(t, m.JustQ())
	assert.Equal(t, n, m.GetJust())
}

func TestMaybeNothingCoherent(t *testing.T) {
	m := MkNothing[int]()

	assert.False(t, m.JustQ())
	assert.True(t, m.NothingQ())
}

func TestLazyIdempotent(t *testing.T) {
	n := 5
	callCount := 0
	k := func() int {
		callCount++

		return n
	}
	lz := MkLazy(k)

	assert.Equal(t, n, lz.Force())
	_ = lz.Force()
	assert.Equal(t, 1, callCount)
}
