# IRIS-SYNC-SERVER
A server that synchronize IRIS blockChain data into a database

# Structure

- `conf`: config of project
- `model`: database model
- `module`: project module
- `sync`: main logic of sync-server, sync data from blockChain and write to database
- `util`: common constants and helper functions
- `main.go`: bootstrap project

# Build And Run

- Build: `make all`
- Run: `make run`
- Cross compilation: `make build-linux` or `make docker-build`