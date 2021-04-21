## 0.4.2 (Unreleased)

BUGFIXES:
* **Resource:** `ad_user`: Fix bug when removing user attributes. ([#77](https://github.com/hashicorp/terraform-provider-ad/pull/77))
* **Resource:** `ad_group`: Use correct command when updating AD groups. ([#83](https://github.com/hashicorp/terraform-provider-ad/pull/83))
* **Resource:** `ad_group`: Fix category name. ([#69](https://github.com/hashicorp/terraform-provider-ad/pull/69))

FEATURES:
* **provider:** Execute commands as current user when running on windows. ([#83](https://github.com/hashicorp/terraform-provider-ad/pull/83))

IMPROVEMENTS:
* **Resource:** `ad_computer`: Add description attribute to resource. ([#85](https://github.com/hashicorp/terraform-provider-ad/pull/85))
* **Resource:** `ad_group`: Add description attribute to resource. ([#93](https://github.com/hashicorp/terraform-provider-ad/pull/93))
* **Resource:** `ad_group, ad_user, ad_computer`: Add a computed field that holds the object's SID. ([#76](https://github.com/hashicorp/terraform-provider-ad/pull/76))
* **provider**: Upgraded the terraform plugin SDK version to 2.5.0
* **provider**: Extract error messages from CLIXML. ([#74](https://github.com/hashicorp/terraform-provider-ad/pull/74))

## 0.4.1 (January 18, 2021)

**BREAKING CHANGES:**

If you are using the `ad_group` or `ad_user` datasources you will have to update some fields in your terraform configuration.

* **Resource:** `ad_group` datasource now use the attribute `group_id` instead of `guid`. ([#69](https://github.com/hashicorp/terraform-provider-ad/pull/69))
* **Resource:** `ad_user` datasource now use the attribute `user_id` instead of `guid`. ([#69](https://github.com/hashicorp/terraform-provider-ad/pull/69))

BUGFIXES:
* **Resource:** `ad_group_membership` uses parameter `Members` instead of `Member`. ([#68](https://github.com/hashicorp/terraform-provider-ad/pull/68))
* **Kerberos Authentication:** Kerberos now respects the protocol setting and correctly uses https when instructed, instead of always using `http`. ([#66](https://github.com/hashicorp/terraform-provider-ad/pull/66))

FEATURES:
* **Resource:** `ad_user`: Added many standard attributes. ([#63](https://github.com/hashicorp/terraform-provider-ad/pull/63))
* **Resource:** `ad_user`: Added support for custom attributes. ([#73](https://github.com/hashicorp/terraform-provider-ad/pull/73))

## 0.4.0 (December 17, 2020)

FEATURES:
* **New Auth method:** The provider now supports Kerberos authentication.
* **New Auth method:** The provider now supports NTLM authentication. ([#56](https://github.com/hashicorp/terraform-provider-ad/pull/56))

## 0.3.0 (November 06, 2020)

FEATURES:
* **New Resource:** `ad_group_membership`

BUGFIXES:
* **Resource:** `ad_user` now supports moving users between containers ([#49](https://github.com/hashicorp/terraform-provider-ad/pull/49))

## 0.2.0 (September 28, 2020)

IMPROVEMENTS:
* Upgraded to the provider SDK v2.0.0 ([#37](https://github.com/hashicorp/terraform-provider-ad/pull/37))

BUGFIXES:
* **Resource:** `ad_gpo_security` now sets the correct Machine Extension Name in the GPO ([#43](https://github.com/hashicorp/terraform-provider-ad/pull/43/))

## 0.1.0 (July 29, 2020)

FEATURES:

* **New Resource:** `ad_user`
* **New Resource:** `ad_group`
* **New Resource:** `ad_computer`
* **New Resource:** `ad_ou`
* **New Resource:** `ad_gpo`
* **New Resource:** `ad_gpo_security`
* **New Resource:** `ad_gplink`

* **New Datasource:**   `ad_user`
* **New Datasource:**   `ad_group`
* **New Datasource:**   `ad_gpo`
* **New Datasource:**   `ad_computer`
* **New Datasource:**   `ad_ou`
