#!/bin/bash
#
# Shamelessly copied from:
#
# http://superuser.com/questions/36645/how-to-rotate-images-automatically-based-on-exif-data

JHEAD=jhead
SED=sed
CONVERT=convert

# Change the location of the image files here.
for f in *.jpg; do
  orientation=$($JHEAD -v $f | $SED -nr 's:.*Orientation = ([0-9]+).*:\1:p')

  if [ -z $orientation ]
  then
    orientation=0
  fi

  if [ $orientation -gt 1 ]
  then
    echo Rotating $f...
    mv $f $f.bak
    $CONVERT -auto-orient $f.bak $f
  fi
done
