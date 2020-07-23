---
page_title: "ad_user Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_user manages User objects in an Active Directory tree.
---

# Resource `ad_user`

`ad_user` manages User objects in an Active Directory tree.

## Example Usage

```terraform
variable principal_name { default = "testuser" }
variable samaccountname { default = "testuser" }

resource "ad_user" "u" {
  principal_name   = var.principal_name
  sam_account_name = var.samaccountname
  display_name     = "Terraform Test User"
}
```

## Schema

### Required

- **display_name** (String, Required) The Display Name of an Active Directory user.
- **principal_name** (String, Required) The Principal Name of an Active Directory user.
- **sam_account_name** (String, Required) The pre-win2k user logon name.

### Optional

- **cannot_change_password** (Boolean, Optional) If set to true the user will not be allowed to change their password.
- **container** (String, Optional) A DN of the container object that will be holding the user.
- **enabled** (Boolean, Optional) If set to false the user will be disabled.
- **id** (String, Optional) The ID of this resource.
- **initial_password** (String, Optional) The user's initial password. This will be set on creation but will *not* be enforced in subsequent plans.
- **password_never_expires** (Boolean, Optional) If set to true the password for this user will not expire.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_user 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
