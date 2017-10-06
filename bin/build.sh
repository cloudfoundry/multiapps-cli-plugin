#!/bin/bash

set -e

echo -e "Generating binaries..."

ROOT_DIR=$(cd $(dirname $(dirname $0)) && pwd)

rm -rf $ROOT_DIR/out/
GOOS=darwin GOARCH=amd64 go build -o $ROOT_DIR/out/mta_plugin_darwin_amd64
GOOS=linux GOARCH=amd64 go build -o $ROOT_DIR/out/mta_plugin_linux_amd64
GOOS=windows GOARCH=amd64 go build -o $ROOT_DIR/out/mta_plugin_windows_amd64.exe
