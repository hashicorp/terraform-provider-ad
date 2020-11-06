## 0.3.0 (Unreleased)

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
