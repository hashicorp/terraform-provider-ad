// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// AuditLog represents the Audit Log section of the Security Settings GPO extension
type AuditLog struct {
	EventLogPolicy `ini:"Audit Log,omitempty,squash"`
}

// SetResourceData populates resource data based on the AuditLog field values
func (p *AuditLog) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p.EventLogPolicy, d)
}

// WriteAuditLog populates an AuditLog struct from resource data
func WriteAuditLog(data interface{}, cfg *SecuritySettings) error {
	elp, err := NewEventLogPolicy(data)
	if err != nil {
		return err
	}
	cfg.AuditLog = &AuditLog{EventLogPolicy: elp}
	return nil
}
