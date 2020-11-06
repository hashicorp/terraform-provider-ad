---
page_title: "ad_group_membership Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_group_membership manages the members of a given Active Directory group.
---

# Resource `ad_group_membership`

`ad_group_membership` manages the members of a given Active Directory group.

## Example Usage

```terraform
variable name { default = "TestOU" }
variable path { default = "dc=yourdomain,dc=com" }
variable description { default = "some description" }
variable protected { default = false }
variable container { default = "CN=Users,dc=yourdomain,dc=com" }

variable name { default = "test group" }
variable sam_account_name { default = "TESTGROUP" }
variable scope { default = "global" }
variable category { default = "security" }

resource "ad_group" "g" {
  name             = var.name
  sam_account_name = var.sam_account_name
  scope            = var.scope
  category         = var.category
  container        = var.container
}

resource ad_group "g2" {
    name             = "${var.name}-2"
    sam_account_name = "${var.sam_account_name}-2"
    container        = var.container
}


resource ad_user "u" {
    display_name     = "test user"
    principal_name   = "testUser"
    sam_account_name = "testUser"
    initial_password = "SuperSecure1234!!"
    container        = var.container
}

resource ad_group_membership "gm" {
    group_id = ad_group.g.id
    group_members  = [ ad_group.g2.id, ad_user.u.id ]
}
```

## Schema

### Required

- **group_id** (String, Required) The ID of the group. This can be a GUID, a SID, a Distinguished Name, or the SAM Account Name of the group.
- **group_members** (Set of String, Required) A list of member AD Principals. Each principal can be identified by its GUID, SID, Distinguished Name, or SAM Account Name. Only one is required

### Optional

- **id** (String, Optional) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# The ID for this resource is the group's UUID plus a random UUID joined 
# by an underscore `_`.
$ terraform import ad_group_membership 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02_E9079B50-95C5-4101-8400-E01CC83CF53B
```
