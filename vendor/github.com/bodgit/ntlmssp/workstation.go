// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ntlmssp

import (
	"os"
	"strings"
)

// DefaultWorkstation returns the current workstation name.
func DefaultWorkstation() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return strings.ToUpper(hostname), nil
}
