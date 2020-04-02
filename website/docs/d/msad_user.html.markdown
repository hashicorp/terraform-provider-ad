---
layout: "msad"
page_title: "msad: msad_user"
sidebar_current: "docs-msad-user"
description: |-
  Sample data source in the Terraform provider msad.
---

# msad_user

Sample data source in the Terraform provider msad.

## Example Usage

```hcl
data "msad_user" "example" {
  sample_attribute = "foo"
}
```

## Attributes Reference

* `sample_attribute` - Sample attribute.
