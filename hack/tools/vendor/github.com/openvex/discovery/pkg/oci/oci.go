// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package oci

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/openvex/go-vex/pkg/vex"
	purl "github.com/package-url/packageurl-go"
)

// IdentifiersBundle is a struct that collects different software identifiers
// and hashes in a structured way
type IdentifiersBundle struct {
	Identifiers map[vex.IdentifierType][]string
	Hashes      map[vex.Algorithm][]vex.Hash
}

// ToStringSlice returns all the identifiers and hashes contained in the bundle
// in a flat string slice
func (bundle *IdentifiersBundle) ToStringSlice() []string {
	ret := []string{}
	if bundle.Identifiers != nil {
		for _, sl := range bundle.Identifiers {
			ret = append(ret, sl...)
		}
	}

	if bundle.Hashes != nil {
		for _, sl := range bundle.Hashes {
			for _, h := range sl {
				ret = append(ret, string(h))
			}
		}
	}

	// Sort the slice to make the return value deterministic
	sort.Strings(ret)

	return ret
}

// GenerateReferenceIdentifiers reads an image reference string and
// generates a list of identifiers that can be used to match an entry
// in VEX a  document.
//
// This function returns the hashes and package urls to match the
// container image specified by the reference string. If the image
// is an index and os and arch are specified, the bundle will include
// purls and hashes for both the arch image and the index fronting it.
//
// For each image, the returned bundle will include a SHA256 hash with
// the image digest and two purls, with and without qualifiers. The
// variant with qualifiers will contain all the data known from the
// registry to match VEX documents with more specific purls.
//
// This function performs calls to the registry to retrieve data such
// as the image digests when needed.
func GenerateReferenceIdentifiers(refString, os, arch string) (IdentifiersBundle, error) {
	var dString, tag string
	bundle := IdentifiersBundle{
		Identifiers: map[vex.IdentifierType][]string{vex.PURL: {}},
		Hashes:      map[vex.Algorithm][]vex.Hash{vex.SHA256: {}},
	}

	ref, err := name.ParseReference(refString)
	if err != nil {
		return bundle, fmt.Errorf("parsing image reference: %w", err)
	}

	identifier := ref.Identifier()
	if strings.HasPrefix(identifier, "sha") {
		dString = identifier
	} else {
		tag = identifier
	}

	// If we dont have the digest in the reference, fetch it
	if dString == "" {
		dString, err = crane.Digest(refString)
		if err != nil {
			return bundle, fmt.Errorf("getting image digest: %w", err)
		}
	}

	bundle.Hashes[vex.SHA256] = append(
		bundle.Hashes[vex.SHA256], vex.Hash(strings.TrimPrefix(dString, "sha256:")),
	)

	pts := strings.Split(ref.Context().RepositoryStr(), "/")
	imageName := pts[len(pts)-1]
	registryPath := ref.Context().RegistryStr() + "/" + strings.ReplaceAll(ref.Context().RepositoryStr(), imageName, "")

	// Generate the variants for the input reference
	identifiers := generateImagePurlVariants(registryPath, imageName, dString, tag, os, arch)
	bundle.Identifiers[vex.PURL] = append(bundle.Identifiers[vex.PURL], identifiers...)

	if os == "" || arch == "" {
		return bundle, nil
	}

	// Now compute the identifiers for the platform specific image
	platform, err := v1.ParsePlatform(os + "/" + arch)
	if err != nil {
		return bundle, fmt.Errorf("parsing platform: %w", err)
	}

	archDString, err := crane.Digest(refString, crane.WithPlatform(platform))
	if err != nil {
		// If there is no arch-specific variant ot the image has not been pushed
		// yet, we simply don't include it. Return what we know.
		if strings.Contains(err.Error(), "no child with platform") ||
			strings.Contains(err.Error(), "MANIFEST_UNKNOWN") {
			return bundle, nil
		}
		return bundle, fmt.Errorf("getting image digest: %w", err)
	}

	// If the single-arch image digest is different, we generate purls for
	// it as we want to match the index and the arch image:
	if archDString != dString && archDString != "" {
		bundle.Identifiers[vex.PURL] = append(
			bundle.Identifiers[vex.PURL], generateImagePurlVariants(registryPath, imageName, archDString, tag, os, arch)...,
		)
		bundle.Hashes[vex.SHA256] = append(
			bundle.Hashes[vex.SHA256], vex.Hash(strings.TrimPrefix(archDString, "sha256:")),
		)
	}

	return bundle, nil
}

// generatePurlVariants
func generateImagePurlVariants(registryString, imageName, digestString, tag, os, arch string) []string {
	purls := []string{}

	// Purl with full qualifiers
	qMap := map[string]string{}
	if registryString != "" {
		qMap["repository_url"] = registryString + imageName
	}

	purls = append(purls,
		// Simple purl, no qualifiers
		purl.NewPackageURL(
			purl.TypeOCI, "", imageName, digestString,
			purl.QualifiersFromMap(qMap), "",
		).String(),
	)

	if tag != "" {
		qMap["tag"] = tag
	}
	if os != "" {
		qMap["os"] = os
	}
	if arch != "" {
		qMap["arch"] = arch
	}

	purls = append(purls,
		// Specific version with full qualifiers
		purl.NewPackageURL(
			purl.TypeOCI, "", imageName, digestString,
			purl.QualifiersFromMap(qMap), "",
		).String(),
	)

	return purls
}

type purlRefConverterOptions struct {
	// DefaultRepository will be added to the purl converter when none is found
	// in the package url qualifiers
	DefaultRepository string

	// Override repository will always be added to the purl, overriding
	// any that was set in the purl
	OverrideRepository string
}

type RefConverterOptions func(*purlRefConverterOptions)

func WithDefaultRepository(reg string) RefConverterOptions {
	return func(opts *purlRefConverterOptions) {
		opts.DefaultRepository = reg
	}
}

func WithOverrideRepository(reg string) RefConverterOptions {
	return func(opts *purlRefConverterOptions) {
		opts.OverrideRepository = reg
	}
}

// PurlToReferenceString reads a Package URL of type OCI and returns an image
// reference string. If the purl does not parse or is not of type oci: an error
// will be returned. The function takes a few options:
//
//	WithDefaultRepository(string)
//	Adds a default repository that will be used if none is defined in the
//	purl qualifiers.
//
//	WithOverrideRepository(string)
//	Overrides repository used in the reference regardless one is set in the purl
//	or not.
func PurlToReferenceString(purlString string, fopts ...RefConverterOptions) (string, error) {
	opts := &purlRefConverterOptions{}
	for _, opt := range fopts {
		opt(opts)
	}

	p, err := purl.FromString(purlString)
	if err != nil {
		return "", fmt.Errorf("parsing purl string: %w", err)
	}

	if p.Type != purl.TypeOCI {
		return "", errors.New("package URL is not of type OCI")
	}

	if p.Name == "" {
		return "", errors.New("parsed package URL did not return a package name")
	}

	qualifiers := p.Qualifiers.Map()

	refString := p.Name
	if v, ok := qualifiers["repository_url"]; ok {
		refString = v
	} else if opts.DefaultRepository != "" {
		refString = fmt.Sprintf(
			"%s/%s", strings.TrimSuffix(opts.DefaultRepository, "/"), p.Name,
		)
	}

	// If a repo override is set, rewrite the reference
	if opts.OverrideRepository != "" {
		refString = fmt.Sprintf(
			"%s/%s", strings.TrimSuffix(opts.OverrideRepository, "/"), p.Name,
		)
	}

	if p.Version != "" {
		refString = fmt.Sprintf("%s@%s", refString, p.Version)
	}

	// We add a tag, bu only if no digest is defined
	if _, ok := qualifiers["tag"]; ok && p.Version == "" {
		refString += ":" + qualifiers["tag"]
	}
	return refString, nil
}
