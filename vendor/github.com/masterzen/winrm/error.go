// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package winrm

import "fmt"

// errWinrm generic error struct
type errWinrm struct {
	message string
}

// ErrWinrm implements the Error type interface
func (e errWinrm) Error() string {
	return fmt.Sprintf("%s", e.message)
}
