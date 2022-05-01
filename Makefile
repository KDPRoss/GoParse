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
# │ Copyright 2022, K.D.P.Ross                                  │
# └─────────────────────────────────────────────────────────────┘

GO=go

.PHONY: fmt test demo

demo:
	rlwrap go run internal/demo/exp.go

fmt:
	go fmt ./...

test:
	go test ./...
