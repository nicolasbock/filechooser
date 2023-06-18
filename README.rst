.. image:: https://readthedocs.org/projects/filechooser/badge/?version=latest
   :target: https://filechooser.readthedocs.io/en/latest/?badge=latest
   :alt: Documentation Status

.. image:: https://badge.fury.io/py/filechooser.svg
   :target: https://badge.fury.io/py/filechooser

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml/badge.svg
   :target: https://github.com/nicolasbock/filechooser/actions/workflows/python-package.yaml
   :alt: Build and test

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml/badge.svg
   :target: https://github.com/nicolasbock/filechooser/actions/workflows/go-package.yaml
   :alt: Build and test

.. image:: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml/badge.svg
   :target: https://github.com/nicolasbock/filechooser/actions/workflows/snap-package.yaml
   :alt: Build and test

Introduction
------------

A script that copies a random selection of files from a set of folders
to a single destination folder.

Installation
------------

There are several options to install this script

.. code::

   snap install pick-files

or

.. code::

   pip install filechooser

Usage Example
-------------

.. code::

   pick-files --number 20 --destination new_folder --suffix .jpg .avi .h -- folder1 folder2
