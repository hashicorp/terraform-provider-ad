---
page_title: "ad_ou Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_ou manages OU objects in an AD tree.
---

# Resource `ad_ou`

`ad_ou` manages OU objects in an AD tree.

## Example Usage

```terraform
resource "ad_ou" "o" { 
    name = "gplinktestOU"
    path = "dc=yourdomain,dc=com"
    description = "OU for gplink tests"
    protected = false
}
```

## Schema

### Required

- **name** (String, Required) Name of the OU.

### Optional

- **description** (String, Optional) Description of the OU.
- **id** (String, Optional) The ID of this resource.
- **path** (String, Optional) DN of the object that contains the OU.
- **protected** (Boolean, Optional) Protect this OU from being deleted accidentaly.

### Read-only

- **dn** (String, Read-only) The OU's DN.
- **guid** (String, Read-only) The OU's GUID.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_ou 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
