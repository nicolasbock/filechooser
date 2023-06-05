filechooser: filechooser.go
	go build -o pick-files

.PHONY: test
test:
	go test -v

.PHONY: build
build:
	rm -rf build
	/usr/bin/env python3 ./setup.py build

.PHONY: dist
dist:
	rm -rf dist
	/usr/bin/env python3 ./setup.py bdist_wheel --universal

.PHONY: upload
upload: dist
	twine upload dist/*
