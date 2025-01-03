# ┌─────────────────────────────────────────────────────────────┐
# │ GoParse: A Golang parser-combinator library.                │
# │                                                             │
# │ This codebase is licensed for the following purposes only:  │
# │                                                             │
# │ - study of the code                                         │
# │                                                             │
# │ - compiling / running an unaltered copy of the code for     │
# │   noncommercial educational and entertainment purposes only │
# │                                                             │
# │ - gratis redistribution of the code in entirety and in      │
# │   unaltered form for any aforementioned purpose             │
# │                                                             │
# │ Copyright 2022-2025, K.D.P.Ross                             │
# └─────────────────────────────────────────────────────────────┘

GO=go

.PHONY: fmt test exp examples

exp:
	rlwrap go run internal/demo/exp.go

examples:
	go run internal/demo/examples.go

fmt:
	go fmt ./...

test:
	go test ./pkg/...
	go test ./internal/motmot/...

.PHONY: markdown-lint
markdown-lint: node_modules
	yarn markdownlint --config .markdownlint.jsonc --ignore **/node_modules/** **/*.md *.md

node_modules:
	pnpm install
