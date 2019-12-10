#!/bin/bash

set -e

echo -e "Installing plugin..."

cd $GOPATH/src/github.com/cloudfoundry-incubator/multiapps-cli-plugin
go install
cf uninstall-plugin multiapps
cf install-plugin $GOPATH/bin/multiapps-cli-plugin -f

