#!/bin/bash

docker run -it -v $(PWD):/go/src/coordinator --rm golang:1.12 bash -c "cd /go/src/coordinator && go test ./... --race"