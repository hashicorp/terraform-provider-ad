// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package adschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GPOSecuritySchemaKeys is a list of all keys defined in the resource's schema
// except from gpo_container
var GPOSecuritySchemaKeys []string

// GpoSecuritySchema returns the GPO Security Settings resource schema
func GpoSecuritySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gpo_container": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The GUID of the container the security settings belong to.",
		},
		"password_policies": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: passwordPoliciesSchema()},
			Description: "Settings related to password policies. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b40db09-d95d-40a6-8467-32aedec8140c)",
		},
		"account_lockout": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: accountLockoutSchema()},
			Description: "Settings related to account lockout. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/2cd39c97-97cd-4859-a7b4-1229dad5f53d)",
		},
		"kerberos_policy": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: kerberosPolicySchema()},
			Description: "Settings related to kerberos policies. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0fce5b92-bcc1-4b96-9c2b-56397c3f144f)",
		},
		"system_log": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: systemLogSchema()},
			Description: "System log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf)",
		},
		"audit_log": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: auditLogSchema()},
			Description: "Audit log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf)",
		},
		"application_log": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: applicationLogSchema()},
			Description: "Application log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf)",
		},
		"event_audit": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: eventAuditSchema()},
			Description: "Event audit related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/01f8e057-f6a8-4d6e-8a00-99bcd241b403). Valid values for all items below are: 0 (None), 1 (Success audits only), 2 (Failure audits only), 3 (Success and failure audits), 4 (None)",
		},
		"restricted_groups": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Resource{Schema: restrictedGroupsSchema()},
			Description: "Settings related to Groups Membership. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/b73d8bae-ed22-48aa-acba-7065ab52d709)",
		},
		"registry_values": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Resource{Schema: registryValuesSchema()},
			Description: "Settings related to Registry Values. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/3a14ca47-a22f-43c5-b35e-6be791003ca7)",
		},
		"system_services": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Resource{Schema: systemServicesSchema()},
			Description: "Settings related to System Services. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/32deea3e-3fa4-414b-ba25-4121ad8c055c)",
		},
		"registry_keys": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Resource{Schema: registryKeysSchema()},
			Description: "Settings related to Registry Keys. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/13712a60-de1e-4642-bd9c-ab054dd86278)",
		},
		"filesystem": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Resource{Schema: filesystemSchema()},
			Description: "Settings related to File System permissions. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/abeebe06-49aa-44d4-ae5b-d6aff458e8e7)",
		},
	}
}

func passwordPoliciesSchema() map[string]*schema.Schema {
	k := map[string]string{
		"maximum_password_age":    "Number of days before password expires (-1-999). If set to -1, it means the password never expires.",
		"minimum_password_age":    "Number of days a password must be used before changing it (0-999).",
		"minimum_password_length": "Minimum number of characters used in a password (0-2^16). If set to 0, it means no password is required.",
		"password_complexity":     "Password must meet complexity requirements (0-2^16). If set to 0, then requirements do not apply, any other value means requirements are applied",
		"clear_text_password":     "Store password with reversible encryption (0-2^16). The password will not be stored with reversible encryption if the value is set to 0. Reversible encryption will be used in any other case.",
		"password_history_size":   "The number of unique new passwords that are required before an old password can be reused in association with a user account (0-2^16).  A value of 0 indicates that the password history is disabled.",
	}
	return generateSettingsSchema(k)
}

func accountLockoutSchema() map[string]*schema.Schema {
	k := map[string]string{
		"force_logoff_when_hour_expire": "Disconnect SMB sessions when logon hours expire.",
		"lockout_duration":              "Number of minutes a locked out account must remain locked out.",
		"lockout_bad_count":             "Number of failed logon attempts until a account is locked.",
		"reset_lockout_count":           "Number of minutes a account will remain locked after a failed logon attempt.",
	}
	return generateSettingsSchema(k)
}

func kerberosPolicySchema() map[string]*schema.Schema {
	k := map[string]string{
		"max_service_age":        "Maximum amount of minutes a ticket must be valid to access a service or resource. Minimum should be 10 and maximum should be equal to `max_ticket_age`.",
		"max_ticket_age":         "Maximum amount of hours a ticket-granting ticket is valid (0-99999).",
		"max_renew_age":          "Number of days during which a ticket-granting ticket can be renewed (0-99999).",
		"max_clock_skew":         "Maximum time difference, in minutes, between the client clock and the server clock. (0-99999).",
		"ticket_validate_client": "Control if the session ticket is validated for every request. A non-zero value disables the policy.",
	}
	return generateSettingsSchema(k)
}

func eventAuditSchema() map[string]*schema.Schema {
	k := map[string]string{
		"audit_system_events":    "Audit system events.",
		"audit_logon_events":     "Audit logon events.",
		"audit_privilege_use":    "Audit user attempts of exercising user rights.",
		"audit_policy_change":    "Audit attempts to change a policy.",
		"audit_account_manage":   "Audit account management events.",
		"audit_process_tracking": "Audit process related events.",
		"audit_ds_access":        "Audit access attempts to AD objects.",
		"audit_object_access":    "Audit access attempts to non-AD objects.",
		"audit_account_logon":    "Audit credential validation.",
	}
	return generateSettingsSchema(k)
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

// System, Audit, and Application log policies share the same key names.
func eventLogSchema() map[string]*schema.Schema {
	k := map[string]string{
		"maximum_log_size":           "Maximum size of log in KiloBytes. (64-4194240)",
		"audit_log_retention_period": "Control log retention. Values: 0: overwrite events as needed, 1: overwrite events as specified specified by `retention_days`, 2: never overwrite events.",
		"retention_days":             "Number of days before new events overwrite old events. (1-365)",
		"restrict_guest_access":      "Restrict access to logs for guest users. A non-zero value restricts access to guest users.",
	}
	return generateSettingsSchema(k)
}

// Since most of the settings are of the same type, we will use the function below
// to generate them, instead of repeating the same thing over and over
func generateSettingsSchema(keys map[string]string) map[string]*schema.Schema {
	sch := map[string]*schema.Schema{}
	for key, description := range keys {
		sch[key] = &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: description,
		}
	}
	return sch
}

func restrictedGroupsSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"group_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the group we are managing.",
		},
		"group_members": {
			Required:    true,
			Type:        schema.TypeString,
			Description: "Comma separated list of group names or SIDs that are members of the group.",
		},
		"group_memberof": {
			Required:    true,
			Type:        schema.TypeString,
			Description: "Comma separated list of group names or SIDs that this group belongs to.",
		},
	}
	return sch
}

func registryValuesSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"key_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fully qualified name of the key (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rrp/97587de7-3524-4291-8527-39517110c0eb)",
		},
		"value_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Data type of the key's value. 1: String, 2: Expand String, 3: Binary, 4: DWORD, 5: MULTI_SZ.",
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The value of the key, matching the type set in `value_type`.",
		},
	}
	return sch
}

func registryKeysSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"key_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fully qualified name of the key (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rrp/97587de7-3524-4291-8527-3951711      0c0eb)",
		},
		"propagation_mode": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Control permission propagation. 0: Propagate permissions to all subkeys, 1: Replace existing permissions on all subkeys, 2: Do not allow permissions to be replaced on the key.",
		},
		"acl": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)",
		},
	}
	return sch
}

func systemServicesSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the service.",
		},
		"startup_mode": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Startup mode of the service. Possible values are 2: Automatic, 3: Manual, 4: Disabled.",
		},
		"acl": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)",
		},
	}
	return sch
}

func filesystemSchema() map[string]*schema.Schema {
	sch := map[string]*schema.Schema{
		"path": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Path of the file or directory.",
		},
		"propagation_mode": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Control permission propagation. 0: Propagate permissions to all subfolders and files, 1: Replace existing permissions on all subfolders and files, 2: Do not allow permissions to be replaced.",
		},
		"acl": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)",
		},
	}
	return sch
}

func init() {
	for k := range GpoSecuritySchema() {
		GPOSecuritySchemaKeys = append(GPOSecuritySchemaKeys, k)
	}
}
