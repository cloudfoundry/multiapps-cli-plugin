#!/bin/bash

set -e

echo -e "Installing plugin..."

go install
cf uninstall-plugin MtaPlugin
cf install-plugin $GOPATH/bin/cf-cli-plugin -f
