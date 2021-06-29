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
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

const (
	machineCreationSucceedReasonCondition  = "MachineCreationSucceeded"
	machineCreationSucceedMessageCondition = "Machine successfully created"
	machineCreationFailedReasonCondition   = "MachineCreationFailed"
)

// reconcileProviderConditions updates condition for the machine and returns new Condition []
func reconcileProviderConditions(conditions []ibmcloudproviderv1.IBMCloudMachineProviderCondition, newCondition ibmcloudproviderv1.IBMCloudMachineProviderCondition) []ibmcloudproviderv1.IBMCloudMachineProviderCondition {
	currTime := metav1.Now()

	currCondition := conditionTypeCheck(conditions, newCondition.Type)

	if currCondition == nil {
		klog.Infof("Adding new provider condition %v", newCondition)
		conditions = append(conditions,
			ibmcloudproviderv1.IBMCloudMachineProviderCondition{
				Type:               newCondition.Type,
				Status:             newCondition.Status,
				Reason:             newCondition.Reason,
				Message:            newCondition.Message,
				LastTransitionTime: currTime,
				LastProbeTime:      currTime,
			},
		)
	} else {
		// Update if new Status is diff from existing Status
		// Update if new Message is diff from existing Message
		// Update if new Reason is diff from existing Reason
		if currCondition.Status != newCondition.Status || currCondition.Message != newCondition.Message || currCondition.Reason != newCondition.Reason {
			klog.Infof("Updating provider condition %v", newCondition)
			currCondition.Status = newCondition.Status
			currCondition.Reason = newCondition.Reason
			currCondition.Message = newCondition.Message
			currCondition.LastProbeTime = currTime

			// Update LastTransitionTime if Status differ
			if currCondition.Status != newCondition.Status {
				currCondition.LastTransitionTime = currTime
			}
		}
	}
	return conditions
}

// conditionTypeCheck checks if new condition is present in conditions [], if not return nil
func conditionTypeCheck(conditions []ibmcloudproviderv1.IBMCloudMachineProviderCondition, conditionType ibmcloudproviderv1.IBMCloudMachineProviderConditionType) *ibmcloudproviderv1.IBMCloudMachineProviderCondition {
	for idx, eachCondition := range conditions {
		if eachCondition.Type == conditionType {
			return &conditions[idx]
		}
	}
	return nil
}
