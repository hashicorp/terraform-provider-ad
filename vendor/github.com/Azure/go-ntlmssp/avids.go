// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ntlmssp

type avID uint16

const (
	avIDMsvAvEOL avID = iota
	avIDMsvAvNbComputerName
	avIDMsvAvNbDomainName
	avIDMsvAvDNSComputerName
	avIDMsvAvDNSDomainName
	avIDMsvAvDNSTreeName
	avIDMsvAvFlags
	avIDMsvAvTimestamp
	avIDMsvAvSingleHost
	avIDMsvAvTargetName
	avIDMsvChannelBindings
)
