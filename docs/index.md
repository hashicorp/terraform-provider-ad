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
 - WinRM enabled.

## Note about Kerberos Authentication

Starting with version 0.4.0, this provider supports Kerberos Authentication for WinRM connections.
The underlying library used for Kerberos authentication supports setting its configuration by parsing
a configuration file as specified in this [page](https://web.mit.edu/kerberos/krb5-1.12/doc/admin/conf_files/krb5_conf.html).
If a configuration file is not supplied then we will use the equivalent of the following config:

```
[libdefaults]
   default_realm = YOURDOMAIN.COM
   dns_lookup_realm = false
   dns_lookup_kdc = false

[realms]
	YOURDOMAIN.COM = {
        kdc 	= 	192.168.1.122
        admin_server = 192.168.1.122
        default_domain = YOURDOMAIN.COM
	}

[domain_realm]
	yourdomain.com = YOURDOMAIN.COM
```

where `YOURDOMAIN.COM` is the value of the `krb_realm` setting, and 192.168.1.122 is the value of `winrm_hostname`.


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

- **krb_conf** (String, Optional) Path to kerberos configuration file. (default: none, environment variable: AD_KRB_CONF)
- **krb_realm** (String, Optional) The name of the kerberos realm (domain) we will use for authentication. (default: "", environment variable: AD_KRB_REALM)
- **krb_spn** (String, Optional) Alternative Service Principal Name. (default: none, environment variable: AD_KRB_SPN)
- **winrm_insecure** (Boolean, Optional) Trust unknown certificates. (default: false, environment variable: AD_WINRM_INSECURE)
- **winrm_use_ntlm** (Boolean, Optional) Use NTLM security. (default: false, environment variable: AD_WINRM_USE_NTLM)
- **winrm_port** (Number, Optional) The port WinRM is listening for connections. (default: 5985, environment variable: AD_PORT)
- **winrm_proto** (String, Optional) The WinRM protocol we will use. (default: http, environment variable: AD_PROTO)
