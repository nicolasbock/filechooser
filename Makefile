pick-files: pick-files.go
	go build -ldflags "-X main.Version=$(shell git describe --tags)" -o pick-files ./...

.PHONY: test
test:
	go test -v

.PHONY: coverage
coverage:
	go test -cover -v
