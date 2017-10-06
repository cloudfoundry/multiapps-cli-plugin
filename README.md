# MTA CF CLI Plugin

[![Build Status](TODO: put a link to the original link to the https://travis-ci.com/)

# Description

This is a Cloud Foundry CLI plugin for performing operations on [multi-target applications (MTAs)](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) in Cloud Foundry, such as deploying, removing, viewing, etc. It is a client for the [CF MTA deploy service](https://github.com/SAP/cf-mta-deploy-service), which is an MTA deployer implementation for Cloud Foundry.

# Requirements

- Installed CloudFoundry CLI - ensure that CloudFoundry CLI is installed and working. For more information about installation of CloudFoundry CLI, please visit the official CloudFoundry [documentation](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).
- Working CF deploy-service - this a CF plugin for the deploy-service application. Thus, a working deploy-service must be available on the CF landscape

# Download and Installation

## Download

The latest version of the plugin can be found in the table below. Select the plugin for your platform(Darwin, Linux, Windows) and download it.

Mac OS X 64 bit | Windows 64 bit | Linux 64 bit
--- | --- | ---
[cf-cli-mta-plugin-1.0.4-darwin.bin](TODO: link released libraries from github.com repo) | [cf-cli-mta-plugin-1.0.4-windows.exe](TODO: link released libraries from github.com repo) | [cf-cli-mta-plugin-1.0.4-linux.bin](TODO: link released libraries from github.com repo) |


## Installation

Install the plugin, using the following command:

```
cf install-plugin <path-to-the-downloaded-plugin> -f
```
:rotating_light: Note: if you are running on an Unix-based system, you need to make the plugin executable before installing it. In order to achieve this, execute the following commad `chmod +x <path-to-the-downloaded-plugin>`

:rotating_light: Check whether you have a previous version installed, using the command: `cf plugins`. If the MtaPlugin is already installed, you need to uninstall it first and then to install the new version. You can uninstall the plugin using command `cf uninstall-plugin MtaPlugin`.

# Configuration

The plugin attempts to deduce the deploy service URL based on the CF API URL. In case of issues, or if you want to target a deploy service instance different from the default one, you can configure the targeted deploy service URL via the `DEPLOY_SERVICE_URL` environment variable.

The MTA CF CLI Plugin supports the following commands:

Command Name | Command Description
--- | ---
`deploy` | Deploy a new multi-target app or sync changes to an existing one
`undeploy` | Undeploy (remove) a multi-target app
`mtas` | List all multi-target apps
`mta` | Display health and status for a multi-target app
`mta-ops` | List active multi-target app operations
`download-mta-op-logs` / `dmol` | Download logs of multi-target app operation
`bg-deploy` | Deploy a multi-target app using blue-green deployment
`purge-mta-config` | Purge stale configuration entries

For more information, see the command help output available via `cf [command] --help` or `cf help [command]`.

## Environment Variables
   `DEBUG=1` - Prints in standart output HTTP requests, which are from CLI client to CF deploy service backend.

## Usage

You need to have the MTA CF CLI Plugin installed on your machine.
Here is an example deployment of the open-sourced com.sap.openSAP.hana5.example:
```
git clone https://github.com/SAP/hana-shine-xsa.git
build it
cf deploy assembly/target/hello-world-java-assembly-0.1.0-SNAPSHOT.mtar
```

# How to obtain support

If you need any support, have any question or have found a bug, please report it in the [GitHub bug tracking system](TODO: point to the github.com repostory issues). We shall get back to you.

# Development

## Building new release version
You can automatically build new release for all supported platforms by calling the build.sh script with the version of the build.
The version will be automatically included in the plugin, so it will be reported by `cf plugins`.

:rotating_light: Note that the version parameter should follow the semver format (e.g. 1.2.3).
```bash
./build.sh 1.2.3
```

This will produce `mta_plugin_linux_amd64`, `mta_plugin_darwin_amd64` and `mta_plugin_windows_amd64` in the repo's root directory.

## Adding dependency into the cf-mta-plugin
#### If you want to add a dependecy which to be used later on during the build and development process, you need to follow these steps:
1.  Make sure that you have godep installed(try to run `godep version`). If you do not have it, run the command: `go get github.com/tools/godep`. !!!IMPORTANT!!! Make sure that you are running on latest version of GO and godep!!!
2.  Get the dependency by executing the command: `go get github.com/<package-full-name>` . If you want to update it use the -u option.
3.  Use the dependecy into your code(just import and use)
4.  Make sure that the dependency is not in the Godeps/Godeps.json file(if it is, delete the entry related to it). Godeps.json file is holding information about all the dependencies which are used in the project.
5.  Run `godep save ./...` - this will add all the newly used dependecies into the Godeps/Godeps.json and into the vendor/ folder.

For more information about the godep tool, please refer to: [godep](https://github.com/tools/godep)

# License

This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the [LICENSE](https://github.com/SAP/cf-mta-plugin/blob/master/LICENSE) file.

