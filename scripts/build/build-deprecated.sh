#!/bin/bash

ProjectPath=/LiteKube
Version=0.1.0
GitBranch=$(git rev-parse --abbrev-ref HEAD)
GitVersion=$(git version)
GitCommit=$(git rev-parse HEAD)
BuildDate=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
LeaderFold=/LiteKube/build/outputs/leader
WorkerFold=/LiteKube/build/outputs/worker

VersionTags="\
    -X \"github.com/litekube/LiteKube/pkg/version.Litekube=$Version\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitBranch=$GitBranch\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitVersion=$GitVersion\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitCommit=$GitCommit\" \
    -X \"github.com/litekube/LiteKube/pkg/version.BuildDate=$BuildDate\" \
"


mkdir -p $LeaderFold
cd $ProjectPath/cmd/leader

Tag=leader-$(uname)-$(arch)-$Version
echo "build $Tag"
go build -ldflags "$VersionTags" -o $Tag . && mv $Tag $LeaderFold

Tag=leader-Linux-arm-$Version
echo "build $Tag"
env CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -ldflags "$VersionTags" -o $Tag . && mv $Tag $LeaderFold


mkdir -p $WorkerFold
cd $ProjectPath/cmd/worker

Tag=worker-$(uname)-$(arch)-$Version
echo "build $Tag"
go build -ldflags "$VersionTags" -o $Tag . && mv $Tag $WorkerFold
Tag=worker-Linux-arm-$Version
echo "build $Tag"
env CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -ldflags "$VersionTags" -o $Tag . && mv $Tag $WorkerFold