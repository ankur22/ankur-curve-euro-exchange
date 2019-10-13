#!/bin/bash

set -e

GIT_COMMIT="blah"
VERSION=$(cat version.txt)

TARGET_DIR=release/$VERSION

mkdir -p $TARGET_DIR

go build -o exchange-$VERSION -ldflags "-X main.CommitHash=$GIT_COMMIT -X main.Version=$VERSION" github.com/ankur22/ankur-curve-euro-exchange/cmd/exchange-server

mv exchange-$VERSION $TARGET_DIR
