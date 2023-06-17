.PHONY: build dist upload

build:
	rm -rf build
	/usr/bin/env python3 ./setup.py build

dist:
	rm -rf dist
	/usr/bin/env python3 ./setup.py bdist_wheel --universal

upload: dist
	twine upload dist/*
