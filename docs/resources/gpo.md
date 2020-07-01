---
page_title: "ad_gpo Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_gpo manages Group Policy Objects (GPOs).
---

# Resource `ad_gpo`

`ad_gpo` manages Group Policy Objects (GPOs).

## Example Usage

```terraform
variable "domain" { default = "yourdomain.com" }
variable "name" { default = "tfGPO" }

resource "ad_gpo" "gpo" {
  name   = var.name
  domain = var.domain
}
```

## Schema

### Required

- **name** (String, Required) Name of the Group Policy Object.

### Optional

- **description** (String, Optional) Description of the GPO.
- **domain** (String, Optional) Domain of the GPO.
- **id** (String, Optional) The ID of this resource.
- **status** (String, Optional) Status of the GPO. Can be one of `AllSettingsEnabled`, `UserSettingsDisabled`, `ComputerSettingsDisabled`, or `AllSettingsDisabled` (case sensitive).

### Read-only

- **dn** (String, Read-only)
- **numeric_status** (Number, Read-only)

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_gpo 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
