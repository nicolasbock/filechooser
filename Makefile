build:
	./setup.py build

dist:
	./setup.py bdist_wheel --universal

upload: clean dist
	twine upload dist/*

clean:
	git clean -dfx
