#!/usr/bin/env bash

set -euo pipefail

# map names of CRD files between the vendored openshift/api repository and the ./install directory
CRDS_MAPPING=( "0000_10_machine.crd.yaml:machine.openshift.io.crd.yaml"
               "0000_10_machineset.crd.yaml:machineset.openshift.io.crd.yaml"
               "0000_10_machinehealthcheck.yaml:machinehealthcheck.openshift.io.crd.yaml" )

for crd in "${CRDS_MAPPING[@]}" ; do
    SRC="${crd%%:*}"
    DES="${crd##*:}"
    cp "vendor/github.com/openshift/api/machine/v1beta1/$SRC" "config/crds/$DES"
done
