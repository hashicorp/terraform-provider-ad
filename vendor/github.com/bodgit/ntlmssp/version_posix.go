// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// +build !windows

package ntlmssp

// DefaultVersion returns a pointer to a NTLM Version struct for the OS which
// will be populated on Windows or nil otherwise.
func DefaultVersion() *Version {
	return nil
}
