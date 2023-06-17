filechooser: filechooser.go
	go build

.PHONY: test
test:
	go test -v

build:
	rm -rf build
	./setup.py build

dist:
	rm -rf dist
	./setup.py bdist_wheel --universal

upload: dist
	twine upload dist/*
