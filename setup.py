#!/usr/bin/env python

import setuptools
import subprocess


def readme():
    with open("README.md") as fd:
        return fd.read()


def get_latest_tag():
    git = subprocess.Popen(['git', 'describe', '--tags'],
                           stdout=subprocess.PIPE,
                           universal_newlines=True)
    git.wait()
    if git.returncode == 0:
        raw_version = git.stdout.readlines()[0][1:].rstrip().split('-')
        version = raw_version[0]
        if len(version) > 1:
            version += ".dev" + raw_version[1]
        return version
    else:
        return "unknown"


setuptools.setup(
    name="filechooser",
    version=get_latest_tag(),
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
