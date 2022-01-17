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

	machinev1 "github.com/openshift/api/machine/v1beta1"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	"github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/util"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
	controllerRuntimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// machineScopeParams defines the input parameters used to create a new MachineScope.
type machineScopeParams struct {
	context.Context

	client           controllerRuntimeClient.Client          // api server controller runtime client
	machine          *machinev1.Machine                      // machine resource
	ibmClientBuilder ibmclient.IbmcloudClientBuilderFuncType // for building ibmcloud client
}

// machineScope defines a scope defined around a machine and its cluster
type machineScope struct {
	context.Context

	client    controllerRuntimeClient.Client // api server controller runtime client
	ibmClient ibmclient.Client               // client for interacting with IBM Cloud

	// Machine resources
	machine            *machinev1.Machine
	machineToBePatched controllerRuntimeClient.Patch
	providerSpec       *ibmcloudproviderv1.IBMCloudMachineProviderSpec
	providerStatus     *ibmcloudproviderv1.IBMCloudMachineProviderStatus

	// origMachine captures original value of machine before it is updated (to
	// skip object updated if nothing is changed)
	origMachine *machinev1.Machine
	// origProviderStatus captures original value of machine provider status
	// before it is updated (to skip object updated if nothing is changed)
	origProviderStatus *ibmcloudproviderv1.IBMCloudMachineProviderStatus
}

// newMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func newMachineScope(params machineScopeParams) (*machineScope, error) {
	if params.Context == nil {
		params.Context = context.Background()
	}

	providerSpec, err := ibmcloudproviderv1.ProviderSpecFromRawExtension(params.machine.Spec.ProviderSpec.Value)
	if err != nil {
		return nil, machineapierrors.InvalidMachineConfiguration("failed to get machine config: %v", err)
	}

	providerStatus, err := ibmcloudproviderv1.ProviderStatusFromRawExtension(params.machine.Status.ProviderStatus)
	if err != nil {
		return nil, machineapierrors.InvalidMachineConfiguration("failed to get machine provider status: %v", err.Error())
	}

	apikey, err := util.GetCredentialsSecret(params.client, params.machine.GetNamespace(), *providerSpec)
	if err != nil {
		return nil, err
	}

	ibmClient, err := params.ibmClientBuilder(apikey, *providerSpec)
	if err != nil {
		return nil, machineapierrors.InvalidMachineConfiguration("error creating ibm client: %v", err.Error())
	}

	return &machineScope{
		Context:   params.Context,
		client:    params.client,
		ibmClient: ibmClient,
		// Deep copy the machine since it is changed outside
		// of the machine scope by consumers of the machine
		// scope (e.g. reconciler).
		machine:            params.machine.DeepCopy(),
		providerSpec:       providerSpec,
		providerStatus:     providerStatus,
		origMachine:        params.machine.DeepCopy(),
		origProviderStatus: providerStatus.DeepCopy(),
		machineToBePatched: controllerRuntimeClient.MergeFrom(params.machine.DeepCopy()),
	}, nil
}

// Close the MachineScope by persisting the machine spec, machine status after reconciling.
func (s *machineScope) Close() error {
	if err := s.setMachineStatus(); err != nil {
		return fmt.Errorf("[machinescope] failed to set provider status for machine %q in namespace %q: %v", s.machine.Name, s.machine.Namespace, err)
	}

	if err := s.setMachineSpec(); err != nil {
		return fmt.Errorf("[machinescope] failed to set machine spec %q in namespace %q: %v", s.machine.Name, s.machine.Namespace, err)
	}

	if err := s.PatchMachine(); err != nil {
		return fmt.Errorf("[machinescope] failed to patch machine %q in namespace %q: %v", s.machine.Name, s.machine.Namespace, err)
	}

	return nil
}

func (s *machineScope) setMachineSpec() error {
	ext, err := ibmcloudproviderv1.RawExtensionFromProviderSpec(s.providerSpec)
	if err != nil {
		return err
	}

	klog.V(4).Infof("Storing machine spec for %q, resourceVersion: %v, generation: %v", s.machine.Name, s.machine.ResourceVersion, s.machine.Generation)
	s.machine.Spec.ProviderSpec.Value = ext

	return nil
}

func (s *machineScope) setMachineStatus() error {
	if equality.Semantic.DeepEqual(s.providerStatus, s.origProviderStatus) && equality.Semantic.DeepEqual(s.machine.Status.Addresses, s.origMachine.Status.Addresses) {
		klog.Infof("%s: status unchanged", s.machine.Name)
		return nil
	}

	klog.V(4).Infof("Storing machine status for %q, resourceVersion: %v, generation: %v", s.machine.Name, s.machine.ResourceVersion, s.machine.Generation)
	ext, err := ibmcloudproviderv1.RawExtensionFromProviderStatus(s.providerStatus)
	if err != nil {
		return err
	}

	s.machine.Status.ProviderStatus = ext
	time := metav1.Now()
	s.machine.Status.LastUpdated = &time

	return nil
}

func (s *machineScope) PatchMachine() error {
	klog.V(3).Infof("%q: patching machine", s.machine.GetName())

	statusCopy := *s.machine.Status.DeepCopy()

	// patch machine
	if err := s.client.Patch(s.Context, s.machine, s.machineToBePatched); err != nil {
		klog.Errorf("Failed to patch machine %q: %v", s.machine.GetName(), err)
		return err
	}

	s.machine.Status = statusCopy

	// patch status
	if err := s.client.Status().Patch(s.Context, s.machine, s.machineToBePatched); err != nil {
		klog.Errorf("Failed to patch machine status %q: %v", s.machine.GetName(), err)
		return err
	}

	return nil
}
