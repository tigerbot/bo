#!/bin/bash
set -e
set -u

function usage() {
  echo -e "Usage:\tbuild.sh [flags...]\n"
  echo "The flags are:"
  echo "  -h, --help        show this help message and quit"
  echo "  -v, --verbose     show additional output"
  echo "  -s, --short       pass the short flag into the go test to skip longer tests"
}

TEST_FLAGS="-cover"
function setArgs() {
  while [ "${1:-}" != "" ]; do
    case $1 in
      "-h" | "--help")
        usage
        exit 0
        ;;
      "-s" | "--short")
        shift
        TEST_FLAGS="${TEST_FLAGS:-} -short"
        ;;
      "-v" | "--verbose")
        shift
        set -x
        TEST_FLAGS="${TEST_FLAGS:-} -v"
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

setArgs "$@"

# this allows the user to have other Go setups and still build here
export GOPATH="$(cd "$(dirname "$0")" && env pwd -P)"
# make sure it puts the bin/ folder in this directory
unset GOBIN

# not all of our packages currently have tests, so look for any directory that has
# files ending in _test.go and only try to test those packages. We do want to know
# which packages we aren't testing though.
UNTESTED="$(find ${GOPATH}/src -name "*.go"      | xargs -n1 dirname | sort | uniq)"
TESTABLE="$(find ${GOPATH}/src -name "*_test.go" | xargs -n1 dirname | sort | uniq)"

# filter out any packages that aren't ours by getting rid of things in our .gitignore file
for ignored in $(grep '^src/' ${GOPATH}/.gitignore); do
  UNTESTED=$(grep -v ${GOPATH}/${ignored} <<< "${UNTESTED}")
  TESTABLE=$(grep -v ${GOPATH}/${ignored} <<< "${TESTABLE}")
done

# and now strip the ${GOPATH}/src prefix to produce names we can actually use with go commands
UNTESTED=$(sed "s|${GOPATH}/src/||g" <<< "${UNTESTED}")
TESTABLE=$(sed "s|${GOPATH}/src/||g" <<< "${TESTABLE}")

FAILED=""
for pkg in ${TESTABLE}; do
  go get -d -t ${pkg} || true
  go test ${TEST_FLAGS:-} ${pkg} || FAILED="$(echo -e "${FAILED}\n\t${pkg}")"
done

FAILED_CNT=$[$(wc -l <<< "${FAILED}")-1]
if [ ${FAILED_CNT} -gt 0 ]; then
  echo -e "\n" >&2
  echo "${FAILED_CNT} packages failed: ${FAILED}" >&2
  exit 1
fi

for pkg in ${TESTABLE}; do
  UNTESTED=$(grep -v ${pkg} <<< "${UNTESTED}")
done
UNTESTED_CNT=$(wc -l <<< "${UNTESTED}")
if [ ${UNTESTED_CNT} -gt 0 ]; then
  echo -e "\n${UNTESTED_CNT} package left untested: ${UNTESTED}"
fi
