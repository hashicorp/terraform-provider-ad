---
layout: "ad"
page_title: "ad: ad_domain"
sidebar_current: "docs-ad-domain"
description: |-
  Sample data source in the Terraform provider ad.
---

# ad_domain

Sample data source in the Terraform provider ad.

## Example Usage

```hcl
data "ad_domain" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
