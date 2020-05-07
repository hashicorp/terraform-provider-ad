package gposec

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// EventAudit represents the event audit policies section of the Security Settings GPO extension
type EventAudit struct {
	AuditAccountManage   string `ini:",omitempty" mapstructure:"audit_account_manage"`
	AuditDSAccess        string `ini:",omitempty" mapstructure:"audit_ds_access"`
	AuditAccountLogon    string `ini:",omitempty" mapstructure:"audit_account_logon"`
	AuditLogonEvents     string `ini:",omitempty" mapstructure:"audit_logon_events"`
	AuditObjectAccess    string `ini:",omitempty" mapstructure:"audit_object_access"`
	AuditPolicyChange    string `ini:",omitempty" mapstructure:"audit_policy_change"`
	AuditPrivilegeUse    string `ini:",omitempty" mapstructure:"audit_privilege_use"`
	AuditProcessTracking string `ini:",omitempty" mapstructure:"audit_process_tracking"`
	AuditSystemEvents    string `ini:",omitempty" mapstructure:"audit_system_events"`
}

// SetResourceData populates resource data based on the EventAudit field values
func (p *EventAudit) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p, d)
}

// WriteEventAudit populates an EventAudit struct from resource data
func WriteEventAudit(data interface{}, cfg *SecuritySettings) error {
	eap := &EventAudit{}
	err := mapstructure.Decode(data.(map[string]interface{}), eap)
	if err != nil {
		return err
	}
	cfg.EventAudit = eap
	return nil
}
