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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	gtypes "github.com/onsi/gomega/types"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	mockibm "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client/mock"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("Reconciler", func() {
	var c client.Client
	var stopMgr context.CancelFunc
	var fakeRecorder *record.FakeRecorder
	var namespace *corev1.Namespace
	var mockCtrl *gomock.Controller

	// Mock
	mockCtrl = gomock.NewController(GinkgoT())
	// defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	gomock.InOrder(
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist")),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist")),
	)
	BeforeEach(func() {

		mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
		Expect(err).ToNot(HaveOccurred())

		r := Reconciler{
			Client: mgr.GetClient(),
			Log:    log.Log,
			getIbmClient: func(_ string, _ ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
				return mockIBMClient, nil
			},
		}

		// mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil).AnyTimes()
		// gomock.InOrder(
		// 	mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("Not Found")),
		// 	// mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("Not Found")),
		// 	mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		// )
		Expect(err).ToNot(HaveOccurred())
		Expect(r.SetupWithManager(mgr, controller.Options{})).To(Succeed())

		fakeRecorder = record.NewFakeRecorder(1)
		r.recorder = fakeRecorder

		c = mgr.GetClient()
		stopMgr = StartTestManager(mgr)

		namespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "test-machineset-ns-"}}
		Expect(c.Create(ctx, namespace)).To(Succeed())
	})

	AfterEach(func() {
		Expect(deleteMachineSets(c, namespace.Name)).To(Succeed())
		// mockCtrl.Finish()
		stopMgr()
	})

	type reconcileTestCase = struct {
		profile             string
		existingAnnotations map[string]string
		expectedAnnotations map[string]string
		expectedEvents      []string
	}

	DescribeTable("when reconciling MachineSets", func(rtc reconcileTestCase) {

		machineSet, err := newTestMachineSet(namespace.Name, rtc.profile, rtc.existingAnnotations)
		Expect(err).ToNot(HaveOccurred())

		Expect(c.Create(ctx, machineSet)).To(Succeed())
		Eventually(func() map[string]string {
			m := &machinev1.MachineSet{}
			key := client.ObjectKey{Namespace: machineSet.Namespace, Name: machineSet.Name}
			err := c.Get(ctx, key, m)
			if err != nil {
				return nil
			}
			annotations := m.GetAnnotations()
			if annotations != nil {
				return annotations
			}
			// Return an empty map to distinguish between empty annotations and errors
			return make(map[string]string)
		}, timeout).Should(Equal(rtc.expectedAnnotations))

		// Check which event types were sent
		Eventually(fakeRecorder.Events, timeout).Should(HaveLen(len(rtc.expectedEvents)))
		receivedEvents := []string{}
		eventMatchers := []gtypes.GomegaMatcher{}
		for _, eachEvent := range rtc.expectedEvents {
			receivedEvents = append(receivedEvents, <-fakeRecorder.Events)
			eventMatchers = append(eventMatchers, ContainSubstring(fmt.Sprintf(" %s ", eachEvent)))
		}
		Expect(receivedEvents).To(ConsistOf(eventMatchers))
	},
		Entry("with no profile set", reconcileTestCase{
			profile:             "",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: make(map[string]string),
			expectedEvents:      []string{"FailedUpdate"},
		}),
		Entry("with a bx2d-4x16", reconcileTestCase{
			profile:             "bx2d-4x16",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				profileKey: "bx2d-4x16",
			},
			expectedEvents: []string{},
		}),
		Entry("with a bx2-2x8", reconcileTestCase{
			profile:             "bx2-2x8",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				profileKey: "bx2-2x8",
			},
			expectedEvents: []string{},
		}),
		Entry("with existing annotations", reconcileTestCase{
			profile: "bx2-2x8",
			existingAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			expectedAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
				profileKey: "bx2-2x8",
			},
			expectedEvents: []string{},
		}),
		Entry("with an invalid instance profile", reconcileTestCase{
			profile: "invalid",
			existingAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			expectedAnnotations: map[string]string{
				"existing": "annotation",
				"annother": "existingAnnotation",
			},
			expectedEvents: []string{"FailedUpdate"},
		}),
	)
})

func deleteMachineSets(c client.Client, namespaceName string) error {
	machineSets := &machinev1.MachineSetList{}
	err := c.List(ctx, machineSets, client.InNamespace(namespaceName))
	if err != nil {
		return err
	}

	for _, ms := range machineSets.Items {
		err := c.Delete(ctx, &ms)
		if err != nil {
			return err
		}
	}

	Eventually(func() error {
		machineSets := &machinev1.MachineSetList{}
		err := c.List(ctx, machineSets)
		if err != nil {
			return err
		}
		if len(machineSets.Items) > 0 {
			return fmt.Errorf("machineSets not deleted")
		}
		return nil
	}, timeout).Should(Succeed())

	return nil
}

func TestReconcile(t *testing.T) {
	testCases := []struct {
		name                string
		profile             string
		existingAnnotations map[string]string
		expectedAnnotations map[string]string
		expectErr           bool
	}{
		{
			name:                "with no instance profile set",
			profile:             "",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: make(map[string]string),
			// Expect no error and only log entry in such case as we don't update
			// instance profile dynamically
			expectErr: false,
		},
		{
			name:                "with a cx2d-32x64",
			profile:             "cx2d-32x64",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				profileKey: "cx2d-32x64",
			},
			expectErr: false,
		},
		{
			name:                "with a mx2-96x768",
			profile:             "mx2-96x768",
			existingAnnotations: make(map[string]string),
			expectedAnnotations: map[string]string{
				profileKey: "mx2-96x768",
			},
			expectErr: false,
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
				profileKey: "bx2d-4x16",
			},
			expectErr: false,
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
		},
	}
	// Mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	gomock.InOrder(
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist")),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceGetProfile(gomock.Any()).Return(false, errors.New("the provided instance profile ID does not exist")),
	)
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			g := NewWithT(tt)

			machineSet, err := newTestMachineSet("default", tc.profile, tc.existingAnnotations)
			g.Expect(err).ToNot(HaveOccurred())

			r := Reconciler{
				recorder: record.NewFakeRecorder(1),
				getIbmClient: func(_ string, _ ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return mockIBMClient, nil
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
