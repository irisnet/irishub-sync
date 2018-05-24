# IRIS-SYNC-SERVER
A server that synchronize IRIS blockChain data into a database

# Structure

- `conf`: config of project
- `model`: database model
- `module`: project module
- `mongodb`: mongodb script to create database
- `sync`: main logic of sync-server, sync data from blockChain and write to database
- `util`: common constants and helper functions
- `main.go`: bootstrap project

# SetUp

## Rewrite config file

1. rename `/conf/db/type.go.example` to `/conf/db/type.go`, `/conf/server/type.go.example` to `/conf/db/type.go`
2. write your own config into `/conf/db/type.go` and `/conf/db/type.go`

## Create mongodb database

run script `mongodb.js` in `mongodb` folder to create database before run project

# Build And Run

- Build: `make all`
- Run: `make run`
- Cross compilation: `make build-linux`