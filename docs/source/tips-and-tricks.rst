Tips and Tricks
===============

Daily copying of files
----------------------

The snap includes example Systemd service and timer files that can be used to
set up a daily update of the picked files.

.. code-block::

   $ ls /snap/pick-files/current/usr/share/pick-files/pick-files-daily*
   pick-files-daily.service pick-files-daily.timer

Run

.. code-block::

   $ systemctl edit --user --force --full pick-files-daily.service

and copy the contents of the example service file from the snap. Repeat with
the timer.

Then

.. code-block::

   $ systemctl enable --user pick-files-daily.timer
   $ systemctl start --user pick-files-daily.timer

will start the timer. The output of the pick-files command will be in the
journal and can be checked with

.. code-block::

   $ journalctl --unit pick-files-daily.service
