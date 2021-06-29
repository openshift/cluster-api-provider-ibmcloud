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
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
)

// Client is a wrapper object for IBM SDK clients
type Client interface {
	// Instances functions
	InstanceGetByID(instanceID string) (*vpcv1.Instance, error)
	InstanceExistsByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (bool, error)
	InstanceGetByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (*vpcv1.Instance, error)
	InstanceDeleteByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) error
	InstanceCreate(machineName string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec, userData string) (*vpcv1.Instance, error)

	// Helper functions
	GetCustomImageByName(imageName string, resourceGroupID string) (string, error)
	GetVPCIDByName(vpcName string, resourceGroupID string) (string, error)
	GetResourceGroupIDByName(resourceGroupName string) (string, error)
	GetSubnetIDbyName(subnetName string, resourceGroupID string) (string, error)
	GetSecurityGroupsByName(securityGroupNames []string, resourceGroupID string, vpcID string) ([]vpcv1.SecurityGroupIdentityIntf, error)
}

// ibmCloudClient makes call to IBM Cloud APIs
type ibmCloudClient struct {
	vpcService             *vpcv1.VpcV1
	resourceManagerService *resourcemanagerv2.ResourceManagerV2
}

// IbmcloudClientBuilderFuncType is function type for building ibm cloud client
type IbmcloudClientBuilderFuncType func(credentialVal string) (Client, error)

// NewClient initilizes a new validated client
func NewClient(credentialVal string) (Client, error) {

	// authenticator
	authenticator := &core.IamAuthenticator{
		ApiKey: credentialVal,
	}

	// IC Virtual Private Cloud (VPC) API
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, err
	}

	// IC Resource Manager API
	resourceManagerService, err := resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, err
	}

	return &ibmCloudClient{
		vpcService:             vpcService,
		resourceManagerService: resourceManagerService,
	}, nil
}

// InstanceExistsByName checks if the instance exist in VPC
func (c *ibmCloudClient) InstanceExistsByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (bool, error) {
	// Get Instance info
	_, err := c.InstanceGetByName(name, machineProviderConfig)

	// Instance found
	if err == nil {
		return true, nil
	}

	// Instance not found
	if err.Error() == "Instance not found" {
		return false, nil
	}

	// Could not retrieve Instances list
	return false, err
}

// InstanceDeleteByName deletes the requested instance
func (c *ibmCloudClient) InstanceDeleteByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) error {
	// Get Instance info
	getInstance, err := c.InstanceGetByName(name, machineProviderConfig)
	if err != nil {
		return err
	}

	// Get instance ID
	instanceID := *getInstance.ID
	if instanceID == "" {
		return fmt.Errorf("Could not get the Instance ID")
	}

	// Initialize New Delete Instance Options
	deleteInstanceOption := c.vpcService.NewDeleteInstanceOptions(instanceID)
	// // Set Instance ID
	// deleteInstanceOption.SetID(instanceID)

	// Delete the Instance
	_, err = c.vpcService.DeleteInstance(deleteInstanceOption)
	if err != nil {
		return err
	}

	return nil
}

// InstanceGetByName retrieves a single instance specified by Instance Name
func (c *ibmCloudClient) InstanceGetByName(name string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec) (*vpcv1.Instance, error) {
	// Get region info
	regionName := machineProviderConfig.Region
	region, _, err := c.vpcService.GetRegion(c.vpcService.NewGetRegionOptions(regionName))
	if err != nil {
		return nil, err
	}

	// Set the Service URL
	err = c.vpcService.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
	if err != nil {
		return nil, err
	}

	// Get Service URL
	serviceURL := c.vpcService.GetServiceURL()
	// Initialize New List Instances Options
	listInstOptions := c.vpcService.NewListInstancesOptions()
	// Set Image Name
	listInstOptions.SetName(name)
	// Set VPC Name
	vpcName := machineProviderConfig.VPC
	listInstOptions.SetVPCName(vpcName)

	// Get Instances list
	instance, _, err := c.vpcService.ListInstances(listInstOptions)
	if err != nil {
		return nil, err
	}

	// Check if instance is not nil
	if instance != nil {
		for _, eachInstance := range instance.Instances {
			if name == *eachInstance.Name {
				return &eachInstance, nil
			}
		}
		return nil, fmt.Errorf("Instance not found")
	}

	return nil, fmt.Errorf("Could not retrieve a list of instances - Name: %v in Region: %v under VPC: %v. Service URL: %v", name, regionName, vpcName, serviceURL)
}

// InstanceGetByID returns retrieves a single instance specified by instanceID
func (c *ibmCloudClient) InstanceGetByID(instanceID string) (*vpcv1.Instance, error) {
	options := c.vpcService.NewGetInstanceOptions(instanceID)

	instance, _, err := c.vpcService.GetInstance(options)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// InstanceCreate creates an instance in VPC
func (c *ibmCloudClient) InstanceCreate(machineName string, machineProviderConfig *ibmcloudproviderv1.IBMCloudMachineProviderSpec, userData string) (*vpcv1.Instance, error) {
	// Get Image ID from Image name
	// Get Subnet ID from Subnet name
	// Get SecurityGroups ID from Security Groups name
	// Get VPC ID from VPC name

	// Get region info
	regionName := machineProviderConfig.Region
	region, _, err := c.vpcService.GetRegion(c.vpcService.NewGetRegionOptions(regionName))
	if err != nil {
		return nil, err
	}

	// Set the Service URL
	err = c.vpcService.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
	if err != nil {
		return nil, err
	}

	// Get Resource Group ID
	resourceGroupName := machineProviderConfig.ResourceGroup
	resourceGroupID, err := c.GetResourceGroupIDByName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get Custom Image ID
	imageID, err := c.GetCustomImageByName(machineProviderConfig.Image, resourceGroupID)
	if err != nil {
		return nil, err
	}

	// Get VPC ID
	vpcName := machineProviderConfig.VPC
	vpcID, err := c.GetVPCIDByName(vpcName, resourceGroupID)
	if err != nil {
		return nil, err
	}

	// Get Subnet ID
	subnetName := machineProviderConfig.PrimaryNetworkInterface.Subnet
	subnetID, err := c.GetSubnetIDbyName(subnetName, resourceGroupID)
	if err != nil {
		return nil, err
	}

	// Get Security Groups
	securityGroups, err := c.GetSecurityGroupsByName(machineProviderConfig.PrimaryNetworkInterface.SecurityGroups, resourceGroupID, vpcID)
	if err != nil {
		return nil, err
	}

	// Create Instance Options
	options := &vpcv1.CreateInstanceOptions{}

	// Set Instance Prototype - Contains all the info necessary to provision an instance
	options.SetInstancePrototype(&vpcv1.InstancePrototype{
		Name: &machineName,
		Image: &vpcv1.ImageIdentity{
			ID: &imageID,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &machineProviderConfig.Profile,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &machineProviderConfig.Zone,
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroupID,
		},
		PrimaryNetworkInterface: &vpcv1.NetworkInterfacePrototype{
			Subnet: &vpcv1.SubnetIdentity{
				ID: &subnetID,
			},
			SecurityGroups: securityGroups,
		},
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
		UserData: &userData,
	})

	// Create a new Instance from an instance prototype object
	instance, _, err := c.vpcService.CreateInstance(options)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// GetVPCIDByName Retrives VPC ID
func (c *ibmCloudClient) GetVPCIDByName(vpcName string, resourceGroupID string) (string, error) {
	// Initialize List Vpcs Options
	vpcOptions := c.vpcService.NewListVpcsOptions()

	// Set Resource Group ID
	vpcOptions.SetResourceGroupID(resourceGroupID)

	// Get a list all VPCs
	vpcList, _, err := c.vpcService.ListVpcs(vpcOptions)
	if err != nil {
		return "", err
	}

	if vpcList != nil {
		var vpcID string
		for _, eachVPC := range vpcList.Vpcs {
			if *eachVPC.Name == vpcName {
				vpcID = *eachVPC.ID
				return vpcID, nil
			}
		}
	}

	return "", fmt.Errorf("Could not retrieve VPC ID of name: %v", vpcName)
}

// GetCustomImageByName retrieves custom image from VPC by region and name
func (c *ibmCloudClient) GetCustomImageByName(imageName string, resourceGroupID string) (string, error) {
	// Initialize List Images Options
	options := c.vpcService.NewListImagesOptions()

	// Private images
	options.SetVisibility(vpcv1.ImageVisibilityPrivateConst)
	// Set Resource Group ID
	options.SetResourceGroupID(resourceGroupID)
	// Set Image name
	options.SetName(imageName)

	// List of all the private images in a region
	privateImages, _, err := c.vpcService.ListImages(options)
	if err != nil {
		return "", err
	}

	if privateImages != nil {
		// Update imageID when found a name match
		var imageID string
		for _, eachImage := range privateImages.Images {
			if *eachImage.Name == imageName && *eachImage.Status == vpcv1.ImageStatusAvailableConst {
				// Get Image ID
				imageID = *eachImage.ID
				return imageID, nil
			}
		}
	}

	return "", fmt.Errorf("Could not retrieve Image ID of name: %v", imageName)
}

// GetResourceGroupIDByName retrives a Resource Group ID
func (c *ibmCloudClient) GetResourceGroupIDByName(resourceGroupName string) (string, error) {
	// Get List of Resource Groups
	resourceGroupList, _, err := c.resourceManagerService.ListResourceGroups(c.resourceManagerService.NewListResourceGroupsOptions())
	if err != nil {
		return "", err
	}

	// Check if resourceGroupList is not nil, in case of a 502, etc
	if resourceGroupList != nil {
		var resourceGroupID string
		for _, eachResource := range resourceGroupList.Resources {
			if *eachResource.Name == resourceGroupName {
				// Get Resource Group ID
				resourceGroupID = *eachResource.ID
				return resourceGroupID, nil
			}
		}
	}
	return "", fmt.Errorf("Could not retrieve Resource Group ID of name: %v", resourceGroupName)
}

// GetSubnetIDbyName retrives a Subnet ID
func (c *ibmCloudClient) GetSubnetIDbyName(subnetName string, resourceGroupID string) (string, error) {
	// Initialize List Subnets Options
	subnetOption := c.vpcService.NewListSubnetsOptions()

	// Set Resource Group ID
	subnetOption.SetResourceGroupID(resourceGroupID)

	// Get a list of all subnets
	subnetList, _, err := c.vpcService.ListSubnets(subnetOption)
	if err != nil {
		return "", err
	}

	if subnetList != nil {
		var subnetID string
		for _, eachSubnet := range subnetList.Subnets {
			if *eachSubnet.Name == subnetName {
				// Get Subnet ID
				subnetID = *eachSubnet.ID
				return subnetID, nil
			}
		}
	}
	return "", fmt.Errorf("Could not retrieve Subnet ID of name: %v", subnetName)
}

// GetSecurityGroupsByName retrieves Security Groups ID
func (c *ibmCloudClient) GetSecurityGroupsByName(securityGroupNames []string, resourceGroupID string, vpcID string) ([]vpcv1.SecurityGroupIdentityIntf, error) {
	// Initialize a map with Security Group Names
	securityGroupMap := map[string]string{}
	for _, item := range securityGroupNames {
		securityGroupMap[item] = ""
	}

	// Initialize List Security Groups Options
	securityGroupOptions := c.vpcService.NewListSecurityGroupsOptions()
	// Set Resource Group ID
	securityGroupOptions.SetResourceGroupID(resourceGroupID)
	// Set VPC ID
	securityGroupOptions.SetVPCID(vpcID)

	// Get a List of Security Groups
	securityGroups, _, _ := c.vpcService.ListSecurityGroups(securityGroupOptions)

	var SecurityGroupIdentityList = make([]vpcv1.SecurityGroupIdentityIntf, len(securityGroupNames))
	if securityGroups != nil {
		idxCounter := 0
		for _, eachSecurityGroup := range securityGroups.SecurityGroups {
			if _, ok := securityGroupMap[*eachSecurityGroup.Name]; ok {
				SecurityGroupIdentityList[idxCounter] = &vpcv1.SecurityGroupIdentityByID{
					ID: eachSecurityGroup.ID,
				}
				idxCounter++
			}
		}
	}

	// Check if retrived all IDs
	if len(securityGroupNames) == len(SecurityGroupIdentityList) {
		return SecurityGroupIdentityList, nil
	}

	return nil, fmt.Errorf("Could not retrieve Security Group IDs of Names: %v", securityGroupNames)

}
