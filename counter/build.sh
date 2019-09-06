#!/bin/bash

docker run -it -v $(PWD):/go/src/counter --rm golang:1.12 bash -c "cd /go/src/counter && go build -v -o bin/counter"