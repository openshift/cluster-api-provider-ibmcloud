module github.com/openshift/cluster-api-provider-ibmcloud

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.4.2
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v0.6.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v1.2.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.5.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/openshift/api v0.0.0-20211217221424-8779abfbd571
	github.com/openshift/machine-api-operator v0.2.1-0.20211220105028-362d5b50beca
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	k8s.io/klog/v2 v2.30.0
	sigs.k8s.io/controller-runtime v0.11.0
	sigs.k8s.io/controller-tools v0.7.0
	sigs.k8s.io/yaml v1.3.0
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20211122165613-e8f29d2a999f
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20211123175116-02b96338fcac
)
