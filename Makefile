pick-files: *.go
	go build -ldflags "-X main.Version=$(shell git describe --tags)" -v -o pick-files ./...

.PHONY: test
test:
	go test -v

.PHONY: coverage
coverage:
	go test -cover -v

.PHONY: clean
clean:
	rm pick-files
