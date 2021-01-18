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
    user_id = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "username" {
    value = data.ad_user.u.sam_account_name
}

output "country" {
    value = data.ad_user.u.country
}

output "trusted_for_delegation" {
    value = data.ad_user.u.trusted_for_delegation
}

data "ad_user" "u2" {
    user_id = "CN=Test User,OU=Users,DC=contoso,DC=com"
}

output "testuser_guid" {
    value = data.ad_user.u2.id
}
```

## Schema

### Required

- **user_id** (String, Required) The user's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **city** (String, Read-only) City assigned to user object.
- **company** (String, Read-only) Company assigned to user object.
- **country** (String, Read-only) Country assigned to user object.
- **department** (String, Read-only) Department assigned to user object.
- **description** (String, Read-only) Description of the user object.
- **display_name** (String, Read-only) The display name of the user object.
- **division** (String, Read-only) Division assigned to user object.
- **email_address** (String, Read-only) Email address assigned to user object.
- **employee_id** (String, Read-only) Employee ID assigned to user object.
- **employee_number** (String, Read-only) Employee Number assigned to user object.
- **fax** (String, Read-only) Fax number assigned to user object.
- **given_name** (String, Read-only) Given name of the user object.
- **home_directory** (String, Read-only) Home directory of the user object.
- **home_drive** (String, Read-only) Home drive of the user object.
- **home_page** (String, Read-only) Home page of the user object.
- **home_phone** (String, Read-only) Home phone of the user object.
- **initials** (String, Read-only) Initials of the user object.
- **mobile_phone** (String, Read-only) Mobile phone of the user object.
- **office** (String, Read-only) Office assigned to user object.
- **office_phone** (String, Read-only) Office phone of the user object.
- **organization** (String, Read-only) Organization assigned to user object.
- **other_name** (String, Read-only) Extra name of the user object.
- **po_box** (String, Read-only) Post office assigned to user object.
- **postal_code** (String, Read-only) Postal code of the user object.
- **principal_name** (String, Read-only) The principal name of the user object.
- **sam_account_name** (String, Read-only) The SAM account name of the user object.
- **smart_card_logon_required** (Boolean, Read-only) Smart card required to logon or not
- **state** (String, Read-only) State of the user object.
- **street_address** (String, Read-only) Address of the user object.
- **surname** (String, Read-only) Surname of the user object.
- **title** (String, Read-only) Title of the user object
- **trusted_for_delegation** (Boolean, Read-only) Check if user is trusted for delegation


