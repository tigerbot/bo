#!/bin/bash

set -e
set -u

# this allows us to build here without conflicting with a more traditional go setup
export GOPATH=$(cd $(dirname "$0") && env pwd -P)
# makes sure the go install command put the results in this directory
unset GOBIN
# allow us to directly use the go utilities we install here
export PATH="${GOPATH}/bin:${PATH}"

# first build the public directory
pushd $GOPATH/browser
npm install
popd

# then conver the public directory into a go
go get -u github.com/jteeuwen/go-bindata/...
go-bindata -o "${GOPATH}/src/bo_server/bindata.go" -prefix "${GOPATH}/public" "${GOPATH}/public/..."

go get bo_server

bo_server
