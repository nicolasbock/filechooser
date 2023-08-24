#!/bin/bash

set -e -u -x

basedir=$(realpath $(dirname $0)/..)

: ${DESTINATION:=${basedir}/artifacts}

for i in a b c; do
  for j in a b c d; do
    mkdir -p ${DESTINATION}/${i}/${j}
    for k in $(seq 10); do
      for s in jpg png tif jpeg; do
        UUID=$(uuidgen)
        echo ${UUID} > ${DESTINATION}/${i}/${j}/$(printf "%03d-${UUID}.${s}" ${k})
      done
    done
  done
done
