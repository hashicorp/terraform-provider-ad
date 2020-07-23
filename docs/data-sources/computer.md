---
page_title: "ad_computer Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Active Directory Computer object.
---

# Data Source `ad_computer`

Get the details of an Active Directory Computer object.

## Example Usage

```terraform
data "ad_computer" "c" {
    dn = "cn=test123,cn=Computers,dc=yourdomain,dc=com"
}

output "computer_guid" {
  value = data.ad_computer.c.guid
}
```

## Schema

### Optional

- **dn** (String, Optional) The Distinguished Name of the computer object.
- **guid** (String, Optional) The GUID of the computer object.
- **id** (String, Optional) The ID of this resource.

### Read-only

- **name** (String, Read-only) The name of the computer object.


