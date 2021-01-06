---
page_title: "ad_user Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Active Directory user object.
---

# Data Source `ad_user`

Get the details of an Active Directory user object.

## Example Usage

```terraform
data "ad_user" "u" {
    guid = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "username" {
    value = data.ad_user.u.sam_account_name
}
```

## Schema

### Required

- **user_id** (String, Required) The user's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **display_name** (String, Read-only) The display name of the user object.
- **principal_name** (String, Read-only) The principal name of the user object.
- **sam_account_name** (String, Read-only) The SAM account name of the user object.


