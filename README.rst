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

`pick-files` is a script that randomly selects a specific number of files from
a set of folders and copies these files to a single destination folder. During
repeat runs the previously selected files are excluded from the selection for
a specific time period that can be specified.

Usage Example
-------------

.. code-block:: console

    pick-files --number 20 \
        --destination output \
        --suffix .jpg --suffix .avi \
        --folder folder1 --folder folder2

Would choose at random 20 files from `folder1` and `folder2` (including
sub-folders) and copy those files into `output`. The `output` is created if it
does not exist already. In this example, only files with suffixes `.jpg` or
`.avi` are considered.

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

Code Development
----------------

The code is hosted on GitHub at <https://github.com/nicolasbock/filechooser>.
If you are interested in contributing in any way, please have a look at the
repository.
