// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package discovery

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"
	"sync"

	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"

	"github.com/openvex/discovery/pkg/probers/oci"

	"github.com/openvex/discovery/pkg/discovery/options"
)

var (
	regMtx  sync.RWMutex
	probers = map[string]VexProbe{}
)

func init() {
	RegisterBuiltInDrivers()
}

// RegisterBuiltInDrivers adds all the built in backend drivers to the
// VexProbers collection.
func RegisterBuiltInDrivers() {
	RegisterDriver(purl.TypeOCI, oci.New())
}

// RegisterDriver adds a new VexProbe to the drivers collection.
func RegisterDriver(purlType string, probe VexProbe) {
	regMtx.Lock()
	probers[purlType] = probe
	regMtx.Unlock()
}

// UnregisterDrivers removes all registered backend drivers from the probers
// collection. This is useful when you want to enable only specific drivers
// or register custom ones.
func UnregisterDrivers() {
	regMtx.Lock()
	probers = map[string]VexProbe{}
	regMtx.Unlock()
}

// Probe is the main object that inspects repositories and looks for security
// documents. To create a new Probe use the `NewProbe` function
type Agent struct {
	impl    agentImplementation
	Options options.Options
}

// NewAgent creates a new discovery agent
func NewAgent() *Agent {
	return &Agent{
		impl:    &defaultAgentImplementation{},
		Options: options.Default,
	}
}

func (agent *Agent) SetImplementation(impl agentImplementation) {
	agent.impl = impl
}

// ProbePURL examines an PackageURL and retrieves all the OpenVEX documents
// it can find by testing known locations of its identifiers and type.
func (agent *Agent) ProbePurl(purlString string) ([]*vex.VEX, error) {
	p, err := agent.impl.ParsePurl(purlString)
	if err != nil {
		return nil, fmt.Errorf("parsing purl: %w", err)
	}

	pkgProbe, err := agent.impl.GetPackageProbe(agent.Options, p)
	if err != nil {
		return nil, fmt.Errorf("getting package probe for purl type %s: %w", p.Type, err)
	}

	docs, err := agent.impl.FindDocumentsFromPurl(agent.Options, pkgProbe, p)
	if err != nil {
		return nil, fmt.Errorf("fetching documents: %w", err)
	}

	return docs, nil
}

// TODO(puerco): ProbeSBOM
// TODO(puerco): ProbeHash
