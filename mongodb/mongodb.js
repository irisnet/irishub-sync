// create collections
db.createCollection("account");
db.createCollection("block");
db.createCollection("stake_role_candidate");
db.createCollection("stake_role_delegator");
db.createCollection("sync_task");
db.createCollection("tx_coin");
db.createCollection("tx_stake");

// create index
db.account.createIndex({"address": 1}, {"unique": true});
db.block.createIndex({"height": -1}, {"unique": true});

db.stake_role_candidate.createIndex({"pub_key": 1}, {"unique": true});
db.stake_role_candidate.createIndex({"address": 1});

db.stake_role_delegator.createIndex({"pub_key": 1});
db.stake_role_delegator.createIndex({"address": 1});
db.stake_role_delegator.createIndex({"address": 1, "pub_key": 1}, {"unique": true});

db.sync_task.createIndex({"chain_id": 1}, {"unique": true});

db.tx_coin.createIndex({"tx_hash": 1}, {"unique": true});
db.tx_coin.createIndex({"from": 1});
db.tx_coin.createIndex({"to": 1});
db.tx_coin.createIndex({"height": -1});
db.tx_coin.createIndex({"time": -1});
db.tx_coin.createIndex({"from": 1, "to": 1, "time": 1});

db.tx_stake.createIndex({"tx_hash": 1});
db.tx_stake.createIndex({"stake_tx.tx_hash": 1});
db.tx_stake.createIndex({"description.moniker": 1});
db.tx_stake.createIndex({"from": 1});
db.tx_stake.createIndex({"pub_key": 1});
db.tx_stake.createIndex({"height": -1});
db.tx_stake.createIndex({"type": 1});
db.tx_stake.createIndex({"time": -1});
db.tx_stake.createIndex({"from": 1, "to": 1, "type": 1, "time": -1});


// drop collection
db.account.drop();
db.block.drop();
db.stake_role_candidate.drop();
db.stake_role_delegator.drop();
db.sync_task.drop();
db.tx_coin.drop();
db.tx_stake.drop();

// remove collection data
db.account.remove({});
db.block.remove({});
db.stake_role_candidate.remove({});
db.stake_role_delegator.remove({});
db.sync_task.remove({});
db.tx_coin.remove({});
db.tx_stake.remove({});










