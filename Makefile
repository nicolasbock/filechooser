filechooser: filechooser.go
	go build

.PHONY: test
test:
	go test -v
