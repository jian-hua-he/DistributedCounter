# Distributed Counter

It is a small distributed service, consisting of multiple micro services (isolated processes) which can count the number of items, grouped by tenants that are delivered through an HTTP restful interface.

The coordinator public API with 2 basic RESTful methods:

```
- POST /items
- GET  /items/{tenant id}/count
```

## Getting Started

All services are running on `docker`. Install `docker` and `docker-compose` before you start it. First, we need to build the bin file by command below:

```
$ make build
```

Then run service by the command below:

```
$ make up
```

The docker will boot all services. Default counter scale is 3. If you want to use a different scale number. You can:

```
$ make up COUNTER_SCALE={num}
```

Stop the services

```
$ make down
```

If you wanna see what happen in services. You can exec the command below:

```
$ make logs
```

The Makefile is simply used `docker-compose` command. You can also use `docker-compose` command to start or stop the services.

## Design
