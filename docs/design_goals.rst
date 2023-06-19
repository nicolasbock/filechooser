Design Goals
============

.. toctree::
   :maxdepth: 2
   :caption: Contents

1. Selection

  - Select `N` random pictures from multiple folders
  - The selection is recursive
  - The highest order folder is picked for selection
  - More than one folder can be picked for selection
  - Several different frame categories combinations exist, e.g. `a`,
    `b`, `c`, `d`. Frame categories may be combined with each other,
    e.g. `a` + `b` or `a` + `c`. Each combination has its own
    selection, download and frame display.
  - When selected the pictures are copied into a specific file path.

2. Selection Time Stamp

The selection should exclude pictures that have been selected within a
certain time interval.

   - When a picture is selected, a time stamp is created for this
     picture
   - The time stamp indicates time of last creation
   - A time stand is needed for each category
   - A data base might be needed for tracking the time stamps

3. Picture exclusion

   - For each selection a time frame in days is chosen that a picture
     cannot be selected again
   - A process running in the background determines which pictures are
     in the picture selection pool for each (combined) frame category.
     For example a picture can be in the selection pool for `a` + `c`,
     but not in the one for `a` + `b`.

4. Other

    - Reset option for the time stamp
    - Call selection by a cron job
    - Adding image folders needs to be possible
    - Remember past choices so that files will not reappear within a
      configurable time window
