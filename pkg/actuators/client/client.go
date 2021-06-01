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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// Client is a wrapper object for actual IBM SDK clients to allow for easier testing.
type Client interface {
	InstancesGet()
}

// IBMCloudClient struct
type IBMCloudClient struct {
	VPCService *vpcv1.VpcV1
	//APIKey          string
	//IAMEndpoint     string
	//ServiceEndPoint string
}

// IbmcloudClientBuilderFuncType is function type for building ibm cloud client
type IbmcloudClientBuilderFuncType func(serviceAccountJSON string) (Client, error)

// NewClient return a new client
func NewClient() error {
	// var err error
	// c.VPCService, err = vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
	// 	Authenticator: &core.IamAuthenticator{
	// 		ApiKey: apiKey,
	// 		URL:    iamEndpoint,
	// 	},
	// 	URL: svcEndpoint,
	// })

	return nil
}
