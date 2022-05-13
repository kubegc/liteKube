#!/bin/bash

version=0.1.0

VersionTags="
    -X \"github.com/litekube/LiteKube/pkg/version.Litekube=$version\"
    -X \"github.com/litekube/LiteKube/pkg/version.GitBranch=`git branch`\"
    -X \"github.com/litekube/LiteKube/pkg/version.GitVersion=`git version`\" 
    -X \"github.com/litekube/LiteKube/pkg/version.GitCommit=`git rev-parse HEAD`\"
    -X \"github.com/litekube/LiteKube/pkg/version.BuildDate=`date -u +'%Y-%m-%dT%H:%M:%SZ'`\"
"

mkdir -p /Litekube/build/outputs/leader
cd /Litekube/cmd/leader
go build -ldflags "$VersionTags" -o leader-${uname}-${arch}-$version .
mv leader-linux-amd64-$version /Litekube/build/outputs/leader/
env CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -ldflags "$VersionTags" -o leader-Linux-arm-$version .
mv leader-linux-arm-$version /Litekube/build/outputs/leader/

mkdir -p /Litekube/build/outputs/worker
cd /Litekube/cmd/worker
go build -ldflags "$VersionTags" -o worker-${uname}-${arch}-$version .
mv worker-linux-amd64-$version /Litekube/build/outputs/worker/
env CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -ldflags "$VersionTags" -o worker-Linux-arm-$version .
mv worker-linux-arm-$version /Litekube/build/outputs/worker/
