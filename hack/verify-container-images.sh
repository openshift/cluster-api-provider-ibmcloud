#!/bin/bash

# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

if [[ "${TRACE-0}" == "1" ]]; then
    set -o xtrace
fi

VERSION=${1}
GO_ARCH="$(go env GOARCH)"

REPO_ROOT=$(git rev-parse --show-toplevel)
"${REPO_ROOT}/hack/ensure-trivy.sh" "${VERSION}"

TRIVY="${REPO_ROOT}/hack/tools/bin/trivy/${VERSION}/trivy"

# Builds all the container images to be scanned and cleans up changes to ./*manager_image_patch.yaml ./*manager_pull_policy.yaml.
make REGISTRY=gcr.io/k8s-staging-capi-ibmcloud PULL_POLICY=IfNotPresent TAG=dev OUTPUT_TYPE=type=docker docker-build
make clean-release-git

make -C hack/ccm ARCH="${GO_ARCH}" TAG=dev build-local

# Scan the images
"${TRIVY}" image -q --exit-code 1 --ignore-unfixed --severity MEDIUM,HIGH,CRITICAL gcr.io/k8s-staging-capi-ibmcloud/cluster-api-ibmcloud-controller-"${GO_ARCH}":dev && R1=$? || R1=$?
"${TRIVY}" image -q --exit-code 1 --ignore-unfixed --severity MEDIUM,HIGH,CRITICAL gcr.io/k8s-staging-capi-ibmcloud/powervs-cloud-controller-manager-"${GO_ARCH}":dev && R2=$? || R2=$?

echo ""
BRed='\033[1;31m'
BGreen='\033[1;32m'
NC='\033[0m' # No

if [ "$R1" -ne "0" ] || [ "$R2"  -ne "0" ]
then
  echo -e "${BRed}Check container images failed! There are vulnerabilities to be fixed${NC}"
  exit 1
fi

echo -e "${BGreen}Check container images passed! No vulnerability found${NC}"
