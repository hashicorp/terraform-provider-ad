---
page_title: "ad_ou Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Organizational Unit Active Directory object.
---

# Data Source `ad_ou`

Get the details of an Organizational Unit Active Directory object.

## Example Usage

```terraform
data "ad_ou" "g" {
    dn = "OU=SomeOU,dc=yourdomain,dc=com"
}

output "ou_uuid" {
    value = data.ad_ou.g.id
}
```

## Schema

### Optional

- **dn** (String, Optional) Distinguished Name of the OU object.
- **id** (String, Optional) The ID of this resource.
- **name** (String, Optional) Name of the OU object. If this is used then the `path` attribute needs to be set as well.
- **path** (String, Optional) Path of the OU object. If this is used then the `Name` attribute needs to be set as well.

### Read-only

- **description** (String, Read-only) The OU's description.
- **protected** (String, Read-only) The OU's protected status.


