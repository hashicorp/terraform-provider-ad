// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !v5

package packet

func init() {
	V5Disabled = true
}
