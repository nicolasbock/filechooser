#!/usr/bin/env python

import setuptools
import subprocess


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
    license="BSD",
    url="https://github.com/nicolasbock/filechooser.git",
    project_urls={
        "Documentation": "https://setuptools.readthedocs.io/"
    },
    scripts=["scripts/autorotate.sh"],
    entry_points={
        "console_scripts": [
            "pick-files = filechooser.main:main"
        ]
    },
    packages=setuptools.find_packages(),
    test_suite="tests"
)
