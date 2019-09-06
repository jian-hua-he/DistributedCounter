#!/bin/bash

## Test that won’t change the server state
docker-compose up -d

docker run -it --rm \
    -v $(PWD):/go/src/testclient \
    --network $(basename $(dirname $(PWD)))_dis_sys_network \
    golang:1.12 bash -c "cd /go/src/testclient && go test api.go health_test.go none_item_test.go"

docker-compose down

# Test that would change the server state
# Every tests should re-launch the services 
docker-compose up -d

docker run -it --rm \
    -v $(PWD):/go/src/testclient \
    --network $(basename $(dirname $(PWD)))_dis_sys_network \
    golang:1.12 bash -c "cd /go/src/testclient && go test api.go multiple_items_test.go"

docker-compose down