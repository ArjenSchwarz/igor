#!/bin/bash
set -ex

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o main

zip -r igor.zip main index.js config.yml
