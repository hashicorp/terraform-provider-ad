package gposec

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// PasswordPolicies represents the password policies section of the Security Settings GPO extension
type PasswordPolicies struct {
	MaximumPasswordAge    string `ini:",omitempty" mapstructure:"maximum_password_age"`
	MinimumPasswordAge    string `ini:",omitempty" mapstructure:"minimum_password_age"`
	MinimumPasswordLength string `ini:",omitempty" mapstructure:"minimum_password_length"`
	PasswordComplexity    string `ini:",omitempty" mapstructure:"password_complexity"`
	ClearTextPassword     string `ini:",omitempty" mapstructure:"clear_text_password"`
	PasswordHistorySize   string `ini:",omitempty" mapstructure:"password_history_size"`
}

// SetResourceData populates resource data based on the PasswordPolicies field values
func (p *PasswordPolicies) SetResourceData(section string, d *schema.ResourceData) error {
	return genericSetResourceData(section, p, d)
}

// WritePasswordPolicies populates a PasswordPolicies struct from resource data
func WritePasswordPolicies(data interface{}, cfg *SecuritySettings) error {
	pp := &PasswordPolicies{}
	err := mapstructure.Decode(data.(map[string]interface{}), pp)
	if err != nil {
		return err
	}

	if cfg.SystemAccess == nil {
		cfg.SystemAccess = &SystemAccess{}
	}
	cfg.SystemAccess.PasswordPolicies = pp
	return nil
}
