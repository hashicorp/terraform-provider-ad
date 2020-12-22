---
page_title: "ad_user Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_user manages User objects in an Active Directory tree.
---

# Resource `ad_user`

`ad_user` manages User objects in an Active Directory tree.

## Example Usage

### Example 1 : basic

```terraform
variable principal_name { default = "testuser" }
variable samaccountname { default = "testuser" }

resource "ad_user" "u" {
  principal_name   = var.principal_name
  sam_account_name = var.samaccountname
  display_name     = "Terraform Test User"
}
```
### Example 2 : All attributes

```terraform
variable principal_name { default = "testuser2" }
variable samaccountname { default = "testuser2" }
variable container      { default = "CN=Users,DC=contoso,DC=com" }

resource "ad_user" "u" {
  principal_name            = var.principal_name
  sam_account_name          = var.samaccountname
  display_name              = "Terraform Test User"
  container                 = var.container
  initial_password          = "Password"
  city                      = "City"
  company                   = "Company"
  country                   = "us"
  department                = "Department"
  description               = "Description"
  division                  = "Division"
  email_address             = "some@email.com"
  employee_id               = "id"
  employee_number           = "number"
  fax                       = "Fax"
  given_name                = "GivenName"
  home_directory            = "HomeDirectory"
  home_drive                = "HomeDrive"
  home_phone                = "HomePhone"
  home_page                 = "HomePage"
  initials                  = "Initia"
  mobile_phone              = "MobilePhone"
  office                    = "Office"
  office_phone              = "OfficePhone"
  organization              = "Organization"
  other_name                = "OtherName"
  po_box                    = "POBox"
  postal_code               = "PostalCode"
  state                     = "State"
  street_address            = "StreetAddress"
  surname                   = "Surname"
  title                     = "Title"
  smart_card_logon_required = false
  trusted_for_delegation    = true
}
```


## Schema

### Required

- **display_name** (String, Required) The Display Name of an Active Directory user.
- **principal_name** (String, Required) The Principal Name of an Active Directory user.
- **sam_account_name** (String, Required) The pre-win2k user logon name.

### Optional

- **cannot_change_password** (Boolean, Optional) If set to true, the user will not be allowed to change their password.
- **container** (String, Optional) A DN of the container object that will be holding the user.
- **enabled** (Boolean, Optional) If set to false, the user will be disabled.
- **id** (String, Optional) The ID of this resource.
- **initial_password** (String, Optional) The user's initial password. This will be set on creation but will *not* be enforced in subsequent plans.
- **password_never_expires** (Boolean, Optional) If set to true, the password for this user will not expire.
- **city** (String, Optional) The city assigned to user object.
- **company** (String, Optional) The company assigned to user object.
- **country** (String, Optional) The country assigned to user object.
- **department** (String, Optional) The department assigned to user object.
- **description** (String, Optional) The description assigned to user object.
- **division** (String, Optional) The division assigned to user object.
- **email_address** (String, Optional) The email address assigned to user object.
- **employee_id** (String, Optional) The employee id assigned to user object.
- **employee_number** (String, Optional) The employee number assigned to user object.
- **fax** (String, Optional) The fax assigned to user object.
- **given_name** (String, Optional) The given name assigned to user object.
- **home_directory** (String, Optional) The home directory assigned to user object.
- **home_drive** (String, Optional) The home drive assigned to user object.
- **home_phone** (String, Optional) The home phone assigned to user object.
- **home_page** (String, Optional) The home page assigned to user object.
- **initials** (String, Optional) Initials assigned to user object.
- **mobile_phone** (String, Optional) The mobile phone assigned to user object.
- **office** (String, Optional) The office assigned to user object.
- **office_phone** (String, Optional) The office phone assigned to user object.
- **organization** (String, Optional) The organization assigned to user object.
- **other_name** (String, Optional) Extra name of the user object.
- **po_box** (String, Optional) The post office assigned to user object.
- **postal_code** (String, Optional) The postal code assigned to user object.
- **smart_card_logon_required** (Boolean, Optional) Smart card required to logon or not.
- **state** (String, Optional) The state of the user object.
- **street_address** (String, Optional) The address of the user object. assigned to user object.
- **surname** (String, Optional) The surname assigned to user object.
- **title** (String, Optional) The title assigned to user object.
- **trusted_for_delegation** (Boolean, Optional) Check if user is trusted for delegation.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_user 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
