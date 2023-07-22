#!/bin/sh
set -e

out="output"
app="quicclient"

mkdir -p $out

go version

export CGO_ENABLED=0

export GOOS=linux
export GOARCH=amd64
go build -trimpath -o $out/${app}_${GOOS} ./...

export GOOS=windows
export GOARCH=386
go build -trimpath -o $out/${app}_${GOOS} ./...

export GOOS=darwin
export GOARCH=amd64
go build -trimpath -o $out/${app}_${GOOS} ./...

cd $out
echo sha256sum
sha256sum * | tee SHA256SUMS
