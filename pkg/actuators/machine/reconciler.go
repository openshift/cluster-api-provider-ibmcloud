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
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"github.com/openshift/machine-api-operator/pkg/metrics"
	apicorev1 "k8s.io/api/core/v1"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	klog "k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	requeueAfterSeconds      = 20
	requeueAfterFatalSeconds = 180
	userDataSecretKey        = "userData"
)

// Reconciler are list of services required by machine actuator, easy to create a fake
type Reconciler struct {
	*machineScope
}

// NewReconciler populates all the services based on input scope
func newReconciler(scope *machineScope) *Reconciler {
	return &Reconciler{
		scope,
	}
}

// Create creates an instance via machine cr which is handled by cluster-api
func (r *Reconciler) create() error {

	if err := validateMachine(*r.machine); err != nil {
		return machinecontroller.InvalidMachineConfiguration("failed validating machine provider spec: %v", err)
	}

	userData, err := r.getUserData()
	if err != nil {
		return fmt.Errorf("failed to get user data: %w", err)
	}

	// Create an instance
	_, err = r.ibmClient.InstanceCreate(r.machine.Name, r.providerSpec, userData)

	if err != nil {
		klog.Errorf("%s: error occured while creating machine: %w", r.machine.Name, err)
		metrics.RegisterFailedInstanceCreate(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})

		if reconcileMachineWithCloudStateErr := r.reconcileMachineWithCloudState(&ibmcloudproviderv1.IBMCloudMachineProviderCondition{
			Type:    ibmcloudproviderv1.MachineCreated,
			Status:  apicorev1.ConditionFalse,
			Reason:  ibmcloudproviderv1.MachineCreationFailed,
			Message: err.Error(),
		}); reconcileMachineWithCloudStateErr != nil {
			klog.Errorf("failed to reconcile machine condtion with cloud state: %v", reconcileMachineWithCloudStateErr)
		}
		return fmt.Errorf("failed to create instance via ibm vpc client: %w", err)
	}

	// Update Machine Spec and status with instance info
	return r.reconcileMachineWithCloudState(nil)
}

// update gets instance details and reconciles the machine resource with its state
func (r *Reconciler) update() error {
	if err := validateMachine(*r.machine); err != nil {
		return machinecontroller.InvalidMachineConfiguration("failed validating machine provider spec: %v", err)
	}

	// Update cloud state
	return r.reconcileMachineWithCloudState(nil)
}

func validateMachine(machine machinev1.Machine) error {
	if machine.Labels[machinev1.MachineClusterIDLabel] == "" {
		return machinecontroller.InvalidMachineConfiguration("machine is missing %q label", machinev1.MachineClusterIDLabel)
	}

	return nil
}

// Returns true if machine exists.
func (r *Reconciler) exists() (bool, error) {
	// check if instance exist
	exist, err := r.ibmClient.InstanceExistsByName(r.machine.GetName(), r.providerSpec)
	return exist, err
}

// delete makes a request to delete an instance
func (r *Reconciler) delete() error {

	// Check if the instance exists
	exists, err := r.exists()
	if err != nil {
		return err
	}

	// Found the instance?
	if !exists {
		klog.Infof("%s: Machine not found during delete, skipping", r.machine.Name)
		return nil
	}

	// Delete the instance
	if err = r.ibmClient.InstanceDeleteByName(r.machine.GetName(), r.providerSpec); err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		return fmt.Errorf("failed to delete instance via ibmClient: %v", err)
	}

	klog.Infof("%s: machine status is exists, requeuing...", r.machine.Name)

	return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
}

// getUserData returns User data ignition config
func (r *Reconciler) getUserData() (string, error) {
	if r.providerSpec == nil || r.providerSpec.UserDataSecret == nil {
		return "", nil
	}

	var userDataSecret apicorev1.Secret

	if err := r.client.Get(context.Background(), client.ObjectKey{Namespace: r.machine.GetNamespace(), Name: r.providerSpec.UserDataSecret.Name}, &userDataSecret); err != nil {
		if apimachineryerrors.IsNotFound(err) {
			return "", machinecontroller.InvalidMachineConfiguration("user data secret %q in namespace %q not found: %v", r.providerSpec.UserDataSecret.Name, r.machine.GetNamespace(), err)
		}
		return "", fmt.Errorf("error getting user data secret %q in namespace %q: %v", r.providerSpec.UserDataSecret.Name, r.machine.GetNamespace(), err)
	}
	data, exists := userDataSecret.Data[userDataSecretKey]
	if !exists {
		return "", machinecontroller.InvalidMachineConfiguration("secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance", r.machine.GetNamespace(), r.providerSpec.UserDataSecret.Name, userDataSecretKey)
	}
	return string(data), nil
}

// reconcileMachineWithCloudState reconcile Machine status and spec with the lastest cloud state
func (r *Reconciler) reconcileMachineWithCloudState(conditionFailed *ibmcloudproviderv1.IBMCloudMachineProviderCondition) error {
	// Update providerStatus.Conditions with the failed condtions
	if conditionFailed != nil {
		r.providerStatus.Conditions = reconcileProviderConditions(r.providerStatus.Conditions, *conditionFailed)
		return nil
	}

	// conditionFailed is nil, get the cloud instance and reconcile the fields
	newInstance, err := r.ibmClient.InstanceGetByName(r.machine.Name, r.providerSpec)
	if err != nil {
		return fmt.Errorf("get instance failed with an error: %q", err)
	}

	// Update Machine Status Addresses
	ipAddr := *newInstance.PrimaryNetworkInterface.PrimaryIpv4Address
	if ipAddr != "" {
		networkAddresses := []apicorev1.NodeAddress{{Type: apicorev1.NodeInternalDNS, Address: r.machine.Name}}
		networkAddresses = append(networkAddresses, apicorev1.NodeAddress{Type: apicorev1.NodeInternalIP, Address: ipAddr})
		r.machine.Status.Addresses = networkAddresses
	} else {
		return fmt.Errorf("could not get the primary ipv4 address of instance: %v", newInstance.Name)
	}

	clusterID := r.machine.Labels[machinev1.MachineClusterIDLabel]
	providerID := fmt.Sprintf("ibmvpc://%s/%s/%s", clusterID, r.providerSpec.Zone, r.machine.GetName())
	currProviderID := r.machine.Spec.ProviderID

	// Provider ID check and update
	if currProviderID != nil && *currProviderID == providerID {
		klog.Infof("%s: provider id already set in the machine Spec with value:%s", r.machine.Name, *currProviderID)
	} else {
		r.machine.Spec.ProviderID = &providerID
		klog.Infof("%s: provider id set at machine spec: %s", r.machine.Name, providerID)
	}

	// Set providerStatus in machine
	r.providerStatus.InstanceState = newInstance.Status
	r.providerStatus.InstanceID = newInstance.ID

	// Update conditions
	conditionSuccess := ibmcloudproviderv1.IBMCloudMachineProviderCondition{
		Type:    ibmcloudproviderv1.MachineCreated,
		Reason:  ibmcloudproviderv1.MachineCreationSucceeded,
		Message: machineCreationSucceedMessageCondition,
		Status:  apicorev1.ConditionTrue,
	}
	r.providerStatus.Conditions = reconcileProviderConditions(r.providerStatus.Conditions, conditionSuccess)

	// Update labels & Annotations
	r.setMachineCloudProviderSpecifics(newInstance)

	// Requeue if status is not Running
	if *newInstance.Status != "running" {
		klog.Infof("%s: machine status is %q, requeuing...", r.machine.Name, *newInstance.Status)
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
	}
	return nil
}

// setMachineCloudProviderSpecifics updates Machine resource labels and Annotations
func (r *Reconciler) setMachineCloudProviderSpecifics(instance *vpcv1.Instance) {
	// Make sure machine labels are present before any updates
	if r.machine.Labels == nil {
		r.machine.Labels = make(map[string]string)
	}

	// Make sure machine Annotations are present before any updates
	if r.machine.Annotations == nil {
		r.machine.Annotations = make(map[string]string)
	}

	// Update annotations
	r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = *instance.Status

	// Update labels
	r.machine.Labels[machinecontroller.MachineRegionLabelName] = r.providerSpec.Region
	r.machine.Labels[machinecontroller.MachineAZLabelName] = r.providerSpec.Zone
	r.machine.Labels[machinecontroller.MachineInstanceTypeLabelName] = r.providerSpec.Profile

}
