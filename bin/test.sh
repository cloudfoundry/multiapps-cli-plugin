#!/bin/bash

set -e

function printStatus {
  if [ $? -eq 0 ]; then
    echo -e "SUITE SUCCESS"
  else
    echo -e "SUITE FAILURE"
  fi
}

trap printStatus EXIT

if [ ! $(which ginkgo) ]; then
  echo -e "Installing ginkgo..."
  go get github.com/onsi/ginkgo/ginkgo
fi

echo -e "Formatting packages..."
go fmt ./...

echo -e "Vetting packages for potential issues..."
for file in $(find {clients,commands,log,testutil,ui,util} \( -name "*.go" -not -iname "*test.go" \))
do
  go tool vet -all $file
done

echo -e "Testing packages..."
ginkgo -r $@

echo -e "Running build script to confirm everything compiles..."
bin/build.sh
