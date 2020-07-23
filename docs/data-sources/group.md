---
page_title: "ad_group Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Active Directory Group object.
---

# Data Source `ad_group`

Get the details of an Active Directory Group object.

## Example Usage

```terraform
data "ad_group" "g" {
    guid = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "groupname" {
    value = data.ad_group.g.name
}
```

## Schema

### Required

- **guid** (String, Required) The GUID of the Group object.

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **category** (String, Read-only) The Group's category.
- **display_name** (String, Read-only) The display name of the Group object.
- **name** (String, Read-only) The name of the Group object.
- **sam_account_name** (String, Read-only) The SAM account name of the Group object.
- **scope** (String, Read-only) The Group's scope.


