# GoParse

## About

GoParse is an implementation of parser combinators in
Golang, based on ‘packrat parsing’. Currently, the code
in this repository is licensed only for study only and may
not be used for any productive purpose.

## Demonstration

Here are some simple parsers demonstrating the capabilities
of this library; see the `internal/demo` directory for a
simple interpreter providing a somewhat-more-extensive
example; the tests are also illustrative.

### Basic Parsers

* The simplest possible parser is one that matches a
  constant string:

  ```go
  fooParse := parse.Txt("foo")
  ```

  * Let's parse something (successfully) with it:

    ```go
    if r := parse.Parse(fooParse, "foo"); r.SuccessQ() {
      v, _ := r.GetSuccess()
      fmt.Printf("v = %q\n", v) // prints `v = "foo"`
    }
    ```

  * Let's see our first parse failure:

    ```go
    if r := parse.Parse(fooParse, "bar"); r.FailureQ() {
      fmt.Println("Ruh roh! Parse failure.")
    }
    ```

* We can also base our parsers on regular expressions ⟦
  These are PCRE regexps—the best sort. ⟧:

  ```go
  baaarParse := parse.Regexp("ba+r")
  ```

  * We might use this like:

    ```go
    if r := parse.Parse(baaarParse, "baaaaaaaaaar"); r.SuccessQ() {
      v, _ := r.GetSuccess()
      fmt.Printf("v = %q\n", v) // prints `v = "baaaaaaaaaar"`
    }
    ```

### Combining Parsers

The combinator-based approach (to parsing … and everything
else) is beautiful and powerful because of compositionality.

* We can sequence parsers:

  ```go
  fooBarParse := parse.Seq(fooParse, baaarParse)
  ```

  * This will match something like `"foobaaaaaaaaar"`—and
    (upon a successful parse) return a pair of values ‘in
    the obvious way’:

    ```go
    fooBarParse := parse.Seq(fooParse, baaarParse)
    if r := parse.Parse(fooBarParse, "foobaaaaaaaaar"); r.SuccessQ() {
      v, _ := r.GetSuccess()
      fmt.Printf("v = (%q, %v)\n", v.First(), v.Second()) // prints `v = ("foo", baaaaaaaaar)`
    }
    ```

* We can post-process the result of a parser:

  ```go
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
  ```

  * ⟦ Building some sort of abstract syntax tree would be
    a rather clever thing to do here! ⟧

* We can repeat some parser separated by a delimiter:

  ```go
  numList := parseext.RepSep(num, parse.Txt(","))
  ```

  ```go
  if r:=parse.Parse(numList, "1,2,12,57"); r.SuccessQ() {
    v, _ := r.GetSuccess()
    fmt.Printf("v = %v\n", v) // prints `v = [1 2 12 57]`
  }
  ```

  * We may also wish to allow spaces—and require matching
    a nonempty sequence, which we can do like this:

    ```go
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
    ```

### Recursion

The ‘elephant in the room’ whenever building any sort of
recursive parser in an applicative language is termination
problems that arise due to self-reference / recursion. Let's
build the ‘list of numbers’ parser that we considered
above … from only the basic combinators:

* First, let's write a ‘forward declaration’ of our
  parser:

  ```go
  var numList data.Lazy[parse.Parser[[]int]]
  ```

  * Notice that this is a *lazy-cell* reference to a parser;
    we'll need to take this into account when *using* the
    parser.

* Next, let's define some utility bindings—a ‘normal
  functional-style `cons`’ function and a trivial
  empty-list parser:

  ```go
  cons := func(p data.Pair[int, []int]) []int {
    return append([]int{p.First()}, p.Second()...)
  }
  emptyList := parse.ParserJust([]int{})
  ```

* Using these, we can define a parser for a nonempty list:

  ```go
  nonemptyList := parse.Proc(parse.Seq(num, parse.Alt(nonemptyTail, emptyList)), cons)
  ```

  * It's a number followed by either a comma *and* a
    nonempty list xor a number; concatenate the list head /
    tail.

* And, a list is either a nonempty list xor an empty list:

  ```go
  parse.Alt(nonemptyList, emptyList)
  ```

* Putting this all together, we have:

  ```go
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
  ```

* And, we can use this like:

  ```go
  if r := parse.Parse(parse.Cache(numListCustom), "1,23,456"); r.SuccessQ() {
    v, _ := r.GetSuccess()
    fmt.Printf("v = %v\n", v) // prints `v = [1 23 456]`
  }
  ```

Returning to the beginning of this section: The lazy cell
combined with the `parse.Cache` combinator allows us to
control evaluation … and is key to self-reference (or
mutual-self-reference / -recursion) in our parsers.

## Observations / Metalevel Discussion

I wrote this library primarily to get a sense of what it
‘feels like’ to write fundamentally-polymorphic code in
Golang. I wrote this code across a few sessions over a
weekend, and I noticed a few things:

* In VSCode, at least, it helped the type system immensely
  to provide seemingly-needless type annotations. Once the
  code was completed, various square-bracketed annotations
  were marked as superfluous, at which point I removed them.
  Reaching this point with any particular polymorphic
  function was *shockingly* satisfying—it's a bit like
  achieving a Tetris.

* It's a bit disappointing that Golang was released with
  polymorphism / ‘template functions’ … but no library
  to leverage them—some version of STL.

* From a functional perspective, it's *extremely weird* that
  ‘cons’ (viz., `append`) adds elements *to the right
  side* of a list rather than to the left. ⟦ Of course,
  slices are *not* cons-cell-based lists at all! ⟧

* Although I'm entirely happy editing OCaml or Motmot or
  Haskell or Tanager code in Emacs, I would *not* want to
  write (polymorphic or otherwise) Golang code without an
  IDE, which suggests that there's a good lot of redundancy
  required.

* Testing is amazing. I write buggy code, and so do you.
  Testing proves—always and only—the presence of bugs,
  never their absence. And, I applaud Golang for its testing
  logistics—which made writing this library very much
  easier. QuickCheck is the best—and you'd better believe
  that I'll be porting it to Golang—but … it's not very
  useful for libraries as abstract as this.

## Caveats

* If Golang had a more capable polymorphic type system, we
  could avoid numerous type annotations that Golang
  requires.

* If Golang supported infix function calls / ‘operator
  overloading’, we could make substantially-shorter /
  -more-readable parsers.

* If Golang could handle polymorphic type aliases /
  `typedef`s, the *implementation* of the combinators could
  be simplified a bit.

* If Golang were lazy, we could dispense with the explicit
  `data.Lazy` cells.

I have, historically, been a Golang hater, but this library
shows that—with, at long last, the addition of
(parametric) polymorphism / ‘generics’—it's *finally*
possibly to write libraries that use monads and abstract
things interestingly. I'm chuffed to bits about being able
to do *this* … rather than the truly-deplorable
unsafe-coercion-based parser-combinator implementation that
I made in pre-1.18 Golang.

## Resources

* [a fantastic ‘Computerphile’ video about parser
  combinators with Graham
  Hutton](https://youtu.be/dDtZLm7HIJs)

* [The Packrat Parsing and Parsing Expression Grammars
  Page](https://bford.info/packrat/)
