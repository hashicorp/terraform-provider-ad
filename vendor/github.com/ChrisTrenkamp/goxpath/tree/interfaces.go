// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tree

import "fmt"

//Result is used for all data types.
type Result interface {
	fmt.Stringer
}

//IsBool is used for the XPath boolean function.  It turns the data type to a bool.
type IsBool interface {
	Bool() Bool
}

//IsNum is used for the XPath number function.  It turns the data type to a number.
type IsNum interface {
	Num() Num
}
