#!/usr/bin/env python

import setuptools
from filechooser import __version__


setuptools.setup(
    name="filechooser",
    version=__version__,
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
