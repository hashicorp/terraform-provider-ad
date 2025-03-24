// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dom

type Namespace struct {
	Prefix string
	Uri    string
}

func (ns *Namespace) SetTo(node *Element) {
	node.SetNamespace(ns.Prefix, ns.Uri)
}
