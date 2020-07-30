# Windows Active Directory (AD) provider for Terraform
<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" align="right" height="50" />
</a>

![Status: Experimental](https://img.shields.io/badge/status-experimental-EAAA32) [![Releases](https://img.shields.io/github/release/hashicorp/terraform-provider-ad.svg)](https://github.com/hashicorp/terraform-provider-ad/releases)
[![LICENSE](https://img.shields.io/github/license/hashicorp/terraform-provider-ad.svg)](https://github.com/hashicorp/terraform-provider-ad/blob/master/LICENSE)
![Unit tests](https://github.com/hashicorp/terraform-provider-ad/workflows/Unit%20tests/badge.svg)

This Windows AD provider for Terraform allows you to manage users, groups and group policies in your AD installation.

Please regard this project as experimental. It still requires extensive testing and polishing to mature into production-ready quality. Please [file issues](https://github.com/hashicorp/terraform-provider-ad/issues/new/choose) generously and detail your experience while using the provider. We welcome your feedback.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x+
* [Windows Server](https://www.microsoft.com/en-us/windows-server) 2012R2 or greater 
* [Go](https://golang.org/doc/install) version 1.14.x+

## Getting Started

If this is your first time here, you can get an overview of the provider by reading our [introductory blog post](). Otherwise, start by downloading a copy of our latest build from the [registry](https://registry.terraform.io/providers/hashicorp/ad/latest).

Once you have the plugin installed, review the [docs](docs/) folder to understand which configuration options are available. You can find examples and more in [our examples folder](examples/). Don't forget to run `terraform init` in your Terraform configuration directory to allow Terraform to detect the provider plugin.

## Contributing

We welcome your contribution. Please understand that the experimental nature of this repository means that contributing code may be a bit of a moving target. If you have an idea for an enhancement or bug fix, and want to take on the work yourself, please first [create an issue](https://github.com/hashicorp/terraform-provider-ad/issues/new/choose) so that we can discuss the implementation with you before you proceed with the work.

You can review our [contribution guide](_about/CONTRIBUTING.md) to begin. You can also check out our [frequently asked questions](_about/FAQ.md).

## Experimental Status

By using the software in this repository (the "Software"), you acknowledge that: (1) the Software is still in development, may change, and has not been released as a commercial product by HashiCorp and is not currently supported in any way by HashiCorp; (2) the Software is provided on an "as-is" basis, and may include bugs, errors, or other issues;  (3) the Software is NOT INTENDED FOR PRODUCTION USE, use of the Software may result in unexpected results, loss of data, or other unexpected results, and HashiCorp disclaims any and all liability resulting from use of the Software; and (4) HashiCorp reserves all rights to make all decisions about the features, functionality and commercial release (or non-release) of the Software, at any time and without any obligation or liability whatsoever.
