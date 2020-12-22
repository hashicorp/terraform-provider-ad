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

- **guid** (String, Required) The GUID of the user object.

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **display_name** (String, Read-only) The display name of the user object.
- **principal_name** (String, Read-only) The principal name of the user object.
- **sam_account_name** (String, Read-only) The SAM account name of the user object.
- **city** (String, Read-only) The city assigned to user object.
- **company** (String, Read-only) The company assigned to user object.
- **country** (String, Read-only) The country assigned to user object.
- **department** (String, Read-only) The department assigned to user object.
- **description** (String, Read-only) The description assigned to user object.
- **division** (String, Read-only) The division assigned to user object.
- **email_address** (String, Read-only) The email address assigned to user object.
- **employee_id** (String, Read-only) The employee id assigned to user object.
- **employee_number** (String, Read-only) The employee number assigned to user object.
- **fax** (String, Read-only) The fax assigned to user object.
- **given_name** (String, Read-only) The given name assigned to user object.
- **home_directory** (String, Read-only) The home directory assigned to user object.
- **home_drive** (String, Read-only) The home drive assigned to user object.
- **home_phone** (String, Read-only) The home phone assigned to user object.
- **home_page** (String, Read-only) The home page assigned to user object.
- **initials** (String, Read-only) Initials assigned to user object.
- **mobile_phone** (String, Read-only) The mobile phone assigned to user object.
- **office** (String, Read-only) The office assigned to user object.
- **office_phone** (String, Read-only) The office phone assigned to user object.
- **organization** (String, Read-only) The organization assigned to user object.
- **other_name** (String, Read-only) Extra name of the user object.
- **po_box** (String, Read-only) The post office assigned to user object.
- **postal_code** (String, Read-only) The postal code assigned to user object.
- **smart_card_logon_required** (Boolean, Read-only) Smart card required to logon or not.
- **state** (String, Read-only) The state of the user object.
- **street_address** (String, Read-only) The address of the user object. assigned to user object.
- **surname** (String, Read-only) The surname assigned to user object.
- **title** (String, Read-only) The title assigned to user object.
- **trusted_for_delegation** (Boolean, Read-only) Check if user is trusted for delegation.

