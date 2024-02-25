# Terraform VMware Workstation Pro Provider

In the Terraform registry you can download the latest version of this provider 

https://registry.terraform.io/providers/elsudano/vmworkstation/latest

This is the repository for the Terraform VMware Workstation Pro Provider, which one can use
with Terraform to work with [VMware Workstation][vmware-workstation] Products.

[vmware-workstation]: https://www.vmware.com/products/workstation-pro.html

This provider is inpirated on [vSphere Provider][vmware-vsphere]

[vmware-vsphere]: https://github.com/terraform-providers/terraform-provider-vsphere

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

# Using the Provider

The current version of this provider requires Terraform v0.10.2 or higher to
run.

Note that you need to run `terraform init` to fetch the provider before
deploying. Read about the provider split and other changes to TF v0.15.0 in the
official release announcement found [here][tf-0.15-announce].

[tf-0.15-announce]: https://www.hashicorp.com/blog/announcing-hashicorp-terraform-0-15-general-availability

## Full Provider Documentation

The provider is documented in full on the [GitHub][github-docs] website.
Check the provider documentation for details on entering your connection information and
how to get started with writing configuration for vSphere resources.

[github-docs]: https://github.com/elsudano/terraform-provider-vmworkstation

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`terraform` block in your Terraform configuration if you have not added one
already. Normally this block is in the versions.tf inside your module or your main.tf file.

The syntax is as follows:

```hcl
terraform {
  required_providers {
    vmworkstation = {
      source  = "elsudano/vmworkstation"
      version = "1.0.4"
    }
  }
  required_version = ">= 0.15.4"
}
provider "vmworkstation" {
  user     = var.vmws_user
  password = var.vmws_password
  url      = var.vmws_url
  https    = false
  debug    = true
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.5.0. [Read
more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

# Building The Provider

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/elsudano/terraform-provider-vmworkstation`:

```sh
mkdir -p $GOPATH/src/github.com/elsudano
cd $GOPATH/src/github.com/elsudano
git clone git@github.com:elsudano/terraform-provider-vmworkstation
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/elsudano/terraform-provider-vmworkstation
make build
```
## Debbuging

To debbug the differents issues in the provider, please use the environment variable TF_LOG_PATH

Normaly, you can enable the logs in Terraform setting this two environment vars:

```sh
export TF_LOG=TRACE # INFO, DEBUG, TRACE
export TF_LOG_PATH="terraform.log"
```

Also, you can enable the debug flag for the VMWORKSTATION provider setting the environment var:

```sh
export VMWS_DEBUG=true
```

## Installing the Local Plugin

After the build is complete, you need configure the .terraformrc file in you home for indicate
which is the correct path where is the binary.

Please run the following commands to complete this tasks:

```sh
nano .terraformrc
now put inside these blocks

  dev_overrides {
      "elsudano/vmworkstation" = "/complete/path/where/is/the/binary/file"
  }
  direct {}
```

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/elsudano/terraform-provider-vmworkstation/issues
[gh-prs]: https://github.com/elsudano/terraform-provider-vmworkstation/pulls

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.9+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

See [Building the Provider](#building-the-provider) for details on building the provider.

# Testing the Provider

Please you want try the provider, download from [Registry Provider Link][latest-version]

[latest-version]: https://registry.terraform.io/providers/elsudano/vmworkstation/latest

## Configuring Environment Variables

**NOTE:** Pending
