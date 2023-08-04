#!/bin/bash

set -e -u -x

basedir=$(realpath $(dirname $0))
builddir=$(realpath ${1:-$(realpath $(mktemp --directory --tmpdir=${basedir}/.. doc-build.XXXXXX))})

python3 -m venv ${builddir}/doc-venv
${builddir}/doc-venv/bin/pip install --upgrade pip
${builddir}/doc-venv/bin/pip install sphinx sphinx_rtd_theme

cd ${basedir}/../docs
${builddir}/doc-venv/bin/sphinx-build -M html source ${builddir}/build
