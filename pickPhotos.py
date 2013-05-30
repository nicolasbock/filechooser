#!/usr/bin/env python

import argparse
import math
import os
import os.path
import random
import shutil
import sys
import tempfile

parser = argparse.ArgumentParser(description = """This script randomly picks a
chosen numer of files from a set of folders and copies those files to a single
destination folder. The files to be considered can be filtered by suffix.""")

parser.add_argument("DIR",
    help = "The directory to recursively consider",
    nargs = "+")

parser.add_argument("-N",
    help = "Choose N photos randomly",
    type = int,
    default = 10)

parser.add_argument("--destination",
    metavar = "DIR",
    help = "Copy photos to DIR")

parser.add_argument("--suffix",
    help = """Only consider files with this SUFFIX. The leading '.' has to be
included, i.e. jpeg files could be added as '.jpg'. By default, all files are
be considered.""",
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

photos = []
for path in options.DIR:
  for root, dirs, files in os.walk(path):
    for i in files:
      if options.suffix != None:
         ( basename, extension ) = os.path.splitext(i)
         if not extension in options.suffix:
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
    backupfolder = tempfile.mkdtemp()
    print("moving old files to %s" % (backupfolder))
    for root, dirs, files in os.walk(options.destination):
      for i in files:
        shutil.move(os.path.join(root, i), backupfolder)

for i in selectedPhotos:
  if options.verbose:
    print("copying %s" % (i))
  if options.destination != None:
    shutil.copy(i, options.destination)
