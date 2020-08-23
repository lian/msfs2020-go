#!/bin/bash
set -e

go generate github.com/lian/msfs2020-go/simconnect
go generate github.com/lian/msfs2020-go/vfrmap

build_time=$(date -u +'%Y-%m-%d_%T')
set +e
build_version=$(git describe --tags)
set -e

[ -f .env ] && source .env
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.mapApiKeyDefault=$MAP_API_KEY -X main.buildVersion=$build_version -X main.buildTime=$build_time" -v github.com/lian/msfs2020-go/vfrmap
