[![Documentation Status](https://readthedocs.org/projects/filechooser/badge/?version=latest)](https://filechooser.readthedocs.io/en/latest/?badge=latest)
[![Build and test](https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml/badge.svg)](https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml)
[![Build and test](https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml/badge.svg)](https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml)
[![Build and test](https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml/badge.svg)](https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml)

[![Python](https://badge.fury.io/py/filechooser.svg)](https://badge.fury.io/py/filechooser)

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/pick-files)

# Introduction

A script that copies a random selection of files from a set of folders
to a single destination folder.

# Installation

There are several options to install this script

```console
snap install pick-files
```

or alternatively to build the binary:

```console
make
```

# Usage Example

```console
pick-files --number 20 --destination new_folder --suffix .jpg .avi -- folder1 folder2
```
