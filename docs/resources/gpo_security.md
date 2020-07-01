---
page_title: "ad_gpo_security Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_gpo_security manages the security settings portion of a Group Policy Object (GPO).
---

# Resource `ad_gpo_security`

`ad_gpo_security` manages the security settings portion of a Group Policy Object (GPO).

## Example Usage

```terraform
variable domain { default = "yourdomain.com" }
variable gpo_name { default = "TestGPO" }

resource "ad_gpo" "gpo" {
  name   = var.gpo_name
  domain = var.domain
}

resource "ad_gpo_security" "gpo_sec" {
  gpo_container = ad_gpo.gpo.id
  password_policies {
    minimum_password_length = 3
  }

  system_services {
    service_name = "TapiSrv"
    startup_mode = "2"
    acl          = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;LA)"
  }

  system_services {
    service_name = "CertSvc"
    startup_mode = "2"
    acl          = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;SY)(A;;CCLCSWLOCRRC;;;IU)S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)"
  }

}
```

## Schema

### Required

- **gpo_container** (String, Required) The GUID of the container the security settings belong to.

### Optional

- **account_lockout** (Block List, Max: 1) Settings related to account lockout. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/2cd39c97-97cd-4859-a7b4-1229dad5f53d) (see [below for nested schema](#nestedblock--account_lockout))
- **application_log** (Block List, Max: 1) Application log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf) (see [below for nested schema](#nestedblock--application_log))
- **audit_log** (Block List, Max: 1) Audit log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf) (see [below for nested schema](#nestedblock--audit_log))
- **event_audit** (Block List, Max: 1) Event audit related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/01f8e057-f6a8-4d6e-8a00-99bcd241b403). Valid values for all items below are: 0 (None), 1 (Success audits only), 2 (Failure audits only), 3 (Success and failure audits), 4 (None) (see [below for nested schema](#nestedblock--event_audit))
- **filesystem** (Block Set) Settings related to File System permissions. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/abeebe06-49aa-44d4-ae5b-d6aff458e8e7) (see [below for nested schema](#nestedblock--filesystem))
- **id** (String, Optional) The ID of this resource.
- **kerberos_policy** (Block List, Max: 1) Settings related to kerberos policies. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0fce5b92-bcc1-4b96-9c2b-56397c3f144f) (see [below for nested schema](#nestedblock--kerberos_policy))
- **password_policies** (Block List, Max: 1) Settings related to password policies. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b40db09-d95d-40a6-8467-32aedec8140c) (see [below for nested schema](#nestedblock--password_policies))
- **registry_keys** (Block Set) Settings related to Registry Keys. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/13712a60-de1e-4642-bd9c-ab054dd86278) (see [below for nested schema](#nestedblock--registry_keys))
- **registry_values** (Block Set) Settings related to Registry Values. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/3a14ca47-a22f-43c5-b35e-6be791003ca7) (see [below for nested schema](#nestedblock--registry_values))
- **restricted_groups** (Block Set) Settings related to Groups Membership. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/b73d8bae-ed22-48aa-acba-7065ab52d709) (see [below for nested schema](#nestedblock--restricted_groups))
- **system_log** (Block List, Max: 1) System log related settings. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b9673a7-ce0a-49b4-912b-591efdb37cdf) (see [below for nested schema](#nestedblock--system_log))
- **system_services** (Block Set) Settings related to System Services. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/32deea3e-3fa4-414b-ba25-4121ad8c055c) (see [below for nested schema](#nestedblock--system_services))

<a id="nestedblock--account_lockout"></a>
### Nested Schema for `account_lockout`

Optional:

- **force_logoff_when_hour_expire** (String, Optional) Disconnect SMB sessions when logon hours expire.
- **lockout_bad_count** (String, Optional) Number of failed logon attempts until a account is locked.
- **lockout_duration** (String, Optional) Number of minutes a locked out account must remain locked out.
- **reset_lockout_count** (String, Optional) Number of minutes a account will remain locked after a failed logon attempt.


<a id="nestedblock--application_log"></a>
### Nested Schema for `application_log`

Optional:

- **audit_log_retention_period** (String, Optional) Control log retention. Values: 0: overwrite events as needed, 1: overwrite events as specified specified by `retention_days`, 2: never overwrite events.
- **maximum_log_size** (String, Optional) Maximum size of log in KiloBytes. (64-4194240)
- **restrict_guest_access** (String, Optional) Restrict access to logs for guest users. A non-zero value restricts access to guest users.
- **retention_days** (String, Optional) Number of days before new events overwrite old events. (1-365)


<a id="nestedblock--audit_log"></a>
### Nested Schema for `audit_log`

Optional:

- **audit_log_retention_period** (String, Optional) Control log retention. Values: 0: overwrite events as needed, 1: overwrite events as specified specified by `retention_days`, 2: never overwrite events.
- **maximum_log_size** (String, Optional) Maximum size of log in KiloBytes. (64-4194240)
- **restrict_guest_access** (String, Optional) Restrict access to logs for guest users. A non-zero value restricts access to guest users.
- **retention_days** (String, Optional) Number of days before new events overwrite old events. (1-365)


<a id="nestedblock--event_audit"></a>
### Nested Schema for `event_audit`

Optional:

- **audit_account_logon** (String, Optional) Audit credential validation.
- **audit_account_manage** (String, Optional) Audit account management events.
- **audit_ds_access** (String, Optional) Audit access attempts to AD objects.
- **audit_logon_events** (String, Optional) Audit logon events.
- **audit_object_access** (String, Optional) Audit access attempts to non-AD objects.
- **audit_policy_change** (String, Optional) Audit attempts to change a policy.
- **audit_privilege_use** (String, Optional) Audit user attempts of exercising user rights.
- **audit_process_tracking** (String, Optional) Audit process related events.
- **audit_system_events** (String, Optional) Audit system events.


<a id="nestedblock--filesystem"></a>
### Nested Schema for `filesystem`

Required:

- **acl** (String, Required) Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)
- **path** (String, Required) Path of the file or directory.
- **propagation_mode** (String, Required) Control permission propagation. 0: Propagate permissions to all subfolders and files, 1: Replace existing permissions on all subfolders and files, 2: Do not allow permissions to be replaced


<a id="nestedblock--kerberos_policy"></a>
### Nested Schema for `kerberos_policy`

Optional:

- **max_clock_skew** (String, Optional) Maximum time difference,in minutes, between the client clock and the server clock. (0-99999).
- **max_renew_age** (String, Optional) Number of days during which a ticket-granting ticket can be renewed (0-99999).
- **max_service_age** (String, Optional) Maximum amount of minutes a ticket must be valid to access a service or resource. Minimum should be 10 and maximum should be equal to `max_ticket_age`.
- **max_ticket_age** (String, Optional) Maximum amount of hours a ticket-granting ticket is valid (0-99999).
- **ticket_validate_client** (String, Optional) Control if the session ticket is validated for every request. A non-zero value disables the policy.


<a id="nestedblock--password_policies"></a>
### Nested Schema for `password_policies`

Optional:

- **clear_text_password** (String, Optional) Store password with reversible encryption (0-2^16). The password will not be stored with reversible encryption if the value is set to 0. Reversible encryption will be used in any other case.
- **maximum_password_age** (String, Optional) Number of days before password expires (-1-999). If set to -1 it means the password never expires.
- **minimum_password_age** (String, Optional) Number of days a password must be used before changing it (0-999).
- **minimum_password_length** (String, Optional) Minimum number of characters used in a password (0-2^16). If set to 0 it means no password is required.
- **password_complexity** (String, Optional) Password must meet complexity requirements (0-2^16). If set to 0 then requirements do not apply, any other value means requirements are applied
- **password_history_size** (String, Optional) The number of unique new passwords that are required before an old password can be reused in association with a user account (0-2^16).  A value of 0 indicates that the password history is disabled.


<a id="nestedblock--registry_keys"></a>
### Nested Schema for `registry_keys`

Required:

- **acl** (String, Required) Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)
- **key_name** (String, Required) Fully qualified name of the key (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rrp/97587de7-3524-4291-8527-3951711      0c0eb)
- **propagation_mode** (String, Required) Control permission propagation. 0: Propagate permissions to all subkeys, 1: Replace existing permissions on all subkeys, 2: Do not allow permissions to be replaced on the key


<a id="nestedblock--registry_values"></a>
### Nested Schema for `registry_values`

Required:

- **key_name** (String, Required) Fully qualified name of the key (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rrp/97587de7-3524-4291-8527-39517110c0eb)
- **value** (String, Required) The value of the key, matching the type set in `value_type`.
- **value_type** (String, Required) Data type of the key's value. 1: String, 2: Expand String, 3: Binary, 4: DWORD, 5: MULTI_SZ.


<a id="nestedblock--restricted_groups"></a>
### Nested Schema for `restricted_groups`

Required:

- **group_memberof** (String, Required) Comma separated list of group names or SIDs that this group belongs to.
- **group_members** (String, Required) Comma separated list of group names or SIDs that are members of the group.
- **group_name** (String, Required) Name of the group we are managing.


<a id="nestedblock--system_log"></a>
### Nested Schema for `system_log`

Optional:

- **audit_log_retention_period** (String, Optional) Control log retention. Values: 0: overwrite events as needed, 1: overwrite events as specified specified by `retention_days`, 2: never overwrite events.
- **maximum_log_size** (String, Optional) Maximum size of log in KiloBytes. (64-4194240)
- **restrict_guest_access** (String, Optional) Restrict access to logs for guest users. A non-zero value restricts access to guest users.
- **retention_days** (String, Optional) Number of days before new events overwrite old events. (1-365)


<a id="nestedblock--system_services"></a>
### Nested Schema for `system_services`

Required:

- **acl** (String, Required) Security descriptor to apply. (https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/f4296d69-1c0f-491f-9587-a960b292d070)
- **service_name** (String, Required) Name of the service.
- **startup_mode** (String, Required) Startup mode of the service. Possible values are 2: Automatic, 3: Manual, 4: Disabled.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_gpo_security 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02_securitysettings
```
