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
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang/mock/gomock"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	mockibm "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client/mock"
	"github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/util"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	controllerfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCreate(t *testing.T) {
	// mock calls
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	cases := []struct {
		testcase          string
		labels            map[string]string
		providerSpec      *ibmcloudproviderv1.IBMCloudMachineProviderSpec
		expectedCondition *ibmcloudproviderv1.IBMCloudMachineProviderCondition
		secret            *corev1.Secret
		expectedError     error
	}{
		{
			testcase: "Create machine succeed ",
			expectedCondition: &ibmcloudproviderv1.IBMCloudMachineProviderCondition{
				Type:    ibmcloudproviderv1.MachineCreated,
				Status:  corev1.ConditionTrue,
				Reason:  machineCreationSucceedReasonCondition,
				Message: machineCreationSucceedMessageCondition,
			},
			providerSpec:  &ibmcloudproviderv1.IBMCloudMachineProviderSpec{},
			expectedError: nil,
		},
		{
			testcase: "Fail on invalid missing machine label",
			labels: map[string]string{
				machinev1.MachineClusterIDLabel: "",
			},
			providerSpec:  &ibmcloudproviderv1.IBMCloudMachineProviderSpec{},
			expectedError: errors.New("failed validating machine provider spec: machine is missing \"machine.openshift.io/cluster-api-cluster\" label"),
		},
		{
			testcase: "Fail on bad user data secret",
			providerSpec: &ibmcloudproviderv1.IBMCloudMachineProviderSpec{
				UserDataSecret: &corev1.LocalObjectReference{
					Name: "invalid_name",
				},
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "invalid_name",
				},
				Data: map[string][]byte{
					"bad_api_key": []byte(""),
				},
			},
			expectedError: errors.New("failed to get user data: secret /invalid_name does not have \"userData\" field set. Thus, no user data applied when creating an instance"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.testcase, func(t *testing.T) {
			providerSpec, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(tc.providerSpec)
			if err != nil {
				t.Fatal(err)
			}

			labels := map[string]string{
				machinev1.MachineClusterIDLabel: "CLUSTERID",
			}
			if tc.labels != nil {
				labels = tc.labels
			}

			// if tc.providerSpec != nil {
			// 	providerSpec = tc.providerSpec
			// }

			machine := &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "",
					Namespace: "",
					Labels:    labels,
				},
				Spec: machinev1.MachineSpec{
					ProviderSpec: machinev1.ProviderSpec{
						Value: providerSpec,
					},
				},
			}
			machineScope, err := newMachineScope(machineScopeParams{
				machine: machine,
				client:  controllerfake.NewFakeClient(),
				// providerSpec:   providerSpec,
				// providerStatus: &ibmcloudproviderv1.IBMCloudMachineProviderStatus{},
				ibmClientBuilder: func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return mockIBMClient, nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			mockIBMClient.EXPECT().InstanceCreate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			mockIBMClient.EXPECT().InstanceGetByName(gomock.Any(), gomock.Any()).Return(stubInstanceGetByName(machine.Name, &ibmcloudproviderv1.IBMCloudMachineProviderSpec{CredentialsSecret: &corev1.LocalObjectReference{Name: credentialsSecretName}})).AnyTimes()
			mockIBMClient.EXPECT().InstanceDeleteByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()

			reconciler := newReconciler(machineScope)

			if tc.secret != nil {
				reconciler.client = controllerfake.NewFakeClientWithScheme(scheme.Scheme, tc.secret)
			}

			err = reconciler.create()

			if tc.expectedCondition != nil {
				if reconciler.providerStatus.Conditions[0].Type != tc.expectedCondition.Type {
					t.Errorf("Expected: %s, got %s", tc.expectedCondition.Type, reconciler.providerStatus.Conditions[0].Type)
				}
				if reconciler.providerStatus.Conditions[0].Status != tc.expectedCondition.Status {
					t.Errorf("Expected: %s, got %s", tc.expectedCondition.Status, reconciler.providerStatus.Conditions[0].Status)
				}
				if reconciler.providerStatus.Conditions[0].Reason != tc.expectedCondition.Reason {
					t.Errorf("Expected: %s, got %s", tc.expectedCondition.Reason, reconciler.providerStatus.Conditions[0].Reason)
				}
				if reconciler.providerStatus.Conditions[0].Message != tc.expectedCondition.Message {
					t.Errorf("Expected: %s, got %s", tc.expectedCondition.Message, reconciler.providerStatus.Conditions[0].Message)
				}
			}

			if tc.expectedError != nil {
				if err == nil {
					t.Error("reconciler was expected to return error")
				}
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected: %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("reconciler was not expected to return error: %v", err)
				}
			}

		})
	}
}

func TestExists(t *testing.T) {
	testCases := []struct {
		name          string
		machine       func() *machinev1.Machine
		ibmClient     func(ctrl *gomock.Controller) ibmclient.Client
		expectedError error
		existsResult  bool
	}{
		{
			name: "Found created instance",
			machine: func() *machinev1.Machine {
				machine, err := stubMachine()
				if err != nil {
					t.Fatalf("unable to build stub machine: %v", err)
				}
				machine.Spec.ProviderSpec = machinev1.ProviderSpec{}
				return machine
			},
			existsResult:  true,
			expectedError: nil,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
				return mockIBMClient
			},
		},
		{
			name: "Cannot find machine instance",
			machine: func() *machinev1.Machine {
				machine, err := stubMachine()
				if err != nil {
					t.Fatalf("unable to build stub machine: %v", err)
				}
				machine.Spec.ProviderSpec = machinev1.ProviderSpec{}

				return machine
			},
			existsResult:  false,
			expectedError: nil,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(false, nil)
				return mockIBMClient
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			machineScope, err := newMachineScope(machineScopeParams{
				machine: tc.machine(),
				client:  controllerfake.NewFakeClient(),
				ibmClientBuilder: func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return tc.ibmClient(mockCtrl), nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			reconciler := newReconciler(machineScope)

			exists, err := reconciler.exists()

			if tc.existsResult != exists {
				t.Errorf("expected reconciler tc.Exists() to return: %v, got %v", tc.existsResult, exists)
			}

			if tc.expectedError != nil {

				if err == nil {
					t.Error("reconciler was expected to return error")
				}

				if err.Error() != tc.expectedError.Error() {
					t.Errorf("expected: %v, got %v", tc.expectedError, err)
				}

			} else {
				if err != nil {
					t.Errorf("reconciler was not expected to return error: %v", err)
				}
			}
		})
	}
}

func machineWithSpec(spec *ibmcloudproviderv1.IBMCloudMachineProviderSpec) *machinev1.Machine {
	providerSpec, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(spec)
	if err != nil {
		panic(err)
	}

	return &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibmc-test",
			Namespace: defaultNamespaceName,
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: providerSpec,
			},
		},
	}
}

func TestGetUserData(t *testing.T) {

	providerSpec := &ibmcloudproviderv1.IBMCloudMachineProviderSpec{
		UserDataSecret: &corev1.LocalObjectReference{
			Name: userDataSecretName,
		},
	}
	testCases := []struct {
		testCase         string
		userDataSecret   *corev1.Secret
		providerSpec     *ibmcloudproviderv1.IBMCloudMachineProviderSpec
		expectedUserdata string
		expectError      bool
	}{
		{
			testCase: "all good",
			userDataSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespaceName,
				},
				Data: map[string][]byte{
					userDataSecretKey: []byte("{}"),
				},
			},
			providerSpec:     providerSpec,
			expectedUserdata: "",
			expectError:      false,
		},
		{
			testCase:       "missing secret",
			userDataSecret: nil,
			providerSpec:   providerSpec,
			expectError:    true,
		},
		{
			testCase: "missing key in secret",
			userDataSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespaceName,
				},
				Data: map[string][]byte{
					"userData": []byte("{}"),
				},
			},
			providerSpec: providerSpec,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			clientObjs := []runtime.Object{}

			if tc.userDataSecret != nil {
				clientObjs = append(clientObjs, tc.userDataSecret)
			}

			client := controllerfake.NewFakeClient(clientObjs...)

			// Can't use newMachineScope because it tries to create an API
			// session, and other things unrelated to these tests.
			ms := &machineScope{
				Context:      context.Background(),
				client:       client,
				machine:      machineWithSpec(tc.providerSpec),
				providerSpec: tc.providerSpec,
			}

			userData, err := util.GetCredentialsSecret(ms.client, ms.machine.GetNamespace(), *ms.providerSpec)
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if userData != tc.expectedUserdata {
				t.Errorf("Got: %q, Want: %q", userData, tc.expectedUserdata)
			}
		})
	}
}

func TestReconcileMachineWithCloudState(t *testing.T) {
	// mock calls
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	zone := "test-zone-us-east-1"
	instanceName := "test-Instance-name"
	labels := map[string]string{
		machinev1.MachineClusterIDLabel: "CLUSTERID",
	}

	providerSpec, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(&ibmcloudproviderv1.IBMCloudMachineProviderSpec{
		Zone: zone,
	})
	if err != nil {
		t.Fatal(err)
	}
	machine := &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: "",
			Labels:    labels,
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: providerSpec,
			},
		},
	}
	machineScope, err := newMachineScope(machineScopeParams{
		machine: machine,
		client:  controllerfake.NewFakeClient(),
		// providerSpec:   providerSpec,
		// providerStatus: &ibmcloudproviderv1.IBMCloudMachineProviderStatus{},
		ibmClientBuilder: func(secretVal string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
			return mockIBMClient, nil
		},
	})

	// instance ID returned from stubInstanceGetByName
	intanceID := "0727_xyz-xyz-cccc-aaba-cacdaccad"
	mockIBMClient.EXPECT().InstanceGetByName(gomock.Any(), gomock.Any()).Return(stubInstanceGetByName(machine.Name, &ibmcloudproviderv1.IBMCloudMachineProviderSpec{CredentialsSecret: &corev1.LocalObjectReference{Name: credentialsSecretName}})).AnyTimes()

	expectedNodeAddresses := []corev1.NodeAddress{
		{
			Type:    "InternalDNS",
			Address: instanceName,
		},
		{
			Type:    "InternalIP",
			Address: "10.0.0.1",
		},
	}
	expectedProviderID := fmt.Sprintf("ibmvpc://CLUSTERID/%s/%s", zone, instanceName)

	r := newReconciler(machineScope)
	if err := r.reconcileMachineWithCloudState(nil); err != nil {
		t.Errorf("reconciler was not expected to return error: %v", err)
	}

	if r.machine.Status.Addresses[0] != expectedNodeAddresses[0] {
		t.Errorf("Expected: %s, got: %s", expectedNodeAddresses[0], r.machine.Status.Addresses[0])
	}
	if r.machine.Status.Addresses[1] != expectedNodeAddresses[1] {
		t.Errorf("Expected: %s, got: %s", expectedNodeAddresses[1], r.machine.Status.Addresses[1])
	}

	if *r.machine.Spec.ProviderID != expectedProviderID {
		t.Errorf("Expected: %s, got: %s", expectedProviderID, *r.machine.Spec.ProviderID)
	}
	if *r.providerStatus.InstanceState != "running" {
		t.Errorf("Expected: %s, got: %s", "running", *r.providerStatus.InstanceState)
	}
	if *r.providerStatus.InstanceID != intanceID {
		t.Errorf("Expected: %s, got: %s", intanceID, *r.providerStatus.InstanceID)
	}
}

func TestSetMachineCloudProviderSpecifics(t *testing.T) {
	dummyStatus := "testStatus"
	dummyProfile := "testProfile"
	dummyZone := "testZone"
	dummyRegion := "testRegion"

	r := Reconciler{
		machineScope: &machineScope{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "",
					Namespace: "",
				},
			},
			providerSpec: &ibmcloudproviderv1.IBMCloudMachineProviderSpec{
				Profile: dummyProfile,
				Region:  dummyRegion,
				Zone:    dummyZone,
			},
		},
	}

	instance := &vpcv1.Instance{
		Status: &dummyStatus,
	}

	r.setMachineCloudProviderSpecifics(instance)

	actualInstanceStateAnnotation := r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName]
	if actualInstanceStateAnnotation != *instance.Status {
		t.Errorf("Expected instance state annotation: %v, got: %v", actualInstanceStateAnnotation, *instance.Status)
	}

	actualMachineProfileLabel := r.machine.Labels[machinecontroller.MachineInstanceTypeLabelName]
	if actualMachineProfileLabel != r.providerSpec.Profile {
		t.Errorf("Expected machine type label: %v, got: %v", actualMachineProfileLabel, r.providerSpec.Profile)
	}

	actualMachineRegionLabel := r.machine.Labels[machinecontroller.MachineRegionLabelName]
	if actualMachineRegionLabel != r.providerSpec.Region {
		t.Errorf("Expected machine region label: %v, got: %v", actualMachineRegionLabel, r.providerSpec.Region)
	}

	actualMachineAZLabel := r.machine.Labels[machinecontroller.MachineAZLabelName]
	if actualMachineAZLabel != r.providerSpec.Zone {
		t.Errorf("Expected machine zone label: %v, got: %v", actualMachineAZLabel, r.providerSpec.Zone)
	}
}

func TestDelete(t *testing.T) {
	// mock calls
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	machineScope := machineScope{
		machine: &machinev1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "",
				Namespace: "",
				Labels: map[string]string{
					machinev1.MachineClusterIDLabel: "CLUSTERID",
				},
			},
		},
		client:         controllerfake.NewFakeClient(),
		providerSpec:   &ibmcloudproviderv1.IBMCloudMachineProviderSpec{},
		providerStatus: &ibmcloudproviderv1.IBMCloudMachineProviderStatus{},
		ibmClient:      mockIBMClient,
	}

	reconciler := newReconciler(&machineScope)

	mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
	mockIBMClient.EXPECT().InstanceDeleteByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	if err := reconciler.delete(); err != nil {
		if _, ok := err.(*machinecontroller.RequeueAfterError); !ok {
			t.Errorf("reconciler was not expected to return error: %v", err)
		}
	}
}
