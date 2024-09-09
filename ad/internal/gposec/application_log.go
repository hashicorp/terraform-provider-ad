// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// ApplicationLog represents the Application Log section of the Security Settings GPO extension
type ApplicationLog struct {
	EventLogPolicy `ini:"Application Log,omitempty,squash"`
}

// SetResourceData populates resource data based on the ApplicationLog field values
func (p *ApplicationLog) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p.EventLogPolicy, d)
}

// WriteApplicationLog populates a WriteApplicationLog struct from resource data
func WriteApplicationLog(data interface{}, cfg *SecuritySettings) error {
	elp, err := NewEventLogPolicy(data)
	if err != nil {
		return err
	}
	cfg.ApplicationLog = &ApplicationLog{EventLogPolicy: elp}
	return nil
}
