// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"fmt"

	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"

	"github.com/openvex/discovery/pkg/discovery/options"
)

//counterfeiter:generate . agentImplementation

type agentImplementation interface {
	ParsePurl(string) (purl.PackageURL, error)
	GetPackageProbe(options.Options, purl.PackageURL) (VexProbe, error)
	FindDocumentsFromPurl(options.Options, VexProbe, purl.PackageURL) ([]*vex.VEX, error)
}

type defaultAgentImplementation struct{}

// ParsePurl checks if a purl is correctly formed
func (pi *defaultAgentImplementation) ParsePurl(purlString string) (purl.PackageURL, error) {
	p, err := purl.FromString(purlString)
	if err != nil {
		return p, err
	}
	return p, nil
}

// GetPackageProbe returns a PackageProbe for the specified purl type
func (pi *defaultAgentImplementation) GetPackageProbe(opts options.Options, p purl.PackageURL) (VexProbe, error) {
	if p, ok := probers[p.Type]; ok {
		p.SetOptions(opts)
		return p, nil
	}
	return nil, fmt.Errorf("purl type %s not supported", p.Type)
}

// FetchDocuments downloads all OpenVEX documents using the PackageProbe for
// the specified purl.
func (pi *defaultAgentImplementation) FindDocumentsFromPurl(opts options.Options, pkgProbe VexProbe, p purl.PackageURL) ([]*vex.VEX, error) {
	docs, err := pkgProbe.FindDocumentsFromPurl(opts, p)
	if err != nil {
		return nil, fmt.Errorf("looking for documents: %w", err)
	}
	return docs, nil
}
