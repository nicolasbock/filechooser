.PHONY: build dist upload clean

build:
	rm -rf build
	./setup.py build

dist:
	rm -rf dist
	./setup.py bdist_wheel --universal

upload: dist
	twine upload dist/*
