# Coordinator

The Coordinator helps all Counters to sync data. Get data from all Counters. If you want the build bin file. Simply execute `./build.sh`.

## Features

- Register Counter: Register all the Counters and provide synced data to them
- Health check: Check all Counters health. If one of Counter failed, remove the Counter from the register
- Update items: Update all items by using 2PC implementation.
- Get tenant count: Get tenant count from Counters

## Opened urls

```
POST /items                 # Update items

GET  /items/{tenant}/count  # Get tenant count from Counter

POST /register              # Register Counter into Coordinator. 
                            # Only accept request from the Count
```
