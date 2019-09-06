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

We have a Coordinator and multiple Counters. The Coordinator is responsible to sync data, register and health check in all Counters. Also, get the data from Counters and send it back to the client.

```

              +-------------+
 request ---> | Coordinator |
              +-------------+
                     |
                     |
       +-------------+-------------+
       |             |             |
       V             v             v
  +---------+   +---------+   +---------+
  | Counter |   | Counter |   | Counter |
  +---------+   +---------+   +---------+

```

### Register counter

Every Counter will register to the Coordinator when the Counter start. If we launch a new Counter. It will also register to the Coordinator and sync the data from the Coordinator.

### Sync data

Every Counter will sync the data when it launched. It will send a request to the Coordinator. The Coordinator will retrieve the data from other existed Counter. After getting the data, the Coordinator sends it back to the Counter who send the sync request.

```

         +-------------+
         | Coordinator |
         +-------------+
         ^  |         |
1.       |  | 3.      | 2.
register |  | Send    | Get items 
& sync   |  | items   | data
         |  V         v
  +---------+        +---------+
  |   New   |        | Counter |
  | Counter |        +---------+
  +---------+

```

### Health check

The Coordinator will check the counter every 10 sec. Not sure it is a common duration. If the request attempts over 3 times. The Coordinator will remove the Counter from the registration table.

It used to prevent some Counter is down. And the update request will keep failing since we use 2PC to update data. Remove the failure Counter from the registration table would be helpful in this situation.

### 2 phase commit

I try to implement 2 phase commit(2PC) to keep the data consistent in all Counters. The steps of 2PC are:

1. The Coordinator sends the query to all Counters
2. The Counters return YES or NO to the Coordinator
3. The Coordinator sends the commit/rollback request to the Coordinator
4. The Counters exec commit/rollback and acknowledge to the Coordinator

If any failure happens, like network timeout, Counter failure, etc. It would be considered as a NO answer in the first phase.

```

1. Query to all Counters. All Counter creates a transaction.

          +-------------+
          | Coordinator |
          +-------------+
            |        |
            |        |
            |        |
            V        V
  +---------+        +---------+
  | Counter |        | Counter |
  +---------+        +---------+

2. Votes.

          +-------------+
          | Coordinator |
          +-------------+
            ^        ^
            |        |
     YES/NO |        | YES/NO
            |        |
  +---------+        +---------+
  | Counter |        | Counter |
  +---------+        +---------+

3-1. If any node vote NO, then rollback and remove Transaction

          +-------------+
          | Coordinator |
          +-------------+
            |        |
            |        |
   Rollback |        | Rollback
            V        V
  +---------+        +---------+
  | Counter |        | Counter |
  +---------+        +---------+

3-1. If all nodes vote YES, then commit and remove Transaction

          +-------------+
          | Coordinator |
          +-------------+
            |        |
            |        |
     Commit |        | Commit
            V        V
  +---------+        +---------+
  | Counter |        | Counter |
  +---------+        +---------+

```

But it cannot handle the error occurred during the commit or rollback.

### Query data

I use `docker-compose` to build all services. The query is relay on `docker-compose` network interface. When you send a request to the Coordinator. it will automatically forward your request to the random Counter.

If the single Counter failed. It won’t effect the query.

## PROS

- Stable for query. Some of the Counter down still can do query
- Easy to launch new Counter node
- Event some Counter is down. The data will keep consistent and queryable

## CONS

- If the Coordinator is down. All services should restart. Because the registration table is only kept in memory
- The update will keep failing when one of the Counter is down. It only succeeds when the Counter recovers or the Coordinator remove the host from the registration table
- Update data would be more slowly if we have more Counters node.
