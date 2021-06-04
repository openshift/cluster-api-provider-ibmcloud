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
	"fmt"

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
	CustomImageGetID(imageName string, regionName string) (*vpcv1.Image, error)
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
	options := c.vpcService.NewGetInstanceOptions(instanceID)

	instance, _, err := c.vpcService.GetInstance(options)
	return instance, err
}

// CustomImageGet retrieves custom image from VPC by region and name
func (c *ibmCloudClient) CustomImageGetID(imageName string, regionName string) (*vpcv1.Image, error) {
	// Get region info
	region, _, err := c.vpcService.GetRegion(c.vpcService.NewGetRegionOptions(regionName))
	if err != nil {
		return nil, err
	}

	options := c.vpcService.NewListImagesOptions()
	// Private images
	options.SetVisibility(vpcv1.ImageVisibilityPrivateConst)

	// Set the Service URL
	err = c.vpcService.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
	if err != nil {
		return nil, err
	}

	// List of all the private images in a region
	privateImages, _, err := c.vpcService.ListImages(options)
	if err != nil {
		return nil, err
	}

	// Update image when found a name match and its status is available
	var image *vpcv1.Image
	for _, eachImage := range privateImages.Images {
		if *eachImage.Name == imageName && *eachImage.Status == vpcv1.ImageStatusAvailableConst {
			image = &eachImage
			return image, nil
		}
	}
	return nil, fmt.Errorf("Image: %s not found in Region: %s or Image may not be available yet", imageName, region)
}
