---
page_title: "ad_gpo Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Active Directory Group Policy Object.
---

# Data Source `ad_gpo`

Get the details of an Active Directory Group Policy Object.

## Example Usage

```terraform
data "ad_gpo" "g" {
    name = "Some GPO"
}

output "gpo_uuid" {
    value = data.ad_gpo.g.guid
}
```

## Schema

### Optional

- **guid** (String, Optional) GUID of the GPO.
- **id** (String, Optional) The ID of this resource.
- **name** (String, Optional) Name of the GPO.

### Read-only

- **domain** (String, Read-only) Domain of the GPO.


