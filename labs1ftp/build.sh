#!/bin/bash
export GOPATH=`pwd`
go get github.com/skorobogatov/input
go get github.com/jlaffaye/ftp
go install ./src/client
