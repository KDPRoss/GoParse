package main

import (
	"testing"

	"github.com/kdpross/GoParse/pkg/parse"
	"github.com/stretchr/testify/assert"
)

func TestKind(t *testing.T) {
	p := parse.SeqLeft(ParseKind, parse.Eof())

	for _, c := range []struct{ lab, s1, s2 string }{
		{"star", "*", "*"},
		{"star brackets", "((*))", "*"},
		{"arrow", "* -> *", "(* -> *)"},
		{"arrow brackets", "((*) -> (*))", "(* -> *)"},
		{"arrow nested", "* -> * -> * -> *", "(* -> (* -> (* -> *)))"},
		{"arrow nested brackets", "((*) -> ((*) -> ((*) -> (*))))", "(* -> (* -> (* -> *)))"},
	} {
		t.Run(c.lab, func(t *testing.T) {
			v := parse.Parse(p, c.s1)
			assert.True(t, v.SuccessQ())
			k, _ := v.GetSuccess()
			assert.Equal(t, c.s2, ShowKind(k))
		})
	}
}
