DOCKER_IMAGE ?= golang:1.12
COUNTER_SCALE ?= 3

.PHONY: build
build:
	cd $(PWD)/coordinator && ./build.sh
	cd $(PWD)/counter && ./build.sh

.PHONY: up
up:
	docker-compose up -d --scale counter=$(COUNTER_SCALE)

.PHONY: down
down:
	docker-compose down

.PHONY: logs
logs:
	docker-compose logs -f coordinator counter