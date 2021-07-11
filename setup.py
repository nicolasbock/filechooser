#!/usr/bin/env python

import os.path
import setuptools
import subprocess


def get_latest_tag():
    version = "unknown"
    try:
        with open('.version') as fd:
            version = fd.readline()
        return version
    except FileNotFoundError:
        pass

    git = subprocess.Popen(['git', 'describe', '--tags'],
                           stdout=subprocess.PIPE,
                           universal_newlines=True)
    git.wait()
    if git.returncode == 0:
        raw_version = git.stdout.readlines()[0][1:].rstrip().split('-')
        version = raw_version[0]
        if len(raw_version) > 1:
            version += ".dev" + raw_version[1]
    with open('.version', 'w') as fd:
        fd.write('{}\n'.format(version))
    return version


setuptools.setup(
    name="filechooser",
    version=get_latest_tag(),
    data_files=['.version'],
    license="BSD",
    url="https://github.com/nicolasbock/filechooser.git",
    project_urls={
        "Documentation": "https://filechooser.readthedocs.io/"
    },
    scripts=["scripts/autorotate.sh"],
    entry_points={
        "console_scripts": [
            "pick-files = filechooser_legacy.main:main",
            "pick-files-new = filechooser.main:main"
        ]
    },
    packages=setuptools.find_packages(),
    test_suite="tests"
)
