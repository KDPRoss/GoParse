package main

import (
	"github.com/kdpross/GoParse/pkg/data"
	"github.com/kdpross/GoParse/pkg/parse"
	parseext "github.com/kdpross/GoParse/pkg/parse_ext"
)

func bracketed[A any](p parse.Parser[A]) parse.Parser[A] {
	return parse.SeqRight(parse.Txt("("), parse.SeqLeft(p, parse.Txt(")")))
}

// This illustrates the subtleties of doing
// mutually-recursive values in a language like Go. We need
// to take care to wrap things in lazy cells in order to
// avoid runtime problems (even though the forward
// declarations will avoid *static* errors).
//
// The present functions could probably be written in a
// more-straightforward way, but this is establishing the
// pattern that we'll use (below) for parsing types.
var ParseKind, ParseKindS = func() (parse.Parser[Kind], parse.Parser[Kind]) {
	// Forward declare the combined parsers; because these are
	// lazy, they may be referenced by the clause parsers.
	var kindP, kindPS data.Lazy[parse.Parser[Kind]]
	// Forward declare the clause parsers; because these are
	// *not* lazy, they may only be referenced by the lazy
	// parsers (so as to avoid capturing the effectively-nil
	// reference).
	var kstar, karr parse.Parser[Kind]

	// Parse 'full' kinds.
	kindP = data.MkLazy(func() parse.Parser[Kind] { return parse.Alt(karr, parse.Cache(kindPS)) })
	// Parse 'simple' kinds.
	kindPS = data.MkLazy(func() parse.Parser[Kind] { return parse.Alt(kstar, bracketed(parse.Cache(kindP))) })

	kstar = parse.Proc(parse.Txt("*"), func(_ string) Kind { return KStar{} })
	karr = parse.Proc(
		parseext.SeqS1(
			parse.Cache(kindPS),
			parseext.SeqRightS1(parse.Txt("->"), parse.Cache(kindP)),
		),
		func(p data.Pair[Kind, Kind]) Kind { return KArr{p.First(), p.Second()} },
	)

	return parse.Cache(kindP), parse.Cache(kindPS)
}()

func alts[A any](p parse.Parser[A], ps ...parse.Parser[A]) parse.Parser[A] {
	res := p

	for _, pP := range ps {
		res = parse.Alt(res, pP)
	}

	return res
}

// And, here, we see how verbose doing this in Go is: If
// we'd had the ability to override operators, all of this
// could be written *substantially* more compactly (as,
// indeed, is the case in OCaml and Motmot). In Haskell, of
// course, we could dispense with the lazy cells entirely.
var ParseType, ParseTypeS = func() (parse.Parser[Type], parse.Parser[Type]) {
	var typeP, typePH, typePS data.Lazy[parse.Parser[Type]]
	var tvar, tabs, tcval, tarr, tapp, ttpl parse.Parser[Type]

	typeP = data.MkLazy(func() parse.Parser[Type] { return alts(tabs, tarr, parse.Cache(typePH)) })
	typePH = data.MkLazy(func() parse.Parser[Type] { return alts(tcval, tapp, parse.Cache(typePS)) })
	typePS = data.MkLazy(func() parse.Parser[Type] { return alts(tvar, ttpl) })

	varP := parse.Regexp("[a-z][a-zA-Z0-9]*")
	consP := parse.Regexp("[A-Z][a-zA-Z0-9]*")
	tvar = parse.Proc(varP, func(s string) Type { return TVar{s} })
	tabs = parse.Proc(
		parseext.SeqS1(
			bracketed(parseext.SeqS1(varP, parseext.SeqRightS1(parse.Txt(":"), ParseKind))),
			parseext.SeqRightS1(parse.Txt("=>"), parse.Cache(typeP)),
		),
		func(t data.Pair[data.Pair[string, Kind], Type]) Type {
			return TAbs{t.First().First(), t.First().Second(), t.Second()}
		},
	)
	tcval = func() parse.Parser[Type] {
		withArgs := parse.Proc(parseext.SeqS1(
			consP,
			parseext.RepSep1(parse.Cache(typePS), parseext.Spaces1),
		), func(p data.Pair[string, []Type]) Type { return TCVal{p.First(), p.Second()} })
		noArgs := parse.Proc(consP, func(s string) Type { return TCVal{s, []Type{}} })
		return parse.Alt(withArgs, noArgs)
	}()
	tarr = parse.Proc(parseext.SeqS1(
		parse.Cache(typePH),
		parseext.SeqRightS1(parse.Txt("->"), parse.Cache(typeP)),
	), func(p data.Pair[Type, Type]) Type { return TArr{p.First(), p.Second()} },
	)
	tapp = parse.Proc(
		parseext.RepSep1(parse.Cache(typePS), parseext.Spaces1),
		func(ts []Type) Type {
			var res Type

			for i, t := range ts {
				if i == 0 {
					res = t
				} else {
					res = TApp{res, t}
				}
			}

			return res
		},
	)
	ttpl = parse.Proc(
		bracketed(parseext.RepSep1(parse.Cache(typeP), parse.Regexp(",[ ]+"))),
		func(ts []Type) Type {
			if len(ts) == 1 {
				return ts[0]
			} else {
				return TTpl{ts}
			}
		},
	)

	return parse.Cache(typeP), parse.Cache(typePS)
}()
