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
	v1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
	configv1.AddToScheme(scheme.Scheme)
}

var (
	userDataSecretName    = "user-data-test"
	credentialsSecretName = "test-ic-credentials"
	defaultNamespaceName  = "test-ns"
	credentialsSecretKey  = "ibmcloud_api_key"
)

func TestActuatorEvents(t *testing.T) {
	g := NewWithT(t)
	timeout := 10 * time.Second

	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crds")},
	}

	cfg, err := testEnv.Start()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cfg).ToNot(BeNil())
	defer func() {
		g.Expect(testEnv.Stop()).To(Succeed())
	}()

	mgr, err := manager.New(cfg, manager.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	if err != nil {
		t.Fatal(err)
	}

	mgrCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		g.Expect(mgr.Start(mgrCtx)).To(Succeed())
	}()

	// K8s client
	k8sClient := mgr.GetClient()
	eventRecorder := mgr.GetEventRecorderFor("ibmcloudcontroller")

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

	providerSpec, err := v1.RawExtensionFromProviderSpec(&v1.IBMCloudMachineProviderSpec{
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: credentialsSecretName,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(providerSpec).ToNot(BeNil())

	cases := []struct {
		name                string
		error               string
		operation           func(actuator *Actuator, machine *machinev1.Machine)
		event               string
		ibmCloudError       bool
		invalidMachineScope bool
	}{
		{
			name: "Create machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				machine.Spec = machinev1.MachineSpec{
					ProviderSpec: machinev1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: []byte{'1'},
						},
					},
				}
				actuator.Create(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: failed to create scope for machine: failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal number into Go value of type v1.IBMCloudMachineProviderSpec",
			invalidMachineScope: true,
			ibmCloudError:       false,
		},
		{
			name: "Create machine event failed, reconciler's create failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				machine.Labels[machinev1.MachineClusterIDLabel] = ""
				actuator.Create(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: reconciler failed to Create machine: failed validating machine provider spec: machine is missing \"machine.openshift.io/cluster-api-cluster\" label",
			invalidMachineScope: false,
			ibmCloudError:       true,
		},
		{
			name: "Create machine event succeed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event:               "Created Machine ibm-actuator-testing-machine",
			invalidMachineScope: false,
			ibmCloudError:       false,
		},
		{
			name: "Update machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				machine.Spec = machinev1.MachineSpec{
					ProviderSpec: machinev1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: []byte{'1'},
						},
					},
				}
				actuator.Update(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: failed to create scope for machine: failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal number into Go value of type v1.IBMCloudMachineProviderSpec",
			invalidMachineScope: true,
			ibmCloudError:       false,
		},
		{
			name: "Update machine event failed, reconciler's update failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				machine.Labels[machinev1.MachineClusterIDLabel] = ""
				actuator.Update(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: reconciler failed to Update machine: failed validating machine provider spec: machine is missing \"machine.openshift.io/cluster-api-cluster\" label",
			invalidMachineScope: false,
			ibmCloudError:       true,
		},
		{
			name: "Update machine event succeed and only one event is created",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.Background(), machine)
				actuator.Update(context.Background(), machine)
			},
			event: "Updated Machine ibm-actuator-testing-machine",
		},
		{
			name: "Delete machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				machine.Spec = machinev1.MachineSpec{
					ProviderSpec: machinev1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: []byte{'1'},
						},
					},
				}
				actuator.Delete(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: failed to create scope for machine: failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal number into Go value of type v1.IBMCloudMachineProviderSpec",
			invalidMachineScope: true,
			ibmCloudError:       false,
		},
		{
			name: "Delete machine event failed, reconciler's delete failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.Background(), machine)
			},
			event:               "ibm-actuator-testing-machine: reconciler failed to Delete machine: requeue in: 20s",
			invalidMachineScope: false,
			ibmCloudError:       true,
		},
		{
			name: "Delete machine event succeed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.Background(), machine)
			},
			event:               "Deleted machine ibm-actuator-testing-machine",
			invalidMachineScope: false,
			ibmCloudError:       false,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIBMClient := mockibm.NewMockClient(mockCtrl)

	gomock.InOrder(
		mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(true, nil),
		mockIBMClient.EXPECT().InstanceExistsByName(gomock.Any(), gomock.Any()).Return(false, nil),
	)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gs := NewWithT(t)

			machine := &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ibm-actuator-testing-machine",
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
			}

			// Create the machine
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

			ibmClientBuilder := func(secretVal string, providerSpec v1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
				return mockIBMClient, nil
			}
			if tc.invalidMachineScope {
				ibmClientBuilder = func(secretVal string, providerSpec v1.IBMCloudMachineProviderSpec) (ibmclient.Client, error) {
					return nil, errors.New("IBM Cloud client error")
				}
			}

			mockIBMClient.EXPECT().InstanceCreate(machine.Name, gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			mockIBMClient.EXPECT().InstanceGetByName(machine.Name, gomock.Any()).Return(stubInstanceGetByName(machine.Name, &v1.IBMCloudMachineProviderSpec{CredentialsSecret: &corev1.LocalObjectReference{Name: credentialsSecretName}})).AnyTimes()
			mockIBMClient.EXPECT().InstanceDeleteByName(machine.Name, gomock.Any()).Return(nil).AnyTimes()

			params := ActuatorParams{
				Client:           k8sClient,
				EventRecorder:    eventRecorder,
				IbmClientBuilder: ibmClientBuilder,
			}

			actuator := NewActuator(params)
			tc.operation(actuator, machine)

			eventList := &corev1.EventList{}
			waitForEvent := func() error {
				err := k8sClient.List(context.Background(), eventList, client.InNamespace(machine.Namespace))
				if err != nil {
					return err
				}

				if len(eventList.Items) != 1 {
					return fmt.Errorf("expected len 1, got %d", len(eventList.Items))
				}
				return nil
			}

			gs.Eventually(waitForEvent, timeout).Should(Succeed())

			gs.Expect(eventList.Items[0].Message).To(Equal(tc.event))

			for i := range eventList.Items {
				gs.Expect(k8sClient.Delete(context.Background(), &eventList.Items[i])).To(Succeed())
			}

		})
	}
}

func TestActuatorExists(t *testing.T) {

	providerSpec, err := v1.RawExtensionFromProviderSpec(&v1.IBMCloudMachineProviderSpec{
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: credentialsSecretName,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name        string
		eventAction string
		event       string
	}{
		{
			name:        "Create event when event action is present",
			eventAction: "testAction",
			event:       "Warning FailedtestAction testError",
		},
		{
			name:        "Don't event when there is no event action",
			eventAction: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				// use fake recorder and store an event into one item long buffer for subsequent check
				EventRecorder: &record.FakeRecorder{
					Events: eventsChannel,
				},
			}

			machine := &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ibm-actuator-testing-machine",
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
			}

			actuator := NewActuator(params)

			actuator.handleMachineError(machine, errors.New("testError"), tc.eventAction)

			select {
			case event := <-eventsChannel:
				if event != tc.event {
					t.Errorf("Expected %q event, got %q", tc.event, event)
				}
			default:
				if tc.event != "" {
					t.Errorf("Expected %q event, got none", tc.event)
				}
			}
		})
	}
}
