# Copyright 2021.
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


ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

GOARCH  ?= $(shell go env GOARCH)
GOOS    ?= $(shell go env GOOS)

VERSION     ?= $(shell git describe --always --abbrev=7)
REPO_PATH   ?= github.com/openshift/cluster-api-provider-ibmcloud
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION) -extldflags "-static"
IMAGE        = origin-ibmcloud-machine-controllers
MUTABLE_TAG ?= latest

# # race tests need CGO_ENABLED, everything else should have it disabled
CGO_ENABLED = 0
unit : CGO_ENABLED = 1

NO_DOCKER ?= 0
ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
  export CGO_ENABLED
else
  DOCKER_CMD = docker run --rm -e CGO_ENABLED=$(CGO_ENABLED) -e GOARCH=$(GOARCH) -e GOOS=$(GOOS) -v "$(PWD)":/go/src/github.com/openshift/cluster-api-provider-ibmcloud:Z -w /go/src/github.com/openshift/cluster-api-provider-ibmcloud openshift/origin-release:golang-1.16
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: check
check: fmt vet lint test # Check your code

.PHONY: generate
generate: gogen goimports
	./hack/verify-diff.sh

gogen:
	$(DOCKER_CMD) go generate ./pkg/... ./cmd/...

.PHONY: fmt
fmt: ## Go fmt your code
	$(DOCKER_CMD) hack/go-fmt.sh .

.PHONY: lint
lint: ## Go lint your code
	$(DOCKER_CMD) hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client/mock')

.PHONY: goimports
goimports: ## Go fmt your code
	$(DOCKER_CMD) hack/goimports.sh .

.PHONY: vet
vet: ## Apply go vet to all go files
	$(DOCKER_CMD) hack/go-vet.sh ./...

.PHONY: crds-sync
crds-sync: ## Sync crds in install with the ones in the vendored oc/api
	$(DOCKER_CMD) hack/crds-sync.sh .

.PHONY: verify-crds-sync
verify-crds-sync: ## Verify that the crds in install and the ones in vendored oc/api are in sync
	$(DOCKER_CMD) hack/crds-sync.sh . && hack/verify-diff.sh .

.PHONY: test
test: ## Run tests
	@echo -e "\033[32mTesting...\033[0m"
	$(DOCKER_CMD) hack/ci-test.sh

.PHONY: unit
unit: # Run unit test
	$(DOCKER_CMD) go test -race -cover ./cmd/... ./pkg/...

.PHONY: build
build: ## build binaries
	$(DOCKER_CMD) CGO_ENABLED=0 go build $(GOGCFLAGS) -o "bin/machine-controller-manager" \
               -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/manager"

.PHONY: images
images: ## Create images
ifeq ($(NO_DOCKER), 1)
	./hack/imagebuilder.sh
endif
	$(IMAGE_BUILD_CMD) -t "$(IMAGE):$(VERSION)" -t "$(IMAGE):$(MUTABLE_TAG)" ./

