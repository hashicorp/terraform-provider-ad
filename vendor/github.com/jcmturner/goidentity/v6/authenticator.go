// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package goidentity

type Authenticator interface {
	Authenticate() (Identity, bool, error)
	Mechanism() string // gives the name of the type of authentication mechanism
}
