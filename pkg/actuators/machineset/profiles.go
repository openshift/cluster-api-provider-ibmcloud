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

package machineset

// Profile is spec of IBM Cloud instance
type Profile struct {
	Profile       string
	VCPU          int64
	MemoryGb      int64
	BandwidthGbps int64
}

// Profiles is a map of IBM Cloud instance resources
var Profiles = map[string]*Profile{
	"bx2-2x8": {
		Profile:       "bx2-2x8",
		VCPU:          2,
		MemoryGb:      8,
		BandwidthGbps: 4,
	},
	"bx2d-2x8": {
		Profile:       "bx2d-2x8",
		VCPU:          2,
		MemoryGb:      8,
		BandwidthGbps: 4,
	},
	"bx2-4x16": {
		Profile:       "bx2-4x16",
		VCPU:          4,
		MemoryGb:      16,
		BandwidthGbps: 8,
	},
	"bx2d-4x16": {
		Profile:       "bx2d-4x16",
		VCPU:          4,
		MemoryGb:      16,
		BandwidthGbps: 4,
	},
	"bx2d-8x32": {
		Profile:       "bx2d-8x32",
		VCPU:          8,
		MemoryGb:      32,
		BandwidthGbps: 16,
	},
	"bx2-8x32": {
		Profile:       "bx2-8x32",
		VCPU:          8,
		MemoryGb:      32,
		BandwidthGbps: 16,
	},
	"bx2-16x64": {
		Profile:       "bx2-16x64",
		VCPU:          16,
		MemoryGb:      64,
		BandwidthGbps: 32,
	},
	"bx2d-16x64": {
		Profile:       "bx2d-16x64",
		VCPU:          16,
		MemoryGb:      64,
		BandwidthGbps: 32,
	},
	"bx2-32x128": {
		Profile:       "bx2-32x128",
		VCPU:          32,
		MemoryGb:      128,
		BandwidthGbps: 64,
	},
	"bx2d-32x128": {
		Profile:       "bx2d-32x128",
		VCPU:          32,
		MemoryGb:      128,
		BandwidthGbps: 64,
	},
	"bx2-48x192": {
		Profile:       "bx2-48x192",
		VCPU:          48,
		MemoryGb:      192,
		BandwidthGbps: 80,
	},
	"bx2d-48x192": {
		Profile:       "bx2d-48x192",
		VCPU:          48,
		MemoryGb:      192,
		BandwidthGbps: 80,
	},
	"bx2d-64x256": {
		Profile:       "bx2d-64x256",
		VCPU:          64,
		MemoryGb:      256,
		BandwidthGbps: 80,
	},
	"bx2-64x256": {
		Profile:       "bx2-64x256",
		VCPU:          64,
		MemoryGb:      256,
		BandwidthGbps: 80,
	},
	"bx2d-96x384": {
		Profile:       "bx2d-96x384",
		VCPU:          96,
		MemoryGb:      384,
		BandwidthGbps: 80,
	},
	"bx2-96x384": {
		Profile:       "bx2-96x384",
		VCPU:          96,
		MemoryGb:      384,
		BandwidthGbps: 80,
	},
	"bx2-128x512": {
		Profile:       "bx2-128x512",
		VCPU:          128,
		MemoryGb:      512,
		BandwidthGbps: 80,
	},
	"bx2d-128x512": {
		Profile:       "bx2d-128x512",
		VCPU:          128,
		MemoryGb:      512,
		BandwidthGbps: 80,
	},
	"cx2-2x4": {
		Profile:       "cx2-2x4",
		VCPU:          2,
		MemoryGb:      4,
		BandwidthGbps: 4,
	},
	"cx2d-2x4": {
		Profile:       "cx2d-2x4",
		VCPU:          2,
		MemoryGb:      4,
		BandwidthGbps: 4,
	},
	"cx2d-4x8": {
		Profile:       "cx2d-4x8",
		VCPU:          4,
		MemoryGb:      8,
		BandwidthGbps: 8,
	},
	"cx2-4x8": {
		Profile:       "cx2-4x8",
		VCPU:          4,
		MemoryGb:      8,
		BandwidthGbps: 8,
	},
	"cx2d-8x16": {
		Profile:       "cx2d-8x16",
		VCPU:          8,
		MemoryGb:      16,
		BandwidthGbps: 16,
	},
	"cx2-8x16": {
		Profile:       "cx2-8x16",
		VCPU:          8,
		MemoryGb:      16,
		BandwidthGbps: 16,
	},
	"cx2-16x32": {
		Profile:       "cx2-16x32",
		VCPU:          16,
		MemoryGb:      32,
		BandwidthGbps: 32,
	},
	"cx2d-16x32": {
		Profile:       "cx2d-16x32",
		VCPU:          16,
		MemoryGb:      32,
		BandwidthGbps: 32,
	},
	"cx2d-32x64": {
		Profile:       "cx2d-32x64",
		VCPU:          32,
		MemoryGb:      64,
		BandwidthGbps: 64,
	},
	"cx2-32x64": {
		Profile:       "cx2-32x64",
		VCPU:          32,
		MemoryGb:      64,
		BandwidthGbps: 64,
	},
	"cx2-48x96": {
		Profile:       "cx2-48x96",
		VCPU:          48,
		MemoryGb:      96,
		BandwidthGbps: 80,
	},
	"cx2d-48x96": {
		Profile:       "cx2d-48x96",
		VCPU:          48,
		MemoryGb:      96,
		BandwidthGbps: 80,
	},
	"cx2-64x128": {
		Profile:       "cx2-64x128",
		VCPU:          64,
		MemoryGb:      128,
		BandwidthGbps: 80,
	},
	"cx2d-64x128": {
		Profile:       "cx2d-64x128",
		VCPU:          64,
		MemoryGb:      128,
		BandwidthGbps: 80,
	},
	"cx2-96x192": {
		Profile:       "cx2-96x192",
		VCPU:          96,
		MemoryGb:      192,
		BandwidthGbps: 80,
	},
	"cx2d-96x192": {
		Profile:       "cx2d-96x192",
		VCPU:          96,
		MemoryGb:      192,
		BandwidthGbps: 80,
	},
	"cx2d-128x256": {
		Profile:       "cx2d-128x256",
		VCPU:          128,
		MemoryGb:      256,
		BandwidthGbps: 80,
	},
	"cx2-128x256": {
		Profile:       "cx2-128x256",
		VCPU:          128,
		MemoryGb:      256,
		BandwidthGbps: 80,
	},
	"mx2d-2x16": {
		Profile:       "mx2d-2x16",
		VCPU:          2,
		MemoryGb:      16,
		BandwidthGbps: 4,
	},
	"mx2-2x16": {
		Profile:       "mx2-2x16",
		VCPU:          2,
		MemoryGb:      16,
		BandwidthGbps: 4,
	},
	"mx2d-4x32": {
		Profile:       "mx2d-4x32",
		VCPU:          4,
		MemoryGb:      32,
		BandwidthGbps: 8,
	},
	"mx2-4x32": {
		Profile:       "mx2-4x32",
		VCPU:          4,
		MemoryGb:      32,
		BandwidthGbps: 8,
	},
	"mx2d-8x64": {
		Profile:       "mx2d-8x64",
		VCPU:          8,
		MemoryGb:      64,
		BandwidthGbps: 16,
	},
	"mx2-8x64": {
		Profile:       "mx2-8x64",
		VCPU:          8,
		MemoryGb:      64,
		BandwidthGbps: 16,
	},
	"mx2d-16x128": {
		Profile:       "mx2d-16x128",
		VCPU:          16,
		MemoryGb:      128,
		BandwidthGbps: 32,
	},
	"mx2-16x128": {
		Profile:       "mx2-16x128",
		VCPU:          16,
		MemoryGb:      128,
		BandwidthGbps: 32,
	},
	"mx2-32x256": {
		Profile:       "mx2-32x256",
		VCPU:          32,
		MemoryGb:      256,
		BandwidthGbps: 64,
	},
	"mx2d-32x256": {
		Profile:       "mx2d-32x256",
		VCPU:          32,
		MemoryGb:      256,
		BandwidthGbps: 64,
	},
	"mx2-48x384": {
		Profile:       "mx2-48x384",
		VCPU:          48,
		MemoryGb:      384,
		BandwidthGbps: 80,
	},
	"mx2d-48x384": {
		Profile:       "mx2d-48x384",
		VCPU:          48,
		MemoryGb:      384,
		BandwidthGbps: 80,
	},
	"mx2d-64x512": {
		Profile:       "mx2d-64x512",
		VCPU:          64,
		MemoryGb:      512,
		BandwidthGbps: 80,
	},
	"mx2-64x512": {
		Profile:       "mx2d-64x512",
		VCPU:          64,
		MemoryGb:      512,
		BandwidthGbps: 80,
	},
	"mx2-96x768": {
		Profile:       "mx2-96x768",
		VCPU:          96,
		MemoryGb:      768,
		BandwidthGbps: 80,
	},
	"mx2d-96x768": {
		Profile:       "mx2d-96x768",
		VCPU:          96,
		MemoryGb:      768,
		BandwidthGbps: 80,
	},
	"mx2-128x1024": {
		Profile:       "mx2-128x1024",
		VCPU:          128,
		MemoryGb:      1024,
		BandwidthGbps: 80,
	},
	"mx2d-128x1024": {
		Profile:       "mx2d-128x1024",
		VCPU:          128,
		MemoryGb:      1024,
		BandwidthGbps: 80,
	},
}
