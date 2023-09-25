# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# basic example
variable principal_name { default = "testuser" }
variable samaccountname { default = "testuser" }

resource "ad_user" "u" {
  principal_name    = var.principal_name
  sam_account_name  = var.samaccountname
  display_name      = "Terraform Test User"
  custom_attributes = jsonencode({
    "carLicense": ["This is", "a multi-value", "attribute"],
    "comment": "and this is a single value attribute"
  })
}


# all user attributes
variable principal_name2 { default = "testuser2" }
variable samaccountname2 { default = "testuser2" }
variable container      { default = "CN=Users,DC=contoso,DC=com" }

resource "ad_user" "u2" {
  principal_name            = var.principal_name2
  sam_account_name          = var.samaccountname2
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

