#!/bin/bash

ProjectPath=../../
Outputs=$ProjectPath/build/outputs
Version=0.1.0
GitBranch=$(git rev-parse --abbrev-ref HEAD)
GitVersion=$(git version)
GitCommit=$(git rev-parse HEAD)
BuildDate=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

VersionTags="\
    -X \"github.com/litekube/LiteKube/pkg/version.Litekube=$Version\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitBranch=$GitBranch\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitVersion=$GitVersion\" \
    -X \"github.com/litekube/LiteKube/pkg/version.GitCommit=$GitCommit\" \
    -X \"github.com/litekube/LiteKube/pkg/version.BuildDate=$BuildDate\" \
"

# build for one kind of arch-os
function rungobuild(){
    cc=$1
    codePath=$2
    fileName=$3
    os=$4
    arch=$5
    archTag=$6
    saveDir=$7

    cd $codePath
    mkdir -p $saveDir

    Tag=$fileName-$os-$archTag-$Version
    echo "build $Tag"
    if ! type $cc >/dev/null 2>&1; then
        echo "$cc not install, skip"
    else
        env CGO_ENABLED=1 GOOS=$os GOARCH=$arch CC=$cc go build -ldflags "$VersionTags -w -s" -o $Tag . && mv $Tag  $saveDir/
    fi
}

function compile(){
    codePath=$1
    fileName=$2
    saveDir=$Outputs/$fileName
    mkdir -p $saveDir

    # build by local
    cd $codePath
    Tag=$fileName-$(uname)-$(arch)-$Version
    echo "build $Tag"
    go build -ldflags "$VersionTags -w -s" -o $Tag . && mv $Tag  $saveDir/

    # build for linux-armv7l
    rungobuild arm-linux-gnueabihf-gcc $codePath $fileName linux arm armv7l $saveDir
}


cmdPath=$ProjectPath/cmd
for file in  `ls $cmdPath`
do
    codePath=$cmdPath/$file
    if [ -d $codePath ]
    then
        compile $codePath $file
    fi
done