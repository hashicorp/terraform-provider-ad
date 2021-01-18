## 0.4.1 (Unreleased)

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
