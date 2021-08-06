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

package machineset

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	mockibm "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client/mock"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

func TestReconcile(t *testing.T) {
	testCases := []struct {
		name                string
		profile             string
		existingAnnotations map[string]string
		expectedAnnotations map[string]string
		expectErr           bool
		ibmClient           func(ctrl *gomock.Controller) ibmclient.Client
	}{
		{
			name:                "with no instance profile set",
			profile:             "",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: make(map[string]string),
			// Expect no error and only log entry in such case as we don't update
			// instance profile dynamically
			expectErr: false,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist"))
				return mockIBMClient
			},
		},
		{
			name:                "with a cx2d-32x64",
			profile:             "cx2d-32x64",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				cpuKey:    "32",
				memoryKey: "65536",
			},
			expectErr: false,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil)
				return mockIBMClient
			},
		},
		{
			name:                "with a mx2-96x768",
			profile:             "mx2-96x768",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				cpuKey:    "96",
				memoryKey: "786432",
			},
			expectErr: false,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil)
				return mockIBMClient
			},
		},
		{
			name:    "with existing annotations",
			profile: "bx2d-4x16",
			existingAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			expectedAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
				cpuKey:     "4",
				memoryKey:  "16384",
			},
			expectErr: false,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil)
				return mockIBMClient
			},
		},
		{
			name:    "with an invalid instance profile",
			profile: "invalid",
			existingAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			expectedAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			// Expect no error and only log entry in such case as we don't update
			// instance profile dynamically
			expectErr: false,
			ibmClient: func(ctrl *gomock.Controller) ibmclient.Client {
				mockCtrl := gomock.NewController(t)
				mockIBMClient := mockibm.NewMockClient(mockCtrl)
				mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist"))
				return mockIBMClient
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			g := NewWithT(tt)

			machineSet, err := newTestMachineSet("default", tc.profile, tc.existingAnnotations)
			g.Expect(err).ToNot(HaveOccurred())

			mockCtrl := gomock.NewController(t)
			r := Reconciler{
				recorder: record.NewFakeRecorder(1),
				getIbmClient: func(_ string, _ ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return tc.ibmClient(mockCtrl), nil
				},
			}

			_, err = r.reconcile(machineSet)
			g.Expect(err != nil).To(Equal(tc.expectErr))
			g.Expect(machineSet.Annotations).To(Equal(tc.expectedAnnotations))
		})
	}
}
func newTestMachineSet(namespace string, profile string, existingAnnotations map[string]string) (*machinev1.MachineSet, error) {
	// Copy anntotations map so we don't modify the input
	annotations := make(map[string]string)
	for k, v := range existingAnnotations {
		annotations[k] = v
	}

	machineProviderSpec := &ibmcloudproviderv1.IBMCloudMachineProviderSpec{
		Profile: profile,
	}
	providerSpec, err := providerSpecFromMachine(machineProviderSpec)
	if err != nil {
		return nil, err
	}

	return &machinev1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Annotations:  annotations,
			GenerateName: "test-machineset-",
			Namespace:    namespace,
		},
		Spec: machinev1.MachineSetSpec{
			Template: machinev1.MachineTemplateSpec{
				Spec: machinev1.MachineSpec{
					ProviderSpec: providerSpec,
				},
			},
		},
	}, nil
}

func providerSpecFromMachine(in *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (machinev1.ProviderSpec, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return machinev1.ProviderSpec{}, err
	}
	return machinev1.ProviderSpec{
		Value: &runtime.RawExtension{Raw: bytes},
	}, nil
}
