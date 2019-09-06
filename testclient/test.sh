#!/bin/bash

docker run -it --rm \
    -v $(PWD):/go/src/testclient \
    --network $(basename $(dirname $(PWD)))_dis_sys_network \
    golang:1.12 bash -c "cd /go/src/testclient && go test ./..."