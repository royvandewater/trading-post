#!/bin/bash

APP_NAME=trading-post
BUILD_DIR=$PWD/dist

build_on_local() {
  local goos="$1"
  local goarch="$2"
  local extension=""

  if [ "$goos" == "windows" ]; then
    extension=".exe"
  fi

  env GOOS="$goos" GOARCH="$goarch" go build -a -tags netgo -installsuffix cgo -ldflags '-w' -o "${BUILD_DIR}/${APP_NAME}-${goos}-${goarch}${extension}" .
}

build_osx_on_local() {
  build_on_local "darwin" "amd64"
}

build_windows_amd64_on_local() {
  build_on_local "windows" "amd64"
}

init() {
  rm -rf $BUILD_DIR/ \
  && mkdir -p $BUILD_DIR/
}

fatal() {
  local message=$1
  echo $message
  exit 1
}

build_all(){
  for goos in darwin linux windows; do
    for goarch in 386 amd64; do
      build_on_local "$goos" "$goarch" > /dev/null
    done
  done
}

main() {
  local mode="$1"

  if [ "$mode" == "all" ]; then
    build_all
    exit $?
  fi

  if [ "$mode" == "osx" ]; then
    build_osx_on_local
    exit $?
  fi

  if [ "$mode" == "windows" ]; then
    build_windows_amd64_on_local
    exit $?
  fi

  echo "Usage: ./build.sh <osx/windows/all>"
  exit 1
}
main $@
