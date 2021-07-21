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
	"fmt"

	"github.com/go-logr/logr"
	ibmclient "github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/client"
	"github.com/openshift/cluster-api-provider-ibmcloud/pkg/actuators/util"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	klog "k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

const (
	// This exposes compute information based on the providerSpec input.
	// This is needed by the autoscaler to foresee upcoming capacity when scaling from zero.
	// https://github.com/openshift/enhancements/pull/186

	profileKey = "machine.openshift.io/profile"
	// cpuKey       = "machine.openshift.io/vCPU"
	// memoryKey    = "machine.openshift.io/memoryMb"
	// bandwidthKey = "machine.openshift.io/bandwidthMbps"
)

// Reconciler - MachineSet Reconciler
type Reconciler struct {
	Client client.Client
	Log    logr.Logger

	recorder record.EventRecorder
	scheme   *runtime.Scheme

	// Allows to mock ibm client for testing
	getIbmClient func(namespace string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error)
}

// SetupWithManager creates a new controller for a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	_, err := ctrl.NewControllerManagedBy(mgr).
		For(&machinev1.MachineSet{}).
		WithOptions(options).
		Build(r)

	if err != nil {
		return fmt.Errorf("failed setting up with a controller manager: %w", err)
	}

	r.recorder = mgr.GetEventRecorderFor("machineset-controller")
	r.scheme = mgr.GetScheme()

	if r.getIbmClient == nil {
		r.getIbmClient = r.getActualIbmClient
	}

	return nil
}

// getActualIbmClient initilizes an ibmclient
func (r *Reconciler) getActualIbmClient(namespace string, providerSpec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {

	apikey, err := util.GetCredentialsSecret(r.Client, namespace, providerSpec)
	if err != nil {
		return nil, err
	}

	ibmClient, err := ibmclient.NewClient(apikey, providerSpec)
	if err != nil {
		return nil, machineapierrors.InvalidMachineConfiguration("error creating ibm client: %v", err.Error())
	}
	return ibmClient, nil
}

// Reconcile implements controller runtime Reconciler interface.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("machineset", req.Name, "namespace", req.Namespace)
	logger.V(3).Info("Reconciling")

	machineSet := &machinev1.MachineSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, machineSet); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Ignore deleted MachineSets, this can happen when foregroundDeletion
	// is enabled
	if !machineSet.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}
	originalMachineSetToPatch := client.MergeFrom(machineSet.DeepCopy())

	result, err := r.reconcile(machineSet)
	if err != nil {
		logger.Error(err, "Failed to reconcile MachineSet")
		r.recorder.Eventf(machineSet, corev1.EventTypeWarning, "ReconcileError", "%v", err)
		// we don't return here so we want to attempt to patch the machine regardless of an error.
	}

	if err := r.Client.Patch(ctx, machineSet, originalMachineSetToPatch); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch machineSet: %v", err)
	}

	if isInvalidConfigurationError(err) {
		// For situations where requeuing won't help we don't return error.
		// https://github.com/kubernetes-sigs/controller-runtime/issues/617
		return result, nil
	}

	return result, err
}

func isInvalidConfigurationError(err error) bool {
	switch t := err.(type) {
	case *machineapierrors.MachineError:
		if t.Reason == machinev1.InvalidConfigurationMachineError {
			return true
		}
	}
	return false
}

func (r *Reconciler) reconcile(machineSet *machinev1.MachineSet) (ctrl.Result, error) {
	providerConfig, err := getproviderConfig(machineSet)
	if err != nil {
		return ctrl.Result{}, machineapierrors.InvalidMachineConfiguration("failed to get providerConfig: %v", err)
	}

	// profile, ok := Profiles[providerConfig.Profile]
	ibmClient, err := r.getIbmClient(machineSet.GetNamespace(), *providerConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = ibmClient.InstanceGetProfile(providerConfig.Profile)

	if err != nil {
		klog.Error("Unable to set annotations: unknown profile: %v", providerConfig.Profile)
		klog.Error("Autoscaling from zero will not work. To fix this, manually populate machine annotations for your instance profile: %v", []string{profileKey})

		// Returning no error to prevent further reconciliation, as user intervention is now required but emit an informational event
		r.recorder.Eventf(machineSet, corev1.EventTypeWarning, "FailedUpdate", "Failed to set autoscaling from zero annotations, instance profile unknown")
		return ctrl.Result{}, nil
	}

	if machineSet.Annotations == nil {
		machineSet.Annotations = make(map[string]string)
	}

	// Set Annotation
	machineSet.Annotations[profileKey] = providerConfig.Profile

	return ctrl.Result{}, nil
}

func getproviderConfig(machineSet *machinev1.MachineSet) (*ibmcloudproviderv1.IBMCloudMachineProviderSpec, error) {
	return ibmcloudproviderv1.ProviderSpecFromRawExtension(machineSet.Spec.Template.Spec.ProviderSpec.Value)
}
