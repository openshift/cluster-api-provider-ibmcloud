/*
   Copyright The containerd Authors.

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

// Package docker is deprecated, and has moved to github.com/distribution/reference.
//
// Deprecated: use github.com/distribution/reference instead.
package docker

import (
	"github.com/distribution/reference"
	"github.com/opencontainers/go-digest"
)

const (
	// NameTotalLengthMax is the maximum total number of characters in a repository name.
	//
	// Deprecated: use [reference.RepositoryNameTotalLengthMax].
	NameTotalLengthMax = reference.RepositoryNameTotalLengthMax
)

var (
	// ErrReferenceInvalidFormat represents an error while trying to parse a string as a reference.
	//
	// Deprecated: use [reference.ErrReferenceInvalidFormat].
	ErrReferenceInvalidFormat = reference.ErrReferenceInvalidFormat

	// ErrTagInvalidFormat represents an error while trying to parse a string as a tag.
	//
	// Deprecated: use [reference.ErrTagInvalidFormat].
	ErrTagInvalidFormat = reference.ErrTagInvalidFormat

	// ErrDigestInvalidFormat represents an error while trying to parse a string as a tag.
	//
	// Deprecated: use [reference.ErrDigestInvalidFormat].
	ErrDigestInvalidFormat = reference.ErrDigestInvalidFormat

	// ErrNameContainsUppercase is returned for invalid repository names that contain uppercase characters.
	//
	// Deprecated: use [reference.ErrNameContainsUppercase].
	ErrNameContainsUppercase = reference.ErrNameContainsUppercase

	// ErrNameEmpty is returned for empty, invalid repository names.
	//
	// Deprecated: use [reference.ErrNameEmpty].
	ErrNameEmpty = reference.ErrNameEmpty

	// ErrNameTooLong is returned when a repository name is longer than NameTotalLengthMax.
	//
	// Deprecated: use [reference.ErrNameTooLong].
	ErrNameTooLong = reference.ErrNameTooLong

	// ErrNameNotCanonical is returned when a name is not canonical.
	//
	// Deprecated: use [reference.ErrNameNotCanonical].
	ErrNameNotCanonical = reference.ErrNameNotCanonical
)

// Reference is an opaque object reference identifier that may include
// modifiers such as a hostname, name, tag, and digest.
//
// Deprecated: use [reference.Reference].
type Reference = reference.Reference

// Field provides a wrapper type for resolving correct reference types when
// working with encoding.
//
// Deprecated: use [reference.Field].
type Field = reference.Field

// AsField wraps a reference in a Field for encoding.
//
// Deprecated: use [reference.AsField].
func AsField(ref reference.Reference) reference.Field {
	return reference.AsField(ref)
}

// Named is an object with a full name
//
// Deprecated: use [reference.Named].
type Named = reference.Named

// Tagged is an object which has a tag
//
// Deprecated: use [reference.Tagged].
type Tagged = reference.Tagged

// NamedTagged is an object including a name and tag.
//
// Deprecated: use [reference.NamedTagged].
type NamedTagged reference.NamedTagged

// Digested is an object which has a digest
// in which it can be referenced by
//
// Deprecated: use [reference.Digested].
type Digested reference.Digested

// Canonical reference is an object with a fully unique
// name including a name with domain and digest
//
// Deprecated: use [reference.Canonical].
type Canonical reference.Canonical

// Domain returns the domain part of the [Named] reference.
//
// Deprecated: use [reference.Domain].
func Domain(named reference.Named) string {
	return reference.Domain(named)
}

// Path returns the name without the domain part of the [Named] reference.
//
// Deprecated: use [reference.Path].
func Path(named reference.Named) (name string) {
	return reference.Path(named)
}

// Parse parses s and returns a syntactically valid Reference.
// If an error was encountered it is returned, along with a nil Reference.
//
// Deprecated: use [reference.Parse].
func Parse(s string) (reference.Reference, error) {
	return reference.Parse(s)
}

// ParseNamed parses s and returns a syntactically valid reference implementing
// the Named interface. The reference must have a name and be in the canonical
// form, otherwise an error is returned.
// If an error was encountered it is returned, along with a nil Reference.
//
// Deprecated: use [reference.ParseNamed].
func ParseNamed(s string) (reference.Named, error) {
	return reference.ParseNamed(s)
}

// WithName returns a named object representing the given string. If the input
// is invalid ErrReferenceInvalidFormat will be returned.
//
// Deprecated: use [reference.WithName].
func WithName(name string) (reference.Named, error) {
	return reference.WithName(name)
}

// WithTag combines the name from "name" and the tag from "tag" to form a
// reference incorporating both the name and the tag.
//
// Deprecated: use [reference.WithTag].
func WithTag(name reference.Named, tag string) (reference.NamedTagged, error) {
	return reference.WithTag(name, tag)
}

// WithDigest combines the name from "name" and the digest from "digest" to form
// a reference incorporating both the name and the digest.
//
// Deprecated: use [reference.WithDigest].
func WithDigest(name reference.Named, digest digest.Digest) (reference.Canonical, error) {
	return reference.WithDigest(name, digest)
}

// TrimNamed removes any tag or digest from the named reference.
//
// Deprecated: use [reference.TrimNamed].
func TrimNamed(ref reference.Named) reference.Named {
	return reference.TrimNamed(ref)
}
