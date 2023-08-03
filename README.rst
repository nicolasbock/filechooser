.. |br| raw:: html

    <br />

.. image:: https://snapcraft.io/static/images/badges/en/snap-store-black.svg
    :target: https://snapcraft.io/pick-files
    :alt: Get it from the Snap Store

|br|

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml
    :alt: Build and test Go code

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml
    :alt: Build and Test snap package

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/debian-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/debian-package.yaml
    :alt: Build and Test Debian package

|br|

.. image:: https://readthedocs.org/projects/filechooser/badge/?version=latest
    :target: https://filechooser.readthedocs.io/en/latest/?badge=latest
    :alt: Documentation Status

Introduction
============

A script that copies a random selection of files from a set of folders
to a single destination folder.

Usage Example
-------------

.. code-block:: console

    pick-files --number 20 --destination new_folder \
        --suffix .jpg --suffix .avi --folder folder1 --folder folder2

Supported options
-----------------

.. code-block:: console

    Usage of pick-files-1.3.1-10-g5e25d21:

    # Introduction

    pick-files is a script that randomly selects a specific number of files from a set of folders and copies these files to a single destination folder. During repeat runs the previously selected files are excluded from the selection for a specific time period that can be specified.

    ## Usage Example

    pick-files --number 20 --destination new_folder --suffix .jpg --suffix .avi --folder folder1 --folder folder2

    Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into new_folder. The new_folder is created if it does not exist already. In this example, only files with suffixes .jpg or .avi are considered.

    -N, --number  (= 1)
        The number of files to choose.
    --append  (= false)
        Append chosen files to existing destination folder.
    --debug  (= false)
        Debug output.
    --delete-existing  (= false)
        Delete existing files in the destination folder instead of moving those files to a new location.
    --destination (= "output")
        The output PATH for the selected files.
    --dry-run  (= false)
        If set then the chosen files are only shown and not copied.
    --folder  (= )
        A folder PATH to consider when picking files; can be used multiple times; works recursively, meaning all sub-folders and their files are included in the selection.
    -h, --help  (= false)
        This help message.
    --print-database  (= false)
        Print the internal database and exit.
    --print-database-format  (= CSV)
        Format of printed database; possible options are CSV and JSON.
    --suffix  (= )
        Only consider files with this SUFFIX. For instance, to only load jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.
    --verbose  (= false)
        Verbose output.
    --version  (= false)
        Print the version of this program.

Installation
------------

There are several options to install this script

Get the snap
------------

.. code-block:: console

    snap install pick-files

Build it from source
--------------------

.. code-block:: console

    make

Which requires ``>= go-1.20`` and ``make``.
