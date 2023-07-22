#!/bin/sh
set -e

out="output"
mkdir -p $out/windows $out/linux $out/darwin

go version

export CGO_ENABLED=0

export GOOS=linux
export GOARCH=amd64
go build -trimpath -o $out/${GOOS}/quicclient ./...

export GOOS=windows
export GOARCH=386
go build -trimpath -o $out/${GOOS}/quicclient ./...

export GOOS=darwin
export GOARCH=amd64
go build -trimpath -o $out/${GOOS}/quicclient ./...

cd $out
echo sha256sum
sha256sum */*
