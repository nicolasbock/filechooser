Supported Options
=================

.. code-block:: console

   Usage of pick-files:

   Introduction
   ============

   pick-files is a script that randomly selects a specific number of files from a set of folders and copies these files to a single destination folder. During repeat runs the previously selected files are excluded from the selection for a specific time period that can be specified.

   Usage Example
   -------------

       pick-files --number 20 \
           --destination output \
           --suffix .jpg --suffix .avi \
           --folder folder1 --folder folder2

   Would choose at random 20 files from folder1 and folder2 (including sub-folders) and copy those files into output. The output is created if it does not exist already. In this example, only files with suffixes .jpg or .avi are considered.

   Tips and Tricks
   ===============

   Daily copying of files
   ----------------------

   The snap includes example Systemd service and timer files that can be used to
   set up a daily update of the picked files.

   .. code-block:: console

      $ ls /snap/pick-files/current/usr/share/pick-files/pick-files-daily*
      pick-files-daily.service pick-files-daily.timer

   Run

   .. code-block:: console

      $ systemctl edit --user --force --full pick-files-daily.service

   and copy the contents of the example service file from the snap. Repeat with
   the timer.

   Then

   .. code-block:: console

      $ systemctl enable --user pick-files-daily.timer
      $ systemctl start --user pick-files-daily.timer

   will start the timer. The output of the pick-files command will be in the
   journal and can be checked with

   .. code-block:: console

      $ journalctl --unit pick-files-daily.service

   Options
   -------

   -N, --number  (= 1)
       The number of files to choose.
   --block-selection (= "")
       Block selection of files for a certain period. Possible units are (s)econds, (m)inutes, (h)ours, (d)days, and (w)weeks.
   --config (= "")
       Use configuration file
   --debug  (= false)
       Debug output.
   --destination (= "output")
       The output PATH for the selected files.
   --destination-option  (= panic)
       What to do when writing to destination; possible options are panic, append, and delete.
   --dry-run  (= false)
       If set then the chosen files are only shown and not copied.
   --dump-configuration  (= false)
       Dump current configuration; output can be used as configuration file.
   --folder  (= )
       A folder PATH to consider when picking files; can be used multiple times; works recursively, meaning all sub-folders and their files are included in the selection.
   -h, --help  (= false)
       This help message.
   --journald  (= false)
       Log to journald.
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
