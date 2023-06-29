.. |br| raw:: html

    <br />

.. image:: https://readthedocs.org/projects/filechooser/badge/?version=latest
    :target: https://filechooser.readthedocs.io/en/latest/?badge=latest
    :alt: Documentation Status

|br|

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml
    :alt: Build and test

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml
    :alt: Build and test

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml/badge.svg
    :target: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml
    :alt: Build and Test

|br|

.. image:: https://badge.fury.io/py/filechooser.svg
    :target: https://badge.fury.io/py/filechooser
    :alt: Python Package

|br|

.. image:: https://snapcraft.io/static/images/badges/en/snap-store-black.svg
    :target: https://snapcraft.io/pick-files
    :alt: Get it from the Snap Store

Introduction
============

A script that copies a random selection of files from a set of folders
to a single destination folder.

Usage Example
-------------

.. code-block:: console

    pick-files --number 20 --destination new_folder \
        --suffix .jpg --suffix .avi --folder folder1 --folder folder2

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
