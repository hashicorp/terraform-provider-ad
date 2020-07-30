## Building and installing

In order to contribute to the provider, you have to build and install it manually as indicated below. Make sure you have a supported version of Go installed and working. Check out or download this repository, then open a terminal and change to its directory.

### Installing the provider to `terraform.d/plugins`
```
$ make build
$ go install
```
This will build the provider and place the provider binary in your 

You will then have to create a symlink in your[plugins directory](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations) in order for Terraform to detect your provider when you run `terraform init`. Note that the plugin path is different between Terraform versions 0.12 and 0.13.

You are now ready to use the provider. You can find example configurations to test with in this repository under the `./examples` folder.

### Using `-plugin-dir` 

Alternatively, you can run:

```
make build
```

This will place the provider binary in the top level of the provider directory. You can then use it with terraform by specifying the `-plugin-dir` option when running `terraform init`

```
terraform init -plugin-dir /path/to/terraform-provider-ad
```
