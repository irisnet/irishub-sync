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

- DB_ADDR: `required` `string` 数据库连接地址（example: `127.0.0.1:27017, 127.0.0.2:27017, ...`）
- DB_USER: `required` `string` 数据库连接用户（example: `user`）
- DB_PASSWD: `required` `string` 数据库密码（example: `DB_PASSWD`）
- DB_DATABASE：`required` `string` 数据库名称（example：`DB_DATABASE`）

### Server config

- SER_BC_FULL_NODE: `required` `string`  全节点地址（example: `tcp://127.0.0.1:26657, tcp://127.0.0.2:26657, ...`）
- SER_BC_CHAIN_ID: `required` `string`  chain id（example: `rainbow-dev`）
- WORKER_NUM_CREATE_TASK: `required` `string` 执行任务创建的线程数（example: `2`）
- WORKER_NUM_EXECUTE_TASK: `required` `string` 执行任务执行的线程数（example: `30`）

- PrefixAccAddr: `option` `string` 普通用户地址前缀（default: `faa`）
- PrefixAccPub: `option` `string` 普通用户 public key 前缀（default: `fap`）
- PrefixConsAddr: `option` `string` 共识用户地址前缀（default: `fca`）
- PrefixConsPub: `option` `string` 共识用户 public key 前缀（default: `fcp`）
- PrefixValAddr: `option` `string` 验证人地址前缀（default: `fva`）
- PrefixValPub: `option` `string` 验证人 public key 前缀（default: `fvp`）
- CRON_SAVE_VALIDATOR_HISTORY: `option` `string` 保存验证人历史的定时任务（default: `@daily`）
