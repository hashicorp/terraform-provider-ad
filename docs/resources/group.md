---
page_title: "ad_group Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_group manages Group objects in an Active Directory tree.
---

# Resource `ad_group`

`ad_group` manages Group objects in an Active Directory tree.

## Example Usage

```terraform
variable name { default = "TestOU" }
variable path { default = "dc=yourdomain,dc=com" }
variable description { default = "some description" }
variable protected { default = false }

variable name { default = "test group" }
variable sam_account_name { default = "TESTGROUP" }
variable scope { default = "global" }
variable category { default = "security" }

resource "ad_ou" "o" {
  name        = var.name
  path        = var.path
  description = var.description
  protected   = var.protected
}


resource "ad_group" "g" {
  name             = var.name
  sam_account_name = var.sam_account_name
  scope            = var.scope
  category         = var.category
  container        = ad_ou.o.dn
}
```

## Schema

### Required

- **container** (String, Required) A DN of a container object holding the group.
- **name** (String, Required) The name of the group.
- **sam_account_name** (String, Required) The pre-win2k name of the group.

### Optional

- **category** (String, Optional) The group's category. Can be one of `system` or `security` (case sensitive).
- **id** (String, Optional) The ID of this resource.
- **scope** (String, Optional) The group's scope. Can be one of `global`, `local`, or `universal` (case sensitive).

### Read-only

- **sid** (String, Read-only) The SID of the group object.
- **description** (String, Read-only) Description of the Group object.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_group 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
