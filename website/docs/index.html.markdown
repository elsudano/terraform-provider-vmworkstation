---
layout: "vmworkstation"
page_title: "Provider: VMware Workstation Pro"
sidebar_current: "docs-vmworkstation-index"
description: |-
  A Terraform provider to work with VMware Workstation Pro, allowing management of virtual machines and other VMware resources. Supports management through VMware Workstation Pro.
---

# VMware Workstation Provider

The VMware Workstation provider gives Terraform the ability to work with VMware Workstation Pro
Products, notably [VMware Workstation][vmware-worstation].
This provider can be used to manage many aspects of a VMware Worstation Pro
environment, including virtual machines, datastores, and more.

[vmware-workstation]: https://www.vmware.com/products/workstation-pro.html

Use the navigation on the left to read about the various resources and data
sources supported by the provider.

~> **NOTE:** This provider requires API REST enable on VMware Worstation Pro.

  [VmWare Workstation API Rest]: https://github.com/elsudano/vmware-workstation-api-client

## Argument Reference

The provider configuration block accepts the following arguments. In general, it's better to set them via indicated environment variables to keep the configuration safe.

##### 1.- *user*
> (Required) Username that will be used to authenticate in the API REST of VmWare Workstation. It is recommended to be set via VMWS_USER environment variable.

##### 2.- *password*
> (Required) Matcnhing password for the user to authenticate in the API REST of VmWare Workstation. It is recommended to be set via VMWS_PASSWORD environment variable.

##### 3.- *url*
> (Required) This is the URL where the API REST of VmWare Workstation are listen, normally in "https://localhost:8697/api". It is recommended to be set via VMWS_URL environment variable.

##### 4.- *https*
> (Optional) This parameter is false for now the API REST of VmWare Workstation just listen in http, but the comunnication it's encrypted, because to use is necessary created a certificate, maybe in the future VmWare change this, for that, is this variable. It is recommended to be set via VMWS_HTTPS environment variable.

##### 5.- *debug*
> (Optional) The last one is the variable to setting the debug mode, when it's a true, the provider module print in the log of Terraform some actions and is easier find bug. It is recommended to be set via VMWS_DEBUG environment variable.

## Example Usage

In file main.tf:

```HLC
terraform {
  required_version = ">= 0.14.4"
  required_providers {
    vmworkstation = {
      source  = "elsudano/vmworkstation"
      version = "0.2.1"
    }
  }
}
resource "vmworkstation_vm" "test_machine" {
  sourceid     = var.vmws_reource_frontend_sourceid
  denomination = var.vmws_reource_frontend_denomination
  description  = var.vmws_reource_frontend_description
  path         = var.vmws_reource_frontend_path
  processors   = var.vmws_reource_frontend_processors
  memory       = var.vmws_reource_frontend_memory
}
```

In file output.tf

```HLC
output vmws_frontend_id {
  value       = vmworkstation_vm.test_machine.id
  description = "This is the id of the VM"
}
```

In file variables.tf

```HLC
variable "vmws_reource_frontend_sourceid" {
  type        = string
  description = "(Required) The ID of the VM that to use for clone at the new"
}
variable "vmws_reource_frontend_denomination" {
  type        = string
  description = "(Required) The Name of VM in WS"
  default     = "NewInstance"
}
variable "vmws_reource_frontend_description" {
  type        = string
  description = "(Required) The Description at later maybe to explain the instance"
}
variable "vmws_frontend_path" {
  type        = string
  description = "(Required) The Path where will be our instance in VmWare"
}
variable "vmws_reource_frontend_processors" {
  type        = string
  description = "(Required) The number of processors of the Virtual Machine"
  default     = "1"
}
variable "vmws_reource_frontend_memory" {
  type        = string
  description = "(Required) The size of memory to the Virtual Machine"
  default     = "512"
}
```
### Useful commands to handle the infrastructure:

```bash
export VMWS_USER="xxxx"; \
export VMWS_PASSWORD="xxxx"; \
export VMWS_URL="https://localhost:8697/api"; \
ansible-vault decrypt terraform/vault/vmw.tfvars; \
terraform plan -state=terraform/envi/vmw/terraform.tfstate -var-file=terraform/vault/vmw.tfvars terraform/envi/vmw/; \
ansible-vault encrypt terraform/vault/vmw.tfvars
```
### Session persistence options

### Debugging options

To find fastly the outputs of the provider in the Terraform logs you can use the follow command:

```bash
watch -t -n 5 "cat <log file> | grep -e VMWS -e WSAPICLI -e $(date +%Y-%m-%d) | tail -15"
```

## Notes on Required Privileges

### Tags

### Events

### Locating Managed Object IDs

## Bug Reports and Contributing