// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DHCPServerCreate d h c p server create
//
// swagger:model DHCPServerCreate
type DHCPServerCreate struct {

	// Optional cidr for DHCP private network
	Cidr *string `json:"cidr,omitempty"`

	// Optional cloud connection uuid to connect with DHCP private network
	CloudConnectionID *string `json:"cloudConnectionID,omitempty"`

	// Optional DNS Server for DHCP service
	DNSServer *string `json:"dnsServer,omitempty"`

	// Optional name of DHCP Service. Only alphanumeric characters and dashes are allowed (will be prefixed by DHCP identifier)
	Name *string `json:"name,omitempty"`

	// Optional network security groups that the DHCP server network interface is a member of. There is a limit of 1 network security group in the array. If not specified, default network security group is used.
	NetworkSecurityGroupIDs []string `json:"networkSecurityGroupIDs"`

	// Indicates if SNAT will be enabled for DHCP service
	SnatEnabled *bool `json:"snatEnabled,omitempty"`
}

// Validate validates this d h c p server create
func (m *DHCPServerCreate) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this d h c p server create based on context it is used
func (m *DHCPServerCreate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DHCPServerCreate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DHCPServerCreate) UnmarshalBinary(b []byte) error {
	var res DHCPServerCreate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
