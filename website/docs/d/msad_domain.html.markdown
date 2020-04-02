---
layout: "msad"
page_title: "msad: msad_domain"
sidebar_current: "docs-msad-domain"
description: |-
  Sample data source in the Terraform provider msad.
---

# msad_domain

Sample data source in the Terraform provider msad.

## Example Usage

```hcl
data "msad_domain" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
