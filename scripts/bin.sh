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

source ${CURRENT_DIR}/init.sh --source-only

BUILD_DIR=${ROOT_DIR}/_build
BUILD=$(git rev-parse --short HEAD || echo "0000000")

build() {
  echo "BUILD ..."
  echo "BUILD: ${BUILD}"

  rm -r ${BUILD_DIR} >/dev/null 2>&1
  mkdir _build

  go build -tags=go_json -o ${BUILD_DIR}/main .

  echo "BUILD=${BUILD}" >> ${BUILD_DIR}/.base.env
  echo "CONFIG_PATH=config.yaml" >> ${BUILD_DIR}/.base.env

  # cp .env ${BUILD_DIR}/
  cp config.yaml ${BUILD_DIR}/
  cp ${SCRIPTPATH}/run.sh ${BUILD_DIR}/

  cd $RUNNING_DIR

  echo "BUILD DONE!"
}

envup() {
  echo "ENVUP ..."

  set -o allexport
    source ${ROOT_DIR}/.env
  set +o allexport

  export BUILD=${BUILD}
  export CONFIG_PATH=${ROOT_DIR}/config.yaml

  echo "ENVUP DONE!"
}

start() {
  echo "START ..."
  echo "BUILD: ${BUILD}"

  envup
  go run -tags=go_json .

  echo "START DONE!"
}

_test() {
  go clean -testcache

  go test -short -p 1 -run=^Test ./... || {
    echo "RUN_TEST FAILED!"
    exit 1
  }
}

test() {
  echo "TEST ..."

  envup

  _test

  echo "RUN_TEST DONE!"
}

code_lint() {
  echo "LINT ..."

  install_golangci_lint
  golangci-lint run $1 --timeout 10m ./...

  echo "LINT DONE!"
}

code_format() {
  echo "FORMAT ..."

  go generate ./...

  gofmt -w internal/ pkg/

  goimports -local tiniapp-backend-oauth-sample -w internal/ pkg/

  echo "FORMAT DONE!"
}

main() {
  case $1 in
    init)
      init
      ;;
    infra)
      infra ${@:2}
      ;;
    build)
      build
      ;;
    test)
      test
      ;;
    start)
      start
      ;;
    code_lint)
      code_lint
      ;;
    code_format)
      code_format
      ;;
    *)
      echo "init|build|start|code_lint|code_format"
      ;;
  esac
}

if [ "${1}" != "--source-only" ]; then
  main "${@}"
fi
