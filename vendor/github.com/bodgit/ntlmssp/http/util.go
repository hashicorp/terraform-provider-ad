// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package http

import "bytes"

func concat(bs ...[]byte) []byte {
	return bytes.Join(bs, nil)
}
