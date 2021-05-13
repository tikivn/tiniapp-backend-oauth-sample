#!/usr/bin/env sh

if [[ "${DEBUG}" =~ ^1|yes|true ]]; then
  echo "DEBUG=true"
  set -o xtrace
fi

SCRIPTPATH="$(
  cd "$(dirname "$0")"
  pwd -P
)"

CURRENT_DIR=$SCRIPTPATH
ROOT_DIR="$(dirname $CURRENT_DIR)"
PROJECT_NAME="$(basename $ROOT_DIR)"

install_goimports() {
  command -v goimports >/dev/null 2>&1 || {
    echo ""
    echo "project is is installing goimports"
    go get -u golang.org/x/tools/cmd/goimports
  }
}

install_mockgen() {
  command -v mockgen >/dev/null 2>&1 || {
    echo ""
    echo "project is installing mockgen"
    go get -u github.com/golang/mock/mockgen
  }
}

install_golangci_lint() {
  command -v golangci-lint >/dev/null 2>&1 || {
    echo ""
    echo "project is installing golangci-lint"
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.40.1
  }
}

init() {
  install_goimports

  install_mockgen

  install_golangci_lint
}

main() {
  init
}

if [ "${1}" != "--source-only" ]; then
  main "${@}"
fi
