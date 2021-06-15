---
page_title: "ad_gmsa Resource - terraform-provider-ad"
subcategory: ""
description: |-
  ad_gmsa manages Gmsa objects in an Active Directory tree.
---

# Resource `ad_gmsa`

`ad_gmsa` manages Gmsa objects in an Active Directory tree.

## Example Usage

```terraform
# basic example
variable principal_name { default = "testGmsa" }
variable samaccountname { default = "testGmsa" }

resource "ad_gmsa" "g" {
  name                                            = var.principal_name
  dns_host_name                                   = var.principal_name
}

# all Gmsa attributes
variable principal_name2 { default = "testGmsa2" }
variable samaccountname2 { default = "testGmsa2$" }
variable container       { default = "CN=Gmsas,DC=contoso,DC=com" }

resource "ad_gmsa" "g2" {
  name                                            = var.principal_name2
  sam_account_name                                = var.samaccountname2
  dns_host_name                                   = var.principal_name2  
  container                                       = var.container
  display_name                                    = "testGmsa2"
  description	                                    = "Some desc 2"
  delegated                                       = false
  managed_password_interval_in_days               = 15
  kerberos_encryption_type                        = [ "aes128","aes256" ]
  expiration                                      = "2021-12-30T00:00:00+00:00"
  service_principal_names                         = [
    "HTTP/Machine3.corp.contoso.com"
  ]
  principals_allowed_to_delegate_to_account       = [
    "CN=group1,DC=groups,DC=contoso,DC=com",
    "CN=computer1,DC=computers,DC=contoso,DC=com"
  ]
  principals_allowed_to_retrieve_managed_password = [
    "CN=group1,DC=groups,DC=contoso,DC=com",
    "CN=computer1,DC=computers,DC=contoso,DC=com"
  ]
}
```

## Schema

### Required

- **name** (String, Required) The Name of an Active Directory Gmsa.
- **dns_host_name** (String, Required) The DNS host name of Gmsa.

### Optional

- **container** (String, Optional) A DN of the container object that will be holding the Gmsa.
- **delegated** (Boolean, Optional) If set to false, the Gmsa will not be delegated to a service. Default value: true.
- **description** (String, Optional) Specifies a description of the object. This parameter sets the value of the Description property for the Gmsa object.
- **display_name** (String, Optional) The Display Name of an Active Directory Gmsa.
- **enabled** (Boolean, Optional) If set to false, the Gmsa will be disabled.
- **expiration** (String, Optional) Expiration date of the gmsa using RFC33339 format (https://tools.ietf.org/html/rfc3339).
- **home_page** (String, Optional) Specifies the URL of the home page of the object. This parameter sets the homePage property of a Gmsa object.
- **kerberos_encryption_type** (Set of String, Optional) This value sets the encryption types supported flags of the Active Directory. Valid: values are: rc4, aes128, aes256.
- **managed_password_interval_in_days** (Int, Optional) This value sets the number of days after which the password is automatically changed. Default values 30.
- **principals_allowed_to_delegate_to_account** (Set of String, Optional) Specifies the accounts which can act on the behalf of users to services running as this Managed Service Account or Group Managed Service Account. Principals must be set in DistinguishedName format.
- **principals_allowed_to_retrieve_managed_password** (Set of String, Optional) Specifies the membership policy for systems which can use a group managed service account. Principals must be set in DistinguishedName format.
- **sam_account_name** (String, Required) The pre-win2k gmsa logon name. Don't forget to add a "$" sign as last character of the string.
- **service_principal_names** (Set of String, Optional) Specifies the service principal names for the account.
- **trusted_for_delegation** (Boolean, Optional) Indicates whether an account is trusted for Kerberos delegation. Default value: false.

### Read-only

- **sid** (String, Read-only) The SID of the Gmsa object.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ad_gmsa 9CB8219C-31FF-4A85-A7A3-9BCBB6A41D02
```
