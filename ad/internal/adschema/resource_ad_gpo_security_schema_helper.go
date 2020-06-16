package adschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// GPOSecuritySchemaKeys is a list of all keys defined in the resource's schema
// except from gpo_container
var GPOSecuritySchemaKeys []string

// GpoSecuritySchema returns the GPO Security Settings resource schema
func GpoSecuritySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gpo_container": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password_policies": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: passwordPoliciesSchema()},
		},
		"account_lockout": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: accountLockoutSchema()},
		},
		"kerberos_policy": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: kerberosPolicySchema()},
		},
		"system_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: systemLogSchema()},
		},
		"audit_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: auditLogSchema()},
		},
		"application_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: applicationLogSchema()},
		},
		"event_audit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem:     &schema.Resource{Schema: eventAuditSchema()},
		},
		"restricted_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: restrictedGroupsSchema()},
		},
		"registry_values": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: registryValuesSchema()},
		},
		"system_services": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: systemServicesSchema()},
		},
		"registry_keys": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: registryKeysSchema()},
		},
		"filesystem": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: filesystemSchema()},
		},
	}
}

func passwordPoliciesSchema() map[string]*schema.Schema {
	keys := []string{"maximum_password_age", "minimum_password_age",
		"minimum_password_length", "password_complexity", "clear_text_password",
		"password_history_size"}
	return generateSettingsSchema(keys)
}

func accountLockoutSchema() map[string]*schema.Schema {
	keys := []string{"force_logoff_when_hour_expire", "lockout_duration",
		"lockout_bad_count", "reset_lockout_count"}
	return generateSettingsSchema(keys)
}

func kerberosPolicySchema() map[string]*schema.Schema {
	keys := []string{"max_service_age", "max_ticket_age", "max_renew_age", "max_clock_skew",
		"ticket_validate_client"}
	return generateSettingsSchema(keys)
}

func eventAuditSchema() map[string]*schema.Schema {
	keys := []string{"audit_system_events", "audit_logon_events", "audit_privilege_use",
		"audit_policy_change", "audit_account_manage", "audit_process_tracking",
		"audit_ds_access", "audit_object_access", "audit_account_logon"}
	return generateSettingsSchema(keys)
}

func systemLogSchema() map[string]*schema.Schema {
	return eventLogSchema()
}

func auditLogSchema() map[string]*schema.Schema {
	return eventLogSchema()
}

func applicationLogSchema() map[string]*schema.Schema {
	return eventLogSchema()
}

func restrictedGroupsSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"group_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"group_members": {
			Required: true,
			Type:     schema.TypeString,
		},
		"group_memberof": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
	return sch
}

func registryValuesSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return sch
}

func registryKeysSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"propagation_mode": {
			Type:     schema.TypeString,
			Required: true,
		},
		"acl": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return sch
}

func systemServicesSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"service_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"startup_mode": {
			Type:     schema.TypeString,
			Required: true,
		},
		"acl": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return sch
}

func filesystemSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"path": {
			Type:     schema.TypeString,
			Required: true,
		},
		"propagation_mode": {
			Type:     schema.TypeString,
			Required: true,
		},
		"acl": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return sch
}

// System, Audit, and Application log policies share the same key names.
func eventLogSchema() map[string]*schema.Schema {
	keys := []string{"maximum_log_size", "audit_log_retention_period", "retention_days",
		"restrict_guest_access"}
	return generateSettingsSchema(keys)
}

// Since most of the settings are of the same type, we will use the function below
// to generate them, instead of repeating the same thing over and over
func generateSettingsSchema(keys []string) map[string]*schema.Schema {
	sch := map[string]*schema.Schema{}
	for _, key := range keys {
		sch[key] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		}
	}
	return sch
}

func init() {
	for k := range GpoSecuritySchema() {
		GPOSecuritySchemaKeys = append(GPOSecuritySchemaKeys, k)
	}
}
