// create collections
db.createCollection("account");
db.createCollection("block");
db.createCollection("stake_role_candidate");
db.createCollection("stake_role_delegator");
db.createCollection("sync_task");
db.createCollection("tx_stake");
db.createCollection("tx_common");

// create index
db.account.createIndex({"address": 1}, {"unique": true});
db.block.createIndex({"height": -1}, {"unique": true});

db.stake_role_candidate.createIndex({"address": 1}, {"unique": true});
db.stake_role_candidate.createIndex({"pub_key": 1});

db.stake_role_delegator.createIndex({"pub_key": 1});
db.stake_role_delegator.createIndex({"address": 1});
db.stake_role_delegator.createIndex({"address": 1, "pub_key": 1}, {"unique": true});

db.sync_task.createIndex({"chain_id": 1}, {"unique": true});

db.tx_stake.createIndex({"height": -1});
db.tx_stake.createIndex({"time": -1});
db.tx_stake.createIndex({"tx_hash": 1});
db.tx_stake.createIndex({"stake_tx.tx_hash": 1});
db.tx_stake.createIndex({"description.moniker": 1});
db.tx_stake.createIndex({"from": 1});
db.tx_stake.createIndex({"to": 1});
db.tx_stake.createIndex({"pub_key": 1});
db.tx_stake.createIndex({"type": 1});
db.tx_stake.createIndex({"status": 1});
db.tx_stake.createIndex({"from": 1, "to": 1, "type": 1, "status": 1, "time": -1});

db.tx_common.createIndex({"height": -1});
db.tx_common.createIndex({"time": -1});
db.tx_common.createIndex({"tx_hash": 1}, {"unique": true});
db.tx_common.createIndex({"from": 1});
db.tx_common.createIndex({"to": 1});
db.tx_common.createIndex({"type": 1});
db.tx_common.createIndex({"status": 1});

// drop collection
db.account.drop();
db.block.drop();
db.stake_role_candidate.drop();
db.stake_role_delegator.drop();
db.sync_task.drop();
db.tx_stake.drop();
db.tx_common.drop();

// remove collection data
db.account.remove({});
db.block.remove({});
db.stake_role_candidate.remove({});
db.stake_role_delegator.remove({});
db.sync_task.remove({});
db.tx_stake.remove({});
db.tx_common.remove({});










