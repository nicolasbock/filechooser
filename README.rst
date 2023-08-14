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

    pick-files --number 20 \
        --destination output \
        --suffix .jpg --suffix .avi \
        --folder folder1 --folder folder2

which copies 20 files chosen randomly from the files in `folder1` and
`folder2` with suffixes `.jpg` and `.avi` into the folder `output`.

For more details on all supported options, see the :doc:`Supported
Options<pick-files-help>` page or run the command with `--help`.

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
