// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"

	"github.com/openvex/discovery/pkg/discovery/options"
)

// VexProbe abstracts a backend driver. The main goal of a probe is to
// capture the logic to find vex documents from software identifiers and
// other references.
//
// The initial version of the VexProbe interface exposes a FindDocumentsFromPurl()
// method that takes an options struct and a Package URL. The VEX probe captures
// the logic to find VEX data for the purl type.
type VexProbe interface {
	FindDocumentsFromPurl(options.Options, purl.PackageURL) ([]*vex.VEX, error)
	SetOptions(options.Options)
}
