package machine

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func stubInstanceGetByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (*vpcv1.Instance, error) {
	returnName := name
	returnID := "0727_xyz-xyz-cccc-aaba-cacdaccad"
	returnPrimaryNetID := "0727-xyz"
	returnPrimaryNetName := "cold-breeze"
	returnPrimaryNetIPv4Add := "10.0.0.1"
	returnRunning := "running"

	return &vpcv1.Instance{
		Name: &returnName,
		ID:   &returnID,
		PrimaryNetworkInterface: &vpcv1.NetworkInterfaceInstanceContextReference{
			ID:                 &returnPrimaryNetID,
			Name:               &returnPrimaryNetName,
			PrimaryIpv4Address: &returnPrimaryNetIPv4Add,
		},
		Status: &returnRunning,
	}, nil
}

func stubMachine() (*machinev1.Machine, error) {
	userDataSecretName := "user-data-test"
	credentialsSecretName := "test-ic-credentials"
	defaultNamespaceName := "test-ns"

	machineSpec := &ibmcloudproviderv1.IBMCloudMachineProviderSpec{
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: credentialsSecretName,
		},
		UserDataSecret: &corev1.LocalObjectReference{
			Name: userDataSecretName,
		},
	}

	providerSpec, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(machineSpec)
	if err != nil {
		return nil, fmt.Errorf("codec.EncodeProviderSpec failed: %v", err)
	}

	machine := &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-testing-machine",
			Namespace: defaultNamespaceName,
			Labels:    map[string]string{},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: "machine.openshift.io/v1beta1",
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: providerSpec,
			},
		},
	}
	return machine, nil
}
