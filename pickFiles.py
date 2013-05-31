#!/usr/bin/env python

import argparse
import math
import os
import os.path
import random
import shutil
import stat
import sys
import tempfile

parser = argparse.ArgumentParser(description = """This script randomly picks a
chosen numer of files from a set of folders and copies those files to a single
destination folder. The files to be considered can be filtered by suffix.""")

parser.add_argument("DIR",
    help = "The directory to recursively consider",
    nargs = "+")

parser.add_argument("-N",
    help = "Choose N files randomly",
    type = int,
    default = 10)

parser.add_argument("--destination",
    metavar = "DIR",
    help = "Copy files to DIR")

parser.add_argument("--delete-existing",
    help = """Delete existing files in the destionation folder instead of
moving those files to a new location.""",
    action = "store_true",
    default = False)

parser.add_argument("--suffix",
    help = """Only consider files with this SUFFIX. The leading '.' has to be
included, i.e. jpeg files could be added as '.jpg'. By default, all files are
be considered. The suffix is case insensitive.""",
    action = "append",
    nargs = "+")

parser.add_argument("--verbose",
    help = "Print a lot of stuff",
    action = "store_true",
    default = False)

options = parser.parse_args()

if options.suffix != None:
  temp = []
  for i in options.suffix:
    temp.extend(i)
  options.suffix = temp

  for i in range(len(options.suffix)):
    options.suffix[i] = options.suffix[i].lower()

photos = []
for path in options.DIR:
  for root, dirs, files in os.walk(path):
    for i in files:
      if options.suffix != None:
         ( basename, extension ) = os.path.splitext(i)
         if not extension.lower() in options.suffix:
           continue
      photos.append(os.path.join(root, i))

if len(photos) == 0:
  print("could not find any files to select from")
  sys.exit(0)

selectedPhotos = []
if len(photos) >= options.N:
  N = options.N
else:
  N = len(photos)

for i in range(N):
  selectedPhotos.append(photos.pop(int(math.floor(random.random()*len(photos)))))

if options.destination != None:
  try:
    os.mkdir(options.destination)
  except OSError as e:
    print("destination path already exists: %s" % (e))

    if options.delete_existing:
      print("deleting existing files")
      for root, dirs, files in os.walk(options.destination):
        for i in files:
          os.remove(os.path.join(root, i))
    else:
      backupfolder = tempfile.mkdtemp()
      print("moving old files to %s" % (backupfolder))
      for root, dirs, files in os.walk(options.destination):
        for i in files:
          shutil.move(os.path.join(root, i), backupfolder)

  # Set very permissible permissions on destination directory.
  os.chmod(options.destination,
      stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR |
      stat.S_IRGRP | stat.S_IWGRP | stat.S_IXGRP |
      stat.S_IROTH | stat.S_IWOTH | stat.S_IXOTH)

for i in selectedPhotos:
  if options.verbose:
    print("copying %s" % (i))
  if options.destination != None:
    shutil.copy(i, options.destination)
