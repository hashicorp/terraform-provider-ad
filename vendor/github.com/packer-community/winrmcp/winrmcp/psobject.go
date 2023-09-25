// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package winrmcp

type pslist struct {
	Objects []psobject `xml:"Object"`
}

type psobject struct {
	Properties []psproperty `xml:"Property"`
	Value      string       `xml:",innerxml"`
}

type psproperty struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",innerxml"`
}
