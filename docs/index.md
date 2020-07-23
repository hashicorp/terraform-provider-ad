---
layout: ""
page_title: "Provider: AD (Active Directory)"
description: |-
  The AD (Active Directory) provider provides resources to interact with an AD domain controller .
---

# AD (Active Directory) Provider

The AD (Active Directory) provider provides resources to interact with an AD domain controller.

## Example Usage

```terraform
resource "ad_ou" "o" { 
    name = "gplinktestOU"
    path = "dc=yourdomain,dc=com"
    description = "OU for gplink tests"
    protected = false
}
    
resource "ad_gpo" "g" {
    name        = "gplinktestGPO"
    domain      = "yourdomain.com"
    description = "gpo for gplink tests"
    status      = "AllSettingsEnabled"
}

resource "ad_gplink" "og" { 
    gpo_guid = ad_gpo.g.id
    target_dn = ad_ou.o.dn
    enforced = true
    enabled = true
    order = 0
}
```

## Schema

### Required

- **winrm_hostname** (String, Required) The hostname of the server we will use to run powershell scripts over WinRM.
- **winrm_password** (String, Required) The password used to authenticate to the the server's WinRM service.
- **winrm_username** (String, Required) The username used to authenticate to the the server's WinRM service.

### Optional

- **winrm_insecure** (Boolean, Optional) Trust unknown certificates. (default: false)
- **winrm_port** (Number, Optional) The port WinRM is listening for connections. (default: 5985)
- **winrm_proto** (String, Optional) The WinRM protocol we will use. (default: http)
