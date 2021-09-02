/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	expectedProviderSpec = IBMCloudMachineProviderSpec{
		VPC:           "test-vpc",
		Image:         "test-instance-image",
		Profile:       "bx2d-32x128",
		DedicatedHost: "dedicated-host-name",
		Zone:          "us-east-1",
		Region:        "us-east",
		ResourceGroup: "aaab-bbbc-cccd-zzzz-xxxx",
		UserDataSecret: &corev1.LocalObjectReference{
			Name: "userData",
		},
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: "credentialKey",
		},
		PrimaryNetworkInterface: NetworkInterface{
			Subnet: "test-subnet",
			SecurityGroups: []string{
				"test-security-group-1",
				"test-security-group-2",
			},
		},
	}
	expectedRawProviderSpec = `{"metadata":{"creationTimestamp":null},"vpc":"test-vpc","image":"test-instance-image","profile":"bx2d-32x128","dedicatedHost":"dedicated-host-name","region":"us-east","zone":"us-east-1","resourceGroup":"aaab-bbbc-cccd-zzzz-xxxx","primaryNetworkInterface":{"subnet":"test-subnet","securityGroups":["test-security-group-1","test-security-group-2"]},"userDataSecret":{"name":"userData"},"credentialsSecret":{"name":"credentialKey"}}`

	instanceID             = "test-instance-id"
	instanceState          = "running"
	expectedProviderStatus = IBMCloudMachineProviderStatus{
		InstanceID:    &instanceID,
		InstanceState: &instanceState,
	}
	expectedRawProviderStatus = `{"instanceId":"test-instance-id","instanceState":"running"}`
)

func TestRawExtensionFromProviderSpec(t *testing.T) {
	rawExtension, err := RawExtensionFromProviderSpec(&expectedProviderSpec)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(rawExtension.Raw) != expectedRawProviderSpec {
		t.Errorf("Expected: %s, got: %s", expectedRawProviderSpec, string(rawExtension.Raw))
	}
}

func TestRawExtensionFromProviderStatus(t *testing.T) {
	rawExtension, err := RawExtensionFromProviderStatus(&expectedProviderStatus)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(rawExtension.Raw) != expectedRawProviderStatus {
		t.Errorf("Expected: %s, got: %s", expectedRawProviderStatus, string(rawExtension.Raw))
	}
}

func TestProviderSpecFromRawExtension(t *testing.T) {
	rawExtension := runtime.RawExtension{
		Raw: []byte(expectedRawProviderSpec),
	}
	providerSpec, err := ProviderSpecFromRawExtension(&rawExtension)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if reflect.DeepEqual(providerSpec, expectedProviderSpec) {
		t.Errorf("Expected: %v, got: %v", expectedProviderSpec, providerSpec)
	}
}

func TestProviderStatusFromRawExtension(t *testing.T) {
	rawExtension := runtime.RawExtension{
		Raw: []byte(expectedRawProviderStatus),
	}
	providerStatus, err := ProviderSpecFromRawExtension(&rawExtension)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if reflect.DeepEqual(providerStatus, expectedProviderStatus) {
		t.Errorf("Expected: %v, got: %v", expectedProviderStatus, providerStatus)
	}
}
