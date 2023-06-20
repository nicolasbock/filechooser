#!/bin/bash

set -e -u -x

for i in a b c; do
  for j in a b c d; do
    for k in $(seq 20); do
      mkdir -p artifacts/${i}/${j}
      UUID=$(uuidgen)
      echo ${UUID} > artifacts/${i}/${j}/$(printf "%03d-${UUID}.txt" ${k})
    done
  done
done
