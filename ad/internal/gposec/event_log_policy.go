package gposec

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// EventLogPolicy is a structure that is used by the next three identical structures
type EventLogPolicy struct {
	MaximumLogSize          string `ini:",omitempty" mapstructure:"maximum_log_size,omitempty"`
	AuditLogRetentionPeriod string `ini:",omitempty" mapstructure:"audit_log_retention_period,omitempty"`
	RetentionDays           string `ini:",omitempty" mapstructure:"retention_days,omitempty"`
	RestrictGuestAccess     string `ini:",omitempty" mapstructure:"restrict_guest_access,omitempty"`
}

// SetResourceData populates resource data based on the EventLogPolicy field values
func (p *EventLogPolicy) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p, d)
}

// NewEventLogPolicy returns an EventLogPolicy structure populated from resource data
func NewEventLogPolicy(data interface{}) (EventLogPolicy, error) {
	elp := EventLogPolicy{}
	err := mapstructure.Decode(data.(map[string]interface{}), &elp)
	if err != nil {
		return EventLogPolicy{}, err
	}
	return elp, nil
}
