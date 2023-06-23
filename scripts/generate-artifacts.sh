#!/bin/bash

set -e -u -x

for i in a b c; do
  for j in a b c d; do
    mkdir -p artifacts/${i}/${j}
    for k in $(seq 10); do
      for s in jpg png tif; do
        UUID=$(uuidgen)
        echo ${UUID} > artifacts/${i}/${j}/$(printf "%03d-${UUID}.${s}" ${k})
      done
    done
  done
done
