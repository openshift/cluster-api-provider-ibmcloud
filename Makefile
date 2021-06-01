
ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

GOARCH  ?= $(shell go env GOARCH)
GOOS    ?= $(shell go env GOOS)

VERSION     ?= $(shell git describe --always --abbrev=7)
REPO_PATH   ?= github.com/openshift/cluster-api-provider-ibmcloud
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION) -extldflags "-static"


NO_DOCKER ?= 0
ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
  export CGO_ENABLED
else
  DOCKER_CMD = docker run --rm -e CGO_ENABLED=0 -e GOARCH=$(GOARCH) -e GOOS=$(GOOS) -v "$(PWD)":/go/src/github.com/openshift/cluster-api-provider-ibmcloud:Z -w /go/src/openshift/cluster-api-provider-ibmcloud openshift/origin-release:golang-1.15
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: build
build: ## build binaries
	$(DOCKER_CMD) CGO_ENABLED=0 go build $(GOGCFLAGS) -o "bin/machine-controller-manager" \
               -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/manager"
	$(DOCKER_CMD) CGO_ENABLED=0 go build  $(GOGCFLAGS) -o "bin/termination-handler" \
	             -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/termination-handler"