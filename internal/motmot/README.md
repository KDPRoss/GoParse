# MotmotLite Types Parser

Here is an extended example using the library to parse
MotmotLite types (and, as a result of their occurring
*within* types, kinds). Although it's not a full
implementation of the language's syntax, it encounters all
of the vaguely-subtle problems that would occur if we *had*
implemented a full parser (including, additionally,
expressions and patterns). In particular, it demonstrates
what's involved in having recursive productions in the
grammar, which, in an eagerly-evaluated language, requires
some care to avoid divergence (or panics due to references'
being captured before they're initialised).
