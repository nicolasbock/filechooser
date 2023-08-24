Supported Options
=================

.. code-block:: console

   Usage of pick-files-1.3.3:

   # Introduction

   pick-files is a script that randomly selects a specific number of files from a set of folders and copies these files to a single destination folder. During repeat runs the previously selected files are excluded from the selection for a specific time period that can be specified.

   ## Usage Example

   pick-files --number 20 --destination new_folder --suffix .jpg --suffix .avi --folder folder1 --folder folder2

   Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into new_folder. The new_folder is created if it does not exist already. In this example, only files with suffixes .jpg or .avi are considered.

   Tips and Tricks
   ===============

   Some tricks....

   -N, --number  (= 1)
       The number of files to choose.
   --block-selection (= "")
       Block selection of files for a certain period. Possible units are (s)econds, (m)inutes, (h)ours, (d)days, and (w)weeks.
   --debug  (= false)
       Debug output.
   --destination (= "output")
       The output PATH for the selected files.
   --destination-option  (= panic)
       What to do when writing to destination; possible options are panic, append, and delete.
   --dry-run  (= false)
       If set then the chosen files are only shown and not copied.
   --folder  (= )
       A folder PATH to consider when picking files; can be used multiple times; works recursively, meaning all sub-folders and their files are included in the selection.
   -h, --help  (= false)
       This help message.
   --print-database (= "")
       Print the internal database to a file and exit; the special name `-` means standard output.
   --print-database-format  (= CSV)
       Format of printed database; possible options are CSV, JSON, and YAML.
   --print-database-statistics  (= false)
       Print some statistics of the internal database.
   --reset-database  (= false)
       Reset the database (re-initialize). Use intended for testing only.
   --suffix  (= )
       Only consider files with this SUFFIX. For instance, to only load jpeg files you would specify either 'jpg' or '.jpg'. By default, all files are considered.
   --verbose  (= false)
       Verbose output.
   --version  (= false)
       Print the version of this program.
