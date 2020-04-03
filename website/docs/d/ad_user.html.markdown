---
layout: "ad"
page_title: "ad: ad_user"
sidebar_current: "docs-ad-user"
description: |-
  Sample data source in the Terraform provider ad.
---

# ad_user

Sample data source in the Terraform provider ad.

## Example Usage

```hcl
data "ad_user" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
