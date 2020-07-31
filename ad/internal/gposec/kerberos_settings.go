package gposec

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// KerberosPolicy represents the kerberos settings section of the Security Settings GPO extension
type KerberosPolicy struct {
	MaxServiceAge        string `ini:",omitempty" mapstructure:"max_service_age"`
	MaxTicketAge         string `ini:",omitempty" mapstructure:"max_ticket_age"`
	MaxRenewAge          string `ini:",omitempty" mapstructure:"max_renew_age"`
	MaxClockSkew         string `ini:",omitempty" mapstructure:"max_clock_skew"`
	TicketValidateClient string `ini:",omitempty" mapstructure:"ticket_validate_client"`
}

// SetResourceData populates resource data based on the KerberosSettings field values
func (p *KerberosPolicy) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p, d)

}

// WriteKerberosPolicy populates a KerberosSettings struct from resource data
func WriteKerberosPolicy(data interface{}, cfg *SecuritySettings) error {
	ks := &KerberosPolicy{}
	err := mapstructure.Decode(data.(map[string]interface{}), ks)
	if err != nil {
		return err
	}
	cfg.KerberosPolicy = ks
	return nil
}
