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
	"testing"

	"github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func TestShouldUpdateCondition(t *testing.T) {
	testCases := []struct {
		oldCondition v1beta1.IBMCloudMachineProviderCondition
		newCondition v1beta1.IBMCloudMachineProviderCondition
		expected     corev1.ConditionStatus
	}{
		{
			oldCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
			newCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
		},
		{
			oldCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
			newCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo2",
				Message: "bar2",
				Status:  corev1.ConditionTrue,
			},
		},
		{
			oldCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
			newCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo2",
				Message: "New Message",
				Status:  corev1.ConditionFalse,
			},
		},
		{
			oldCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
			newCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "New Reason",
				Message: "New Message",
				Status:  corev1.ConditionTrue,
			},
		},
		{
			oldCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionFalse,
			},
			newCondition: v1beta1.IBMCloudMachineProviderCondition{
				Reason:  "foo",
				Message: "bar",
				Status:  corev1.ConditionTrue,
			},
		},
	}

	for _, tc := range testCases {

		conditions := []v1beta1.IBMCloudMachineProviderCondition{}
		conditions = append(conditions, tc.oldCondition)
		returnCondition := reconcileProviderConditions(conditions, tc.newCondition)

		if returnCondition[0].Reason != tc.newCondition.Reason &&
			returnCondition[0].Message != tc.newCondition.Message &&
			returnCondition[0].Status != tc.newCondition.Status {
			t.Errorf("Expected %v, got %v", tc.newCondition, returnCondition)
		}
	}
}
