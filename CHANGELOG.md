#CHANGELOG

## 0.0.1 (August 15, 2019)

NOTES:

* Was generate this documentation and the first structure of folders for the provider module.

## 0.1.9 (January 05, 2021)

NOTES:

* We add some functions and refactoring
* We change the method to use with the API VmWare Version 16.02

## 0.2.0 (April 18, 2021)

NOTES:

* We have solve some issues  
* Also, we have added a new functionality to register the VM in the GUI of VmWare
* Completed methods to create VM instances  
* Refactored `CreateVM` and `ReadVM` to read directly from VMX file  
* Updated API client and publication flow

## 0.2.1 (December 11, 2021)

NOTES:

* We have completed the methods that we will use to create the instances
* Also, we have modified the methods CreateVM and ReadVM in order to read directly of the VMX file
* Moreover, we have modified the manner for publish the new version of the provider
* and finally, we have modified the internal API client that request information at VmWare Workstation API Rest

## 0.2.2 (April 18, 2021)

NOTES:

* Updated documentation  
* Solved issues with Terraform registry

## 0.2.3 (June 04, 2021)

NOTES:

* Documentation improvements and PowerShell Makefile added  
* Introduced Windows support in build process

## 0.2.4 (June 05, 2021)

NOTES:

* Added `goreleaser` for provider release process

## 1.0.2 (December 11, 2021)

NOTES:

* Enhanced Makefile and integrated `goreleaser`  
* Upgraded external tools and Go mod commands  
* Improved internal API client

## 1.0.3 (December 11, 2021)

NOTES:

* Refactored Makefile with best practices  
* Manual release process stabilized

## 1.0.4 (February 25, 2024)

NOTES:

* Fixed deprecated flags and typos  
* Improved DEBUG mode and error handling  
* Solved Windows-specific issues

## 1.1.5 (March 14, 2024)

NOTES:

* Added IP and power state info to `ReadVM`  
* Improved update method to support multi-field changes  
* Optimized Makefile and dependency management

## 1.1.6 (March 24, 2024)

NOTES:

* Improved publish task and wait status for VM readiness  
* Removed unnecessary register step on VM creation  
* Upgraded dependencies

## 1.2.6 (May 21, 2025)

NOTES:

* Added method to power on VM  
* Updated `goreleaser` to latest formats
* Introduced a new control point for enhanced stability
