package gposec

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SystemLog represents the system log section of the Security Settings GPO extension
type SystemLog struct {
	EventLogPolicy `ini:"System Log,omitempty,squash"`
}

// SetResourceData populates resource data based on the SystemLog field values
func (p *SystemLog) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p.EventLogPolicy, d)
}

// WriteSystemLog populates a SystemLog struct from resource data
func WriteSystemLog(data interface{}, cfg *SecuritySettings) error {
	elp, err := NewEventLogPolicy(data)
	if err != nil {
		return err
	}
	cfg.SystemLog = &SystemLog{EventLogPolicy: elp}
	return nil
}
