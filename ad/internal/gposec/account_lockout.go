// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// AccountLockout represents the account lockout section of the Security Settings GPO extension
type AccountLockout struct {
	ForceLogoffWhenHourExpire string `ini:",omitempty" mapstructure:"force_logoff_when_hour_expire"`
	LockoutDuration           string `ini:",omitempty" mapstructure:"lockout_duration"`
	LockoutBadCount           string `ini:",omitempty" mapstructure:"lockout_bad_count"`
	ResetLockoutCount         string `ini:",omitempty" mapstructure:"reset_lockout_count"`
}

// SetResourceData populates resource data based on the AccountLockout field values
func (p *AccountLockout) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p, d)
}

// WriteAccountLockout populates an AccountLockout struct from resource data
func WriteAccountLockout(data interface{}, cfg *SecuritySettings) error {
	al := &AccountLockout{}
	err := mapstructure.Decode(data.(map[string]interface{}), al)
	if err != nil {
		return err
	}

	if cfg.SystemAccess == nil {
		cfg.SystemAccess = &SystemAccess{}
	}
	cfg.SystemAccess.AccountLockout = al
	return nil
}
