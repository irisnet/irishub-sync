# IRISHUB-SYNC
A server that synchronize IRIS blockChain data into a database

# Structure

- `conf`: config of project
- `module`: project module
- `mongodb`: mongodb script to create database
- `service`: main logic of sync-server, sync data from blockChain and write to database
- `store`: database model
- `util`: common constants and helper functions
- `main.go`: bootstrap project

# SetUp

## Create mongodb database

run script `mongodb.js` in `mongodb` folder to create database before run project

# Build And Run

- Build: `make all`
- Run: `make run`
- Cross compilation: `make build-linux`

## Env Variables

### Db config

- DB_ADDR: `required` `string` mongodb addrs（example: `127.0.0.1:27017, 127.0.0.2:27017, ...`）
- DB_USER: `required` `string` mongodb user（example: `user`）
- DB_PASSWD: `required` `string` mongodb password（example: `DB_PASSWD`）
- DB_DATABASE：`required` `string` mongodb database name（example：`DB_DATABASE`）

### Server config

- SER_BC_FULL_NODE: `required` `string`  full node url（example: `tcp://127.0.0.1:26657, tcp://127.0.0.2:26657, ...`）
- SER_BC_CHAIN_ID: `required` `string`  chain id（example: `rainbow-dev`）
- WORKER_NUM_CREATE_TASK: `required` `string` num of worker to create tasks（example: `2`）
- WORKER_NUM_EXECUTE_TASK: `required` `string` num of worker to execute tasks（example: `30`）

- NETWORK: `option` `string` network type（example: `testnet,mainnet`）

## Note

If you synchronizes irishub data from specify block height(such as:17908 current time:1576208532)
1. At first, stop the irishub-sync and run follow sql in mongodb 

```
db.sync_task.insert({'start_height':NumberLong(17908),'end_height':NumberLong(0),'current_height':NumberLong(0),'status':'unhandled','last_update_time':NumberLong(1576208532)})
```

2. Then, start irishub-sync

