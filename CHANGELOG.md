## v2.5.0
* Added --namespace flag (experimental) to 'deploy'/'undeploy' and related commands

## v2.4.1
* Fixed normalization of paths when creating MANIFEST.MF

## v2.4.0

* Added support for non-normalized paths in descriptor
* Additional --mta option for 'dmol' and 'mta-ops'
* Logging the operation id in the beginning of the (un)deployment
* Add --delete-service-keys option for the 'undeploy' command
* Fix overriding of artifacts with the same filename

## v2.3.1

* Fixed deploy help message
## v2.3.0
* Strategy flag for deploy
* Link project to Kanban Board
* Modified build.sh to accommodate static builds

## v2.2.1
* Send mtaId when deployment starts
## v2.2.0
* Introduce configurable retry on operation failure

## v2.1.3
* Add validation for env variable `CHUNK_SIZE_IN_MB`. The minimum value is computed based on the formula: `MIN = MTAR_SIZE / 50`. The maximum value is 50
* Fix backend URL discovery when `-u` option is used

## v2.1.2
* Avoid duplication of error messages in output when process fails
* Allow users to verify archive signature via `--verify-archive-signature` optional parameter

## v2.1.1
* Fixed a DNS lookup timeout issue experienced by some users

## v2.1.0
* Prepare for adoption in [CF-Community](https://github.com/cloudfoundry/cli-plugin-repo) plugins repo
* Rename plugin name: MtaPlugin -> multiapps
* Add builds for linux32 and win32 platforms

## v2.0.13
* Large MTARs are not uploaded as a single unit, but are rather split up into smaller chunks that are uploaded separately. This is done in order to prevent failed uploads due to gorouter's request timeout.
The chunk's size is now configurable through the env variable CHUNK_SIZE_IN_MB. The value of the variable must be a positive integer and the default is 45. Smaller size may be preferable for slower internet connections.

## v2.0.12
* Fix selective deployment on Windows
* Fix selective deployment with modules sharing the same binary
* Stop DDoS attacks against the multiapps-controller
* Retry on "Invalid CSRF token" errors

## v2.0.11
* Fix 'cf mta' command when there are non-staged apps

## v2.0.10
* Add support for all-modules and all-resources

## v2.0.9
* Add support for selective deployment

## v2.0.8
* Increase TLS Handshake timeout
* Remove deploy attributes from /mtas API
* Remove no-longer supported process parameter
* Allow users to skip the ownership validation via `--skip-ownership-validation` optional parameter

## v2.0.7
* Refactor progress output
* Display error messages from non-successful REST calls
* Fix an issue where deployment was not possible in space with a lot of operations

## v2.0.6

* Show reason for failed uploads

## v2.0.1

* Fix computation of deploy service URL

## v2.0.0

* Bump version to 2.0.0
* Print dmol command for finished and aborted processes
* Print the ID of the monitored process
* Add support for --abort-on-error option
* Change Message to Text in models.Message
* Add support for retryable mta rest client
* Introduce function for getting deploy-service URL
* Add support for providing session tokens
* Remove unused fakes
* Re-generate the client for log/content
* Refactor service id provider
* Fix errors in commands
* Remove non-used methods from restclient
* Replace slmp and slpp clients with mta client in commands
* Delete slppclient and slmpclient
* Update version of go-openapi
* Add implementation details to the new client
* Add auth info
* Add method for executing actions
* Add mta_rest.yaml as a client definition
* Print warning when using a custom deploy service URL
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md

## Initial public version 1.0.5

* Supported MTA Operations:
    * deploy - Deploy MTA
    * undeploy - Undeploy MTA
    * bg-deploy - Deploy MTA using blue-green approach
    * mta/mtas - List existing MTA/MTAs
    * mta-ops - Show MTA operations
    * download-mta-op-logs - Download process logs for MTA
