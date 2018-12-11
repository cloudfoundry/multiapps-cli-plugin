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
