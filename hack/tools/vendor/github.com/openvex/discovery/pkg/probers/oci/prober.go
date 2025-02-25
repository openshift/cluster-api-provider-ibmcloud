// SPDX-FileCopyrightText: Copyright 2023 The OpenVEX Authors
// SPDX-License-Identifier: Apache-2.0

package oci

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	purl "github.com/package-url/packageurl-go"

	"github.com/sigstore/cosign/v2/pkg/cosign"
	"github.com/sigstore/cosign/v2/pkg/oci"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
	"github.com/sigstore/cosign/v2/pkg/types"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/openvex/discovery/pkg/discovery/options"
	doci "github.com/openvex/discovery/pkg/oci"
	"github.com/openvex/go-vex/pkg/attestation"
	"github.com/openvex/go-vex/pkg/vex"
)

type Prober struct {
	Options options.Options
	impl    ociImplementation
}

func New() *Prober {
	p := &Prober{
		impl:    &defaultImplementation{},
		Options: options.Default,
	}
	p.Options.ProberOptions[purl.TypeOCI] = localOptions{}
	return p
}

//counterfeiter:generate . ociImplementation
type ociImplementation interface {
	VerifyOptions(*options.Options) error
	PurlToReference(options.Options, purl.PackageURL) (name.Reference, error)
	ResolveImageReference(options.Options, name.Reference) (oci.SignedEntity, error)
	DownloadDocuments(options.Options, oci.SignedEntity) ([]*vex.VEX, error)
}

type defaultImplementation struct{}

type localOptions struct {
	Platform           string
	TagPrefix          string // Attestation "image" prefix
	Repository         string
	RepositoryOverride string // COSIGN_REPOSITORY or other repo that overrides the purl repo
}

type platformList []struct {
	hash     v1.Hash
	platform *v1.Platform
}

func (pl *platformList) String() string {
	r := []string{}
	for _, p := range *pl {
		r = append(r, p.platform.String())
	}
	return strings.Join(r, ", ")
}

// FindDocumentsFromPurl implements the logic to search for OpenVEX documents
// attached to a container image
func (prober *Prober) FindDocumentsFromPurl(opts options.Options, p purl.PackageURL) ([]*vex.VEX, error) {
	if err := prober.impl.VerifyOptions(&prober.Options); err != nil {
		return nil, fmt.Errorf("verifying options: %w", err)
	}

	ref, err := prober.impl.PurlToReference(prober.Options, p)
	if err != nil {
		return nil, fmt.Errorf("translating purl to image reference: %w", err)
	}

	if ref == nil {
		return nil, fmt.Errorf("could not resolve image reference from %s", p)
	}

	image, err := prober.impl.ResolveImageReference(prober.Options, ref)
	if err != nil {
		return nil, fmt.Errorf("resolving image reference: %w", err)
	}

	docs, err := prober.impl.DownloadDocuments(prober.Options, image)
	if err != nil {
		return nil, fmt.Errorf("downloading documents from registry for %s: %w", p.String(), err)
	}

	return docs, nil
}

// PurlToReference reads a purl and generates an image reference. It uses GGCR's
// name package to parse it and returns the reference.
func (di *defaultImplementation) PurlToReference(opts options.Options, p purl.PackageURL) (name.Reference, error) {
	refString, err := doci.PurlToReferenceString(p.String())
	if err != nil {
		return nil, err
	}

	ref, err := name.ParseReference(refString)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %s: %w", refString, err)
	}
	return ref, nil
}

// getIndexPlatforms returns the platforms of the single arch images fronted by
// an image index.
func getIndexPlatforms(idx oci.SignedImageIndex) (platformList, error) {
	im, err := idx.IndexManifest()
	if err != nil {
		return nil, fmt.Errorf("fetching index manifest: %w", err)
	}

	platforms := platformList{}
	for i := range im.Manifests {
		if im.Manifests[i].Platform == nil {
			continue
		}
		platforms = append(platforms, struct {
			hash     v1.Hash
			platform *v1.Platform
		}{im.Manifests[i].Digest, im.Manifests[i].Platform})
	}
	return platforms, nil
}

// ResolveImageReference takes an image ref returns the signed entity it is
// pointing to. This process involves checking if the image is an index, a
// single or multi arch image, if we have an arch in the options, etc, etc.
func (di *defaultImplementation) ResolveImageReference(opts options.Options, ref name.Reference) (oci.SignedEntity, error) {
	if ref == nil {
		return nil, fmt.Errorf("got nil value when trying to resolve OCI image reference")
	}

	ociremoteOpts := []ociremote.Option{}

	// TODO(puerco): Support relevant registry options
	// o := options.RegistryOptions{}
	// ociremoteOpts := []ociremote.Option{ociremote.WithRemoteOptions(o.GetRegistryClientOpts(ctx)...)}

	// Support tag prefix
	if opts.ProberOptions[purl.TypeOCI].(localOptions).TagPrefix != "" {
		ociremoteOpts = append(ociremoteOpts, ociremote.WithPrefix(opts.ProberOptions[purl.TypeOCI].(localOptions).TagPrefix))
	}

	// Registry override. From options or env
	var targetRepoOverride name.Repository
	var err error
	if opts.ProberOptions[purl.TypeOCI].(localOptions).RepositoryOverride != "" {
		targetRepoOverride, err = name.NewRepository(opts.ProberOptions[purl.TypeOCI].(localOptions).RepositoryOverride)
		if err != nil {
			return nil, fmt.Errorf("parsing override repository option")
		}
	} else {
		targetRepoOverride, err = ociremote.GetEnvTargetRepository()
		if err != nil {
			return nil, fmt.Errorf("fetching repository from environment")
		}
	}
	if (targetRepoOverride != name.Repository{}) {
		ociremoteOpts = append(ociremoteOpts, ociremote.WithTargetRepository(targetRepoOverride))
	}

	se, err := ociremote.SignedEntity(ref, ociremoteOpts...)
	if err != nil {
		return nil, err
	}

	idx, isIndex := se.(oci.SignedImageIndex)

	// We only allow --platform on multiarch indexes
	if opts.ProberOptions[purl.TypeOCI].(localOptions).Platform != "" && !isIndex {
		return nil, fmt.Errorf("specified reference is not a multiarch image")
	}

	// If a platform was specified, then we return the corresponding
	// single arch image if there is one
	if opts.ProberOptions[purl.TypeOCI].(localOptions).Platform != "" && isIndex {
		opts.Logger.DebugContext(
			opts.Context, "Reference is an index and arch %s defined", "imageRef", ref.String(),
		)
		targetPlatform, err := v1.ParsePlatform(opts.ProberOptions[purl.TypeOCI].(localOptions).Platform)
		if err != nil {
			return nil, fmt.Errorf("parsing platform: %w", err)
		}
		platforms, err := getIndexPlatforms(idx)
		if err != nil {
			return nil, fmt.Errorf("getting available platforms: %w", err)
		}

		platforms = matchPlatform(targetPlatform, platforms)
		if len(platforms) == 0 {
			return nil, fmt.Errorf("unable to find an attestation for %s", targetPlatform.String())
		}
		if len(platforms) > 1 {
			return nil, fmt.Errorf(
				"platform spec matches more than one image architecture: %s",
				platforms.String(),
			)
		}

		nse, err := idx.SignedImage(platforms[0].hash)
		if err != nil {
			return nil, fmt.Errorf("searching for %s image: %w", platforms[0].hash.String(), err)
		}
		if nse == nil {
			return nil, fmt.Errorf("unable to find image %s", platforms[0].hash.String())
		}
		se = nse
	}

	return se, nil
}

// matchPlatform filters a list of platforms returning only those matching
// a base. "Based" on ko's internal equivalent while it moves to GGCR.
// https://github.com/google/ko/blob/e6a7a37e26d82a8b2bb6df991c5a6cf6b2728794/pkg/build/gobuild.go#L1020
func matchPlatform(base *v1.Platform, list platformList) platformList {
	ret := platformList{}
	for _, p := range list {
		if base.OS != "" && base.OS != p.platform.OS {
			continue
		}
		if base.Architecture != "" && base.Architecture != p.platform.Architecture {
			continue
		}
		if base.Variant != "" && base.Variant != p.platform.Variant {
			continue
		}

		if base.OSVersion != "" && p.platform.OSVersion != base.OSVersion {
			if base.OS != "windows" {
				continue
			} else { //nolint: revive
				if pcount, bcount := strings.Count(base.OSVersion, "."), strings.Count(p.platform.OSVersion, "."); pcount == 2 && bcount == 3 {
					if base.OSVersion != p.platform.OSVersion[:strings.LastIndex(p.platform.OSVersion, ".")] {
						continue
					}
				} else {
					continue
				}
			}
		}
		ret = append(ret, p)
	}

	return ret
}

// DownloadDocuments retrieves attested or attached document from the registry
func (di *defaultImplementation) DownloadDocuments(opts options.Options, se oci.SignedEntity) ([]*vex.VEX, error) {
	docs := []*vex.VEX{}

	attestations, err := cosign.FetchAttestations(se, vex.Context)
	if err != nil {
		// If the image has no attestations attached, cosign returns  an
		// error. Trap it here and handle properly
		if err.Error() == "found no attestations" || strings.Contains(err.Error(), "no attestations with predicate type") {
			opts.Logger.DebugContext(opts.Context, "image has no attestations attached")
			return docs, nil
		}

		return nil, fmt.Errorf("fetching attestations: %w", err)
	}

	opts.Logger.DebugContext(
		opts.Context, fmt.Sprintf("image has %d OpenVEX attestations", len(attestations)),
	)
	for i, att := range attestations {
		if att.PayloadType != types.IntotoPayloadType {
			continue
		}

		pload, err := base64.StdEncoding.DecodeString(att.PayLoad)
		if err != nil {
			opts.Logger.WarnContext(
				opts.Context, fmt.Sprintf("error decoding openvex attestation %d from base64, ignoring", i),
			)
			continue
		}

		statement := attestation.Attestation{}

		// If the attestation could not be parsed, then we ignore it and
		// issue a warning
		if err := json.Unmarshal(pload, &statement); err != nil {
			opts.Logger.WarnContext(
				opts.Context, fmt.Sprintf("error parsing openvex attestation #%d, ignoring", i),
			)
			continue
		}

		docs = append(docs, &statement.Predicate)
	}

	return docs, nil
}

// VerifyOptions checks the options and returns an error if there is something wrong
func (di *defaultImplementation) VerifyOptions(opts *options.Options) error {
	if opts.ProberOptions == nil {
		opts.ProberOptions = map[string]interface{}{}
	}
	if _, ok := opts.ProberOptions[purl.TypeOCI]; !ok {
		opts.ProberOptions[purl.TypeOCI] = localOptions{}
	}
	return nil
}

// SetOptions sets the probe's options
func (prober *Prober) SetOptions(opts options.Options) {
	prober.Options = opts
}
