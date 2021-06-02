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

package client

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// Client is a wrapper object for IBM SDK clients
type Client interface {
	InstanceGet(instanceID string) (*vpcv1.Instance, error)
	// InstanceListAll()
	// InstanceCreate()
	// InstanceDelete()
	// InstanceUpdate()
}

// ibmCloudClient makes call to IBM Cloud APIs
type ibmCloudClient struct {
	vpcService *vpcv1.VpcV1
}

// IbmcloudClientBuilderFuncType is function type for building ibm cloud client
type IbmcloudClientBuilderFuncType func(credentialVal string) (Client, error)

// NewClient initilizes a new validated client
func NewClient(credentialVal string) (Client, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: credentialVal,
	}

	options := &vpcv1.VpcV1Options{
		Authenticator: authenticator,
	}

	vpcService, vpcServiceErr := vpcv1.NewVpcV1(options)

	if vpcServiceErr != nil {
		panic(vpcServiceErr)
	}
	return &ibmCloudClient{
		vpcService: vpcService,
	}, nil
}

// InstanceGet returns retrieves a single instance specified by instanceID
func (c *ibmCloudClient) InstanceGet(instanceID string) (*vpcv1.Instance, error) {
	options := &vpcv1.GetInstanceOptions{}
	options.SetID(instanceID)

	instance, _, err := c.vpcService.GetInstance(options)
	return instance, err
}
