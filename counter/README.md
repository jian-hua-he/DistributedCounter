# Counter

The Counter is used for store tenant data. The Coordinator will get tenant count from the Counter or send store items to it. when we boot the Counter, it will send register request to the Coordinator to ensure the Coordinator would update data to it. If you want the build bin file. Simply execute `./build.sh`.

## Features

- Sync: The Counter will register to Coordinator when the Counter started. And also sync the data from Coordinator
- Store and retrive item: Store items in memory and allow the Coordinator to retrive it
- Health check: Provide end point for health check

## URLs

```
GET  /items                 # Get all items

GET  /items/{tenant}/count  # Get tenant count

POST /vote                  # First phase of 2PC

POST /commit                # Second phase (commit) of 2PC

POST /rollback              # Second phase (rollback) of 2PC

GET  /health                # For health check
```