---
layout: ""
page_title: "Provider: AD (Active Directory)"
description: |-
  The AD (Active Directory) provider provides resources to interact with an AD domain controller .
---

# AD (Active Directory) Provider

The AD (Active Directory) provider provides resources to interact with an AD domain controller.

Requirements:
 - Windows Server 2012R2 or greater.
 - WinRM enabled. Currently only local accounts with `Basic` authentication are supported.


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

- **winrm_hostname** (String, Required) The hostname of the server we will use to run powershell scripts over WinRM. (Environment variable: AD_HOSTNAME)
- **winrm_password** (String, Required) The password used to authenticate to the server's WinRM service. (Environment variable: AD_PASSWORD)
- **winrm_username** (String, Required) The username used to authenticate to the server's WinRM service. (Environment variable: AD_USER)

### Optional

- **winrm_insecure** (Boolean, Optional) Trust unknown certificates. (default: false, environment variable: AD_WINRM_INSECURE)
- **winrm_use_ntlm** (Boolean, Optional) Use NTLM authentication. (default: false, environment variable: AD_WINRM_USE_NTLM)
- **winrm_port** (Number, Optional) The port WinRM is listening for connections. (default: 5985, environment variable: AD_PORT)
- **winrm_proto** (String, Optional) The WinRM protocol we will use. (default: http, environment variable: AD_PROTO)
