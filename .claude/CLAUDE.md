# Project Configuration

## Basic rules

1. Prioritize std go packages to third-party
2. Use idiomatic go patterns
3. Avoid `panic` in situations when simple `result, err := ...` would work
4. Avoid putting too much of the defensive code, like checking each `err` coming from Writers
5. Avoid inline comments, keep them on the function level
   - a too complex function should be refactored instead of putting loads of inline comments
   - not all functions need `/// docs` but most of the non-trivial functions do

## Tools

1. Do not make web requests to arbitrary pages, only standard Golang repos and places

## Code Style

1. Follow `go fmt` rules
