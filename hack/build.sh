#!/bin/bash
CGO_ENABLED=0
GOOS=linux
#GOOS=darwin
GOARCH=amd64
cd ../
go build . -o kubeutil-linux
