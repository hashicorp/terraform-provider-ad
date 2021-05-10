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
`Basic` remains the default authentication method, although this may change in the future. The provider will use
Kerberos as its authentication when `krb_realm` is set.

## Note about Local execution (Windows only)

It is possible to execute commands locally if the OS on which terraform is running is Windows.
In such case, your need to put the following settings in the provider configuration :

- Set winrm_username to null
- Set winrm_password to null
- Set winrm_hostname to null

Note: it will set to local only `if all 3 parameters are set to null`

### Example
```terraform
provider "ad" {
  winrm_hostname = ""
  winrm_username = ""
  winrm_password = ""
}
```

 ## Example Usage

```terraform
variable "hostname" { default = "ad.yourdomain.com" }
variable "username" { default = "user" }
variable "password" { default = "password" }

// remote using Basic authentication
provider "ad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
}

// remote using NTLM authentication
provider "ad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  winrm_use_ntlm = true
}

// remote using NTLM authentication and HTTPS
provider "ad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  winrm_use_ntlm = true
  winrm_port     = 5986
  winrm_proto    = "https"
  winrm_insecure = true
}

// remote using Kerberos authentication
provider "ad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  krb_realm      = "YOURDOMAIN.COM"
}

// remote using Kerberos authentication with krb5.conf file
provider "ad" {
  winrm_hostname = var.hostname
  winrm_username = var.username
  winrm_password = var.password
  krb_conf       = "/etc/krb5.conf"
}

// local (windows only)
provider "ad" {
  winrm_hostname = ""
  winrm_username = ""
  winrm_password = ""
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
- **winrm_port** (Number, Optional) The port WinRM is listening for connections. (default: 5985, environment variable: AD_PORT)
- **winrm_proto** (String, Optional) The WinRM protocol we will use. (default: http, environment variable: AD_PROTO)
- **winrm_use_ntlm** (Boolean, Optional) Use NTLM authentication. (default: false, environment variable: AD_WINRM_USE_NTLM)
