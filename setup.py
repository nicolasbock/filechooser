#!/usr/bin/env python

import setuptools


def readme():
    with open("README.md") as fd:
        return fd.read()


setuptools.setup(
    name="filechooser",
    version="0.1.2",
    author="Nicolas Bock",
    author_email="nicolasbock@gmail.com",
    description="A script that copies a random selection of files from "
    "a set of folders to a single destination folder",
    long_description=readme(),
    license="BSD",
    url="https://github.com/nicolasbock/filechooser.git",
    scripts=["scripts/autorotate.sh"],
    entry_points={
        "console_scripts": [
            "pick-files = filechooser.main:main"
        ]
    },
    packages=setuptools.find_packages(),
    test_suite="tests"
)
