#!/usr/bin/env python

import argparse
import math
import os
import os.path
import random
import shutil
import sys
import tempfile

parser = argparse.ArgumentParser()

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

options = parser.parse_args()

photos = []
for path in options.DIR:
  for root, dirs, files in os.walk(path):
    for i in files:
      photos.append(os.path.join(root, i))

selectedPhotos = []
for i in range(options.N):
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
  print("copying %s" % (i))
  if options.destination != None:
    shutil.copy(i, options.destination)
