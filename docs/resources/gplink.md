---
page_title: "ad_gplink Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_gplink manages links between GPOs and container objects such as OUs.
---

# Resource `ad_gplink`

`ad_gplink` manages links between GPOs and container objects such as OUs.

## Example Usage

```terraform
resource "ad_ou" "o" {
  name        = "gplinktestOU"
  path        = "dc=yourdomain,dc=com"
  description = "OU for gplink tests"
  protected   = false
}

resource "ad_gpo" "g" {
  name        = "gplinktestGPO"
  domain      = "yourdomain.com"
  description = "gpo for gplink tests"
  status      = "AllSettingsEnabled"
}

resource "ad_gplink" "og" {
  gpo_guid  = ad_gpo.g.id
  target_dn = ad_ou.o.dn
  enforced  = true
  enabled   = true
}
```

## Schema

### Required

- **gpo_guid** (String, Required) The GUID of the GPO that will be linked to the container object.
- **target_dn** (String, Required) The DN of the object the GPO will be linked to.

### Optional

- **enabled** (Boolean, Optional) Controls the state of the GP link between a GPO and a container object.
- **enforced** (Boolean, Optional) If set to true, the GPO will be enforced on the container object.
- **id** (String, Optional) The ID of this resource.
- **order** (Number, Optional) Sets the precedence between multiple GPOs linked to the same container object.

## Import

Import is supported using the following syntax:

```shell
# The ID for this resource is comprised of the GPO GUID and the container (OU) GUID separated by 
# an underscore
$ terraform import ad_gplink 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02_5AB72AD7-1AA0-4D97-923B-6FBCDD143CB2
```
