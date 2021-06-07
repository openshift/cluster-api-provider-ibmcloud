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
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	klog "k8s.io/klog/v2"
)

const (
	requeueAfterSeconds      = 20
	requeueAfterFatalSeconds = 180
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

func (r *Reconciler) create() error {
	klog.Infof("%s: creating machine", r.machine.Name)

	if err := validateMachine(*r.machine, *r.providerSpec); err != nil {
		return machinecontroller.InvalidMachineConfiguration("failed validating machine provider spec: %v", err)
	}

	return r.requeueIfInstancePending()
}

func (r *Reconciler) update() error {
	return nil
}

// Returns true if machine exists.
func (r *Reconciler) exists() (bool, error) {
	return true, nil
}

// delete deletes machine
func (r *Reconciler) delete() error {
	return nil
}

func validateMachine(machine machinev1.Machine, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) error {

	return nil
}

func (r *Reconciler) requeueIfInstancePending() error {
	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	return nil
}
