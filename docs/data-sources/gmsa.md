---
page_title: "ad_gmsa Data Source - terraform-provider-ad"
subcategory: ""
description: |-
  Get the details of an Active Directory gmsa object.
---

# Data Source `ad_gmsa`

Get the details of an Active Directory gmsa object.

## Example Usage

```terraform
data "ad_gmsa" "g" {
    gmsa_id = "testGmsa$"
}

output "sam" {
    value = data.ad_gmsa.g.sam_account_name
}

output "trusted_for_delegation" {
    value = data.ad_gmsa.g.trusted_for_delegation
}

data "ad_gmsa" "g2" {
    gmsa_id = "CN=testGmsa,OU=gmsas,DC=contoso,DC=com"
}

output "testgmsa_guid" {
    value = data.ad_gmsa.g2.id
}
```

## Schema

### Required

- **gmsa_id** (String, Required) The gmsa's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.

### Optional

- **id** (String, Optional) The ID of this resource.

### Read-only

- **delegated** (Boolean, Read-only) If set to false, the Gmsa will not be delegated to a service. Default value: true.
- **description** (String, Read-only) Description of the gmsa object.
- **display_name** (String, Read-only) The display name of the gmsa object.
- **dns_host_name** (String, Required) The DNS host name of Gmsa.
- **enabled** (Boolean, Read-only) Check if gmsa is enabled.
- **expiration** (String, Read-only) Expiration date of the gmsa using RFC33339 format (https://tools.ietf.org/html/rfc3339).
- **home_page** (String, Read-only) Home page of the gmsa object.
- **kerberos_encryption_type** (Set of String, Read-only) The list of encryption types supported flags of the Active Directory.
- **managed_password_interval_in_days** (Int, Read-only) The value the number of days after which the password is automatically changed.
- **name** (String, Required) The Name of an Active Directory Gmsa.
- **principals_allowed_to_delegate_to_account** (Set of String, Read-only) List of accounts which can act on the behalf of users to services running as this Managed Service Account or Group Managed Service Account
- **principals_allowed_to_retrieve_managed_password** (Set of String, Read-onlyional) List of principals allowed to retrieve managed password.
 **sam_account_name** (String, Read-only) The SAM account name of the gmsa object.
- **service_principal_names** (Set of String, Read-only) List of SPN's.
- **sid** (String, Read-only) The SID of the Gmsa object.
- **trusted_for_delegation** (Boolean, Read-only) Check if gmsa is trusted for delegation


