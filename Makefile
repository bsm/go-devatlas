default: test

deps:
	go get -t ./...

test: deps
	go test ./...

bench: deps
	go test ./... -bench=.
