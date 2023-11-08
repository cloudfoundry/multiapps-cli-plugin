<p align="center"><img width="335" height="281" src="logo.png" alt="MultiApps logo"></p>

# MultiApps CF CLI Plugin [![Multiapps CLI Plugin build](https://github.com/cloudfoundry/multiapps-cli-plugin/actions/workflows/pull-request-build.yml/badge.svg)](https://github.com/cloudfoundry/multiapps-cli-plugin/actions/workflows/pull-request-build.yml)
This is a Cloud Foundry CLI plugin (formerly known as CF MTA Plugin) for performing operations on [Multitarget Applications (MTAs)](https://www.sap.com/documents/2021/09/66d96898-fa7d-0010-bca6-c68f7e60039b.html) in Cloud Foundry, such as deploying, removing, viewing, etc. It is a client for the [CF MultiApps Controller](https://github.com/cloudfoundry-incubator/multiapps-controller) (known also as CF MTA Deploy Service), which is an MTA deployer implementation for Cloud Foundry. The business logic and actual processing of MTAs happens into CF MultiApps Controller backend.

# Requirements
- Installed CloudFoundry CLI - ensure that CloudFoundry CLI is installed and working. For more information about installation of CloudFoundry CLI, please visit the official CloudFoundry [documentation](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html). You need to have CF CLI v7 or v8 (recommended one)
- Working [CF MultiApps Controller](https://github.com/cloudfoundry-incubator/multiapps-controller) - this a CF plugin for the MultiApps Controller application. Thus, a working MultiApps Controller must be available on the CF landscape

# Download and Installation

:rotating_light: Check whether you have a previous version installed, using the command: `cf plugins`. If the MtaPlugin is already installed, you need to uninstall it first and then to install the new version. You can uninstall the plugin using command `cf uninstall-plugin MtaPlugin`.

## CF Community Plugin Repository

The MultiApps CF CLI Plugin is now also available on the CF Community Repository. To install the latest available version of the MultiApps CLI Plugin execute the following:

`cf install-plugin multiapps`

If you do not have the community repository in your CF CLI you can add it first by executing:

`cf add-plugin-repo CF-Community https://plugins.cloudfoundry.org`

## Manual Installation

Alternatively you can install any version of the plugin by manually downloading it from the releases page and installing the binaries for your specific operating system.

### Download
The latest version of the plugin can also be downloaded from the project's [releases](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/). Download the plugin for your platform (Darwin, Linux, Windows). 


Mac OS X 64 bit | Mac OS X Arm64 | Windows 32 bit | Windows 64 bit | Linux 32 bit | Linux 64 bit | Linux Arm64
--- | --- | --- | --- | --- | --- | ---
[multiapps-plugin.osx](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.osx) | [multiapps-plugin.osxarm64](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.osxarm64) | [multiapps-plugin.win32.exe](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.win32.exe) | [multiapps-plugin.win64.exe](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.win64.exe) | [multiapps-plugin.linux32](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.linux32) | [multiapps-plugin.linux64](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.linux64) | [multiapps-plugin.linuxarm64](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/releases/latest/download/multiapps-plugin.linuxarm64)

### Installation
Install the plugin, using the following command:
```
cf install-plugin <path-to-the-plugin> -f
```
:rotating_light: Note: if you are running on an Unix-based system, you need to make the plugin executable before installing it. In order to achieve this, execute the following commad `chmod +x <path-to-the-plugin>`

## Usage
The MultiApps CF plugin supports the following commands:

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

Here is an example deployment of the open-sourced [spring-music](https://github.com/nvvalchev/spring-music):
```
git clone https://github.com/nvvalchev/spring-music.git
cf deploy mta-assembly/spring-music.mtar -e config.mtaext
```

# Configuration     
The configuration of the MultiApps CF plugin is done via env variables. The following are supported:
* `DEBUG=1` - Enables the logging of HTTP requests in `STDOUT` and `STDERRR`.
* `MULTIAPPS_CONTROLLER_URL=<URL>` - By default, the plugin attempts to deduce the multiapps-controller URL based on the available shared domains. In case of issues, or if you want to target a non-default multiapps-controller instance, you can configure the targeted URL via this env variable.
* `MULTIAPPS_UPLOAD_CHUNK_SIZE=<POSITIVE_INTEGER>` - By default, large MTARs are not uploaded as a single unit, but are split up into smaller chunks of 45 MBs that are uploaded separately. The goal is to prevent failed uploads due to [gorouter](https://github.com/cloudfoundry/gorouter)'s request timeout. In case the default chunk size is still too large, you can configure it via this environment variable. **The specified values are in megabytes.**
:rotating_light: Note: The total number of chunks in which an MTAR is split cannot exceed 50, since the multiapps-controller could interpret bigger values as a denial-of-service attack. For this reason, the minimum value for this environment variable is computed based on the formula: `MIN = MTAR_SIZE / 50`
For example, with a 100MB MTAR the minimum value for this environment variable would be 2, and for a 400MB MTAR it would be 8. Finally, the minimum value cannot grow over 50, so with a 4GB MTAR, the minimum value would be 50 and not 80.
* `MULTIAPPS_UPLOAD_CHUNKS_SEQUENTIALLY=<BOOLEAN>` - By default, MTAR chunks are uploaded in parallel for better performance. In case of a bad internet connection, the option to upload them sequentially will lessen network load.
* `MULTIAPPS_DISABLE_UPLOAD_PROGRESS_BAR=<BOOLEAN>` - By default, the file upload shows a progress bar. In case of CI/CD systems where console text escaping isn't supported, the bar can be disabled to reduce unnecessary logs.

# How to contribute
* [Did you find a bug?](CONTRIBUTING.md#did-you-find-a-bug)
* [Do you have a question or need support?](CONTRIBUTING.md#do-you-have-a-question-or-need-support)
* [How to develop, test and contribute to MultiApps CF Plugin](CONTRIBUTING.md#do-you-want-to-contribute-to-the-code-base)

# Development

*WARNING* : with [Issue 117](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/issues/117) the master branch of this repository as well as other artifacts will be renamed. Adaptation to any CI/CD infrastructure & scritps will be required.

## Cloning the repository
To clone the project in your Go workspace `$GOPATH/src/github.com/cloudfoundry-incubator/multiapps-cli-plugin` execute the following commands:
```
mkdir -p $GOPATH/src/github.com/cloudfoundry-incubator
cd $GOPATH/src/github.com/cloudfoundry-incubator
git clone git@github.com:cloudfoundry-incubator/multiapps-cli-plugin.git
```

## Building new release version
You can automatically build new release for all supported platforms by calling the build.sh script with the version of the build.
The version will be automatically included in the plugin, so it will be reported by `cf plugins`.

:rotating_light: Note that the version parameter should follow the semver format (e.g. 1.2.3).
```bash
./build.sh 1.2.3
```

This will produce `mta_plugin_linux_amd64`, `mta_plugin_darwin_amd64` and `mta_plugin_windows_amd64` in the repo's root directory.

# Further reading
Presentations, documents, and tutorials:
- [Managing Distributed Cloud Native Applications Made Easy (CF Summit EU 2017 slides)](https://www.slideshare.net/NikolayValchev/managing-distributedcloudapps-80697059)
- [Managing Distributed Cloud Native Applications Made Easy (CF Summit EU 2017 video)](https://www.youtube.com/watch?v=d07DZCuUXyk&feature=youtu.be)
- [CF orchestration capabilities by tutorial & example](https://github.com/SAP-samples/cf-mta-examples)

# License

This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the [LICENSE](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/blob/master/LICENSE) file.

