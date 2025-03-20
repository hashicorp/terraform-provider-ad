// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ntlmssp

import "bytes"

func concat(bs ...[]byte) []byte {
	return bytes.Join(bs, nil)
}

func zeroBytes(length int) []byte {
	return make([]byte, length, length)
}

func zeroPad(s []byte, length int) []byte {
	d := zeroBytes(length)
	copy(d, s)

	return d
}
