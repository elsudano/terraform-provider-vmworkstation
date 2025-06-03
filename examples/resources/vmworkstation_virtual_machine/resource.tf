resource "vmworkstation_resource_vm" "vm" {
  sourceid     = data.vmworkstation_datasource_vm.parentvm.id
  denomination = "test01"
  description  = "This VM is just for testing purpose"
  path         = "D:\\VirtualMachines\\test01\\test01.vmx" // Windows annotation
  processors   = 2
  memory       = 1024
  state        = "off"
}
