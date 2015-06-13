#!/bin/sh

if [ "$1" = "internal" ]; then
  go get bitbucket.org/chrj/smtpd
  export GOARCH=amd64 
  for platform in windows linux darwin; do
    suffix=""
    if [ "$platform" = "windows" ]; then
      suffix=".exe"
    fi
    export GOOS=$platform
    go build -o fakemail-${GOOS}-${GOARCH}${suffix}
  done
else
  docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:cross /usr/src/myapp/build.sh internal
fi
