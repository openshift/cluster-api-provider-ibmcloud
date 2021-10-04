module github.com/openshift/cluster-api-provider-ibmcloud

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.4.2
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v0.6.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.4.4
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/openshift/api v0.0.0-20210416115537-a60c0dc032fd
	github.com/openshift/machine-api-operator v0.2.1-0.20210504014029-a132ec00f7dd
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.0-beta.1.0.20210512131817-ce2f0c92d77e
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20201216171336-0b00fb8d96ac
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20201209184807-075372e2ed03
)
