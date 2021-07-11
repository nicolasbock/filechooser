Welcome to filechooser's documentation!
=======================================

.. toctree::
   :maxdepth: 2
   :caption: Contents

   design_goals

Introduction
------------

.. image:: https://readthedocs.org/projects/filechooser/badge/?version=latest
   :target: https://filechooser.readthedocs.io/en/latest/?badge=latest
   :alt: Documentation Status

.. image:: https://badge.fury.io/py/filechooser.svg
    :target: https://badge.fury.io/py/filechooser

A script that copies a random selection of files from a set of folders
to a single destination folder.

Installation
------------

The easiest way to install this script is to run

.. code::

   pip install filechooser

Usage Example
-------------

.. code::

   pick-files -N 20 --destination new_folder --suffix .jpg .avi .h -- folder1 folder2

Indices and tables
------------------

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
