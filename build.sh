#!/bin/bash

set -e
set -u

# this allows us to build here without conflicting with a more traditional go setup
export GOPATH=$(cd $(dirname "$0") && env pwd -P)
# makes sure the go install command put the results in this directory
unset GOBIN
# allow us to directly use the go utilities we install here
export PATH="${GOPATH}/bin:${PATH}"

_dev=false

function usage() {
  echo -e "Usage:\tbuild.sh [flags...]\n"
  echo "The flags are:"
  echo "  -h, --help        show this help message and quit"
  echo "  -d, --dev         don't run gulp, run go-bindata with, start the server"
}

function set_args() {
  while [ "${1:-}" != "" ]; do
    case $1 in
      "-h" | "--help")
        usage
        exit 0
        ;;
      "-d" | "--dev")
        shift
        _dev=true
        flag="-debug"
        ;;
      --)
        shift
        break
        ;;
      *)
        echo "Unrecognized argument: $1"
        usage
        exit 1
        ;;
    esac
  done
}

set_args "$@"

# first build the public directory (unless we are in dev mode, in which case we'll just call
# gulp manually when needed)
if ! $_dev; then
	pushd $GOPATH/browser
	npm install
	popd
fi

# then package the public directory into a go file in the main package.
go get -u github.com/jteeuwen/go-bindata/...
go-bindata ${flag:-} -o "${GOPATH}/src/bo_server/bindata.go" -prefix "${GOPATH}/public" "${GOPATH}/public/..."

go get bo_server

if $_dev; then
	echo "starting server"
	bo_server
fi
