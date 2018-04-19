#!/bin/bash

function info() {
    echo -e "\033[1;34m$1 \033[0m"
}

function warn() {
    echo  -e "\033[0;33m$1 \033[0m"
}

function error() {
    echo  -e "\033[0;31m$1 \033[0m"
}

function createGoPath() {
    if [ ! -d $1 ];
        then
        mkdir -p $1
    fi
    if [ ! -d "$1/src" ];
    then
        mkdir "$1/src"
    fi
    if [ ! -d "$1/bin" ];
    then
        mkdir "$1/bin"
    fi
    if [ ! -d "$1/pkg" ];
    then
        mkdir "$1/pkg"
    fi
}

if [ ! -n $GOROOT ];
then
    warn "not exists GOROOT"
    exit
fi

gopath="`pwd`"
warn "Use $gopath as golang workspace..."
info "export GOPATH=$gopath"
exportGopath="GOPATH=$gopath"
export $exportGopath
info "go build project"
cd "$gopath/src/winky.com/crawler/"
go install
