#!/bin/bash

set -e
set -x

GOOS=linux GOARCH=amd64 go build -o main

zip -r igor.zip main index.js config.yml