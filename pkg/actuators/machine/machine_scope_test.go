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

package machine

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	mockibm "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client/mock"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	controllerfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func init() {
	// Add types to scheme
	configv1.AddToScheme(scheme.Scheme)
	machinev1.AddToScheme(scheme.Scheme)
}

func TestNewMachineScope(t *testing.T) {
	g := NewWithT(t)

	// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	// User Data Secret
	userDataSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespaceName,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte("userDataBlob"),
		},
	}

	// Credentials Secret
	credentialsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credentialsSecretName,
			Namespace: defaultNamespaceName,
		},
		Data: map[string][]byte{
			credentialsSecretKey: []byte("{\"ibmcloud_api_key\": \"test\"}"),
		},
	}

	fakeClient := controllerfake.NewFakeClient(userDataSecret, credentialsSecret)

	providerSpec, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(&ibmcloudproviderv1.IBMCloudMachineProviderSpec{
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: credentialsSecretName,
		},
	})

	g.Expect(err).ToNot(HaveOccurred())

	ibmClientBuilder := func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
		return mockIBMClient, nil
	}
	invalidIbmClientBuilder := func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
		return nil, errors.New("ibmc test error")
	}

	cases := []struct {
		name          string
		params        machineScopeParams
		expectedError error
	}{
		{
			name: "successfully create machine scope",
			params: machineScopeParams{
				client:           fakeClient,
				ibmClientBuilder: ibmClientBuilder,
				machine: &machinev1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ibmc-test",
						Namespace: defaultNamespaceName,
						Labels: map[string]string{
							machinev1.MachineClusterIDLabel: "CLUSTERID",
						},
					},
					Spec: machinev1.MachineSpec{
						ProviderSpec: machinev1.ProviderSpec{
							Value: providerSpec,
						},
					},
				},
			},
		},
		{
			name: "fail to get provider spec",
			params: machineScopeParams{
				client:           fakeClient,
				ibmClientBuilder: ibmClientBuilder,
				machine: &machinev1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ibmc-test",
						Namespace: defaultNamespaceName,
						Labels: map[string]string{
							machinev1.MachineClusterIDLabel: "CLUSTERID",
						},
					},
					Spec: machinev1.MachineSpec{
						ProviderSpec: machinev1.ProviderSpec{
							Value: &runtime.RawExtension{
								Raw: []byte{'1'},
							},
						},
					},
				},
			},
			expectedError: errors.New("failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal number into Go value of type v1.IBMCloudMachineProviderSpec"),
		},
		{
			name: "fail to get provider status",
			params: machineScopeParams{
				client:           fakeClient,
				ibmClientBuilder: ibmClientBuilder,
				machine: &machinev1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ibmc-test",
						Namespace: defaultNamespaceName,
						Labels: map[string]string{
							machinev1.MachineClusterIDLabel: "CLUSTERID",
						},
					},
					Spec: machinev1.MachineSpec{
						ProviderSpec: machinev1.ProviderSpec{
							Value: providerSpec,
						},
					},
					Status: machinev1.MachineStatus{
						ProviderStatus: &runtime.RawExtension{
							Raw: []byte{'1'},
						},
					},
				},
			},
			expectedError: errors.New("failed to get machine provider status: error unmarshalling providerStatus: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal number into Go value of type v1.IBMCloudMachineProviderStatus"),
		},
		{
			name: "fail to get credentials secret",
			params: machineScopeParams{
				client:           controllerfake.NewFakeClient(userDataSecret),
				ibmClientBuilder: ibmClientBuilder,
				machine: &machinev1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ibmc-test",
						Namespace: defaultNamespaceName,
						Labels: map[string]string{
							machinev1.MachineClusterIDLabel: "CLUSTERID",
						},
					},
					Spec: machinev1.MachineSpec{
						ProviderSpec: machinev1.ProviderSpec{
							Value: providerSpec,
						},
					},
				},
			},
			expectedError: errors.New("error getting credentials secret \"test-ic-credentials\" in namespace \"test-ns\": secrets \"test-ic-credentials\" not found"),
		},
		{
			name: "fail to create ibm client",
			params: machineScopeParams{
				client:           fakeClient,
				ibmClientBuilder: invalidIbmClientBuilder,
				machine: &machinev1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ibmc-test",
						Namespace: defaultNamespaceName,
						Labels: map[string]string{
							machinev1.MachineClusterIDLabel: "CLUSTERID",
						},
					},
					Spec: machinev1.MachineSpec{
						ProviderSpec: machinev1.ProviderSpec{
							Value: providerSpec,
						},
					},
				},
			},
			expectedError: errors.New("error creating ibm client: ibmc test error"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gs := NewWithT(t)
			scope, err := newMachineScope(tc.params)

			if tc.expectedError != nil {
				gs.Expect(err).To(HaveOccurred())
				gs.Expect(err.Error()).To(Equal(tc.expectedError.Error()))
			} else {
				gs.Expect(err).ToNot(HaveOccurred())
				gs.Expect(scope.Context).To(Equal(context.Background()))
			}
		})
	}

}

func TestPatchMachine(t *testing.T) {
	g := NewWithT(t)

	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "vendor", "github.com", "openshift", "api", "machine", "v1beta1")},
	}

	cfg, err := testEnv.Start()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cfg).ToNot(BeNil())
	defer func() {
		g.Expect(testEnv.Stop()).To(Succeed())
	}()

	k8sClient, err := client.New(cfg, client.Options{})
	g.Expect(err).ToNot(HaveOccurred())

	// Namespace
	defaultNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultNamespaceName,
		},
	}
	g.Expect(k8sClient.Create(context.Background(), defaultNamespace)).To(Succeed())
	defer func() {
		g.Expect(k8sClient.Delete(context.Background(), defaultNamespace)).To(Succeed())
	}()

	// UserDataSecret
	userDataSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespaceName,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte("userDataBlob"),
		},
	}
	g.Expect(k8sClient.Create(context.Background(), userDataSecret)).To(Succeed())
	defer func() {
		g.Expect(k8sClient.Delete(context.Background(), userDataSecret)).To(Succeed())
	}()

	// CredentialsSecret
	credentialsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credentialsSecretName,
			Namespace: defaultNamespaceName,
		},
		Data: map[string][]byte{
			credentialsSecretKey: []byte("{\"ibmcloud_api_key\": \"test\"}"),
		},
	}
	g.Expect(k8sClient.Create(context.Background(), credentialsSecret)).To(Succeed())
	defer func() {
		g.Expect(k8sClient.Delete(context.Background(), credentialsSecret)).To(Succeed())
	}()

	failedPhase := "Failed"

	testCases := []struct {
		name   string
		mutate func(*machinev1.Machine)
		expect func(*machinev1.Machine) error
	}{
		{
			name: "Test changing labels",
			mutate: func(m *machinev1.Machine) {
				if m.Labels == nil {
					m.Labels = make(map[string]string)
				}
				m.Labels["testlabel"] = "test"
			},
			expect: func(m *machinev1.Machine) error {
				if m.Labels["testlabel"] != "test" {
					return fmt.Errorf("label \"testlabel\" %q not equal expected \"test\"", m.ObjectMeta.Labels["test"])
				}
				return nil
			},
		},
		{
			name: "Test setting phase",
			mutate: func(m *machinev1.Machine) {
				m.Status.Phase = &failedPhase
			},
			expect: func(m *machinev1.Machine) error {
				if m.Status.Phase != nil && *m.Status.Phase == failedPhase {
					return nil
				}
				return fmt.Errorf("phase is nil or not equal expected \"Failed\"")
			},
		},
		{
			name: "Test setting provider status",
			mutate: func(m *machinev1.Machine) {
				instanceID := "123"
				instanceState := "running"
				providerStatus, err := ibmcloudproviderv1.RawExtensionFromProviderStatus(&ibmcloudproviderv1.IBMCloudMachineProviderStatus{
					InstanceID:    &instanceID,
					InstanceState: &instanceState,
				})
				if err != nil {
					panic(err)
				}
				m.Status.ProviderStatus = providerStatus
			},
			expect: func(m *machinev1.Machine) error {
				providerStatus, err := ibmcloudproviderv1.ProviderStatusFromRawExtension(m.Status.ProviderStatus)
				if err != nil {
					return fmt.Errorf("unable to get provider status: %v", err)
				}

				if providerStatus.InstanceID == nil || *providerStatus.InstanceID != "123" {
					return fmt.Errorf("instanceID is nil or not equal expected \"123\"")
				}

				if providerStatus.InstanceState == nil || *providerStatus.InstanceState != "running" {
					return fmt.Errorf("instanceState is nil or not equal expected \"running\"")
				}

				return nil
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gs := NewWithT(t)
			timeout := 10 * time.Second

			machine, err := stubMachine()
			gs.Expect(err).ToNot(HaveOccurred())
			gs.Expect(machine).ToNot(BeNil())

			// Create machine
			gs.Expect(k8sClient.Create(context.Background(), machine)).To(Succeed())
			defer func() {
				gs.Expect(k8sClient.Delete(context.Background(), machine)).To(Succeed())
			}()

			// Ensure the machine has synced to the cache
			getMachine := func() error {
				machineKey := types.NamespacedName{Namespace: machine.Namespace, Name: machine.Name}
				return k8sClient.Get(context.Background(), machineKey, machine)
			}

			gs.Eventually(getMachine, timeout).Should(Succeed())

			machineScope, err := newMachineScope(machineScopeParams{
				client:  k8sClient,
				machine: machine,
				ibmClientBuilder: func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return nil, nil
				},
			})

			if err != nil {
				t.Fatal(err)
			}

			tc.mutate(machineScope.machine)

			// close() the machine and check the expectation from the test case
			gs.Expect(machineScope.Close()).To(Succeed())
			checkExpectation := func() error {
				if err := getMachine(); err != nil {
					return err
				}
				return tc.expect(machine)
			}
			gs.Eventually(checkExpectation, timeout).Should(Succeed())

			// Check that resource version doesn't change if we call patchMachine() again
			machineResourceVersion := machine.ResourceVersion

			gs.Expect(machineScope.PatchMachine()).To(Succeed())
			gs.Eventually(getMachine, timeout).Should(Succeed())
			gs.Expect(machine.ResourceVersion).To(Equal(machineResourceVersion))
		})
	}
}
