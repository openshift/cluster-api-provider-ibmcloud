FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.16-openshift-4.10 AS builder
WORKDIR /go/src/github.com/openshift/cluster-api-provider-ibmcloud
COPY . .
# VERSION env gets set in the openshift/release image and refers to the golang version, which interfers with our own
RUN unset VERSION \
 && GOPROXY=off NO_DOCKER=1 make build

FROM registry.ci.openshift.org/ocp/4.10:base
COPY --from=builder /go/src/github.com/openshift/cluster-api-provider-ibmcloud/bin/machine-controller-manager /
