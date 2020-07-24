---
page_title: "ad_computer Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_computer manages computer objects in an AD tree.
---

# Resource `ad_computer`

`ad_computer` manages computer objects in an AD tree.

## Example Usage

```terraform
variable "name" { default = "test" }
variable "pre2kname" { default = "TEST" }

resource "ad_computer" "c" {
  name      = var.name
  pre2kname = var.pre2kname
}
```

## Schema

### Required

- **name** (String, Required) The name for the computer account.

### Optional

- **container** (String, Optional) The DN of the container used to hold the computer account.
- **id** (String, Optional) The ID of this resource.
- **pre2kname** (String, Optional) The pre-win2k name for the computer account.

### Read-only

- **dn** (String, Read-only)
- **guid** (String, Read-only)

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_computer 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
