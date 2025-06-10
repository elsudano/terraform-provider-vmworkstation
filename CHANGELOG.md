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

## 2.0.1 (June 15, 2025)

NOTES:

### **Features**  
- **Major Release**: Introduced a new major version with the [Terraform Framework][framework].  
- **VM Management**: Implemented `CreateVM`, `UpdateVM`, and `DeleteVM` methods.  
- **Logging Improvements**: Upgraded logging system for increased verbosity and switched to a new library.  
- **Power Control**: Added a new method to **power on VMs**.  
- **API Client Refactor**: Restructured packages within the API client for better organization.  
- **Debuggin**: Now we are handling more errors and we have improved the messages when Terraform found some issue.

[framework]: https://developer.hashicorp.com/terraform/plugin/framework

### **Fixes**  
- **Parallelism Issues**: Resolved API concurrency problems with **VMWare Workstation REST API**.  
- **API Client**: Fixed incorrect package and method references when calling the API client.  
- **Datasource Handling**: Added VM name-based lookup support in datasources (required due to Terraform limitations).  
- **Documentation**: Fixed broken references in documentation.  
- **Testing**: Adjusted test values to align with the provider’s new behavior.  
- **Build & Tests**: Corrected resource tests and Makefile issues.  

### **Documentation**  
- **Full Docs**: Generated complete documentation for the new release.  
- **Code Clarity**: Added inline comments to improve code readability.  
- **Structure**: Renamed files for future resource support and enhanced documentation behavior.  

### **Maintenance & Improvements**  
- **Dependencies**: Updated Go and third-party libraries.  
- **API Client**: Modified value types in the API client for better compatibility.  
- **Code Cleanup**:  
  - Removed example files.  
  - Renamed provider folder to follow correct naming conventions.  
  - Deleted legacy files from previous versions.  
- **VMWare API**: Upgraded the **VMWare Workstation API client**.  
- **Validation**: Added guardrails to verify datasource correctness.  
- **Legal**: Updated **LICENSE** file.  
- **Release Process**: Enhanced **goreleaser.yaml** with new formats.  
- **Go Upgrade**: Updated to a newer Go version.  

Known Issues:

* We know that we can't change the Denomination and Description of the VM's
  because the API Rest fails, we are working on this issue
  trying to fix it a soon as possible.
* The parallelism, as you know Terraform has the option to create different
  resources at the same time, but the API Rest of VmWare Workstation PRO
  hasn't the option to create multiple resources at the same time, for that
  reason, you will need to use the flag -parallelism=1 when you run the command Terraform.