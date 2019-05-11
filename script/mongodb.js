//create database and user
// use sync-iris
// db.createUser(
//     {
//         user:"iris",
//         pwd:"irispassword",
//         roles:[{role:"root",db:"admin"}]
//     }
// )

// create collections
db.createCollection("token_flow");
db.createCollection("block");
db.createCollection("stake_role_candidate");
db.createCollection("sync_task");
db.createCollection("tx_common");
db.createCollection("proposal");
db.createCollection("tx_msg");
db.createCollection("power_change");//explorer
db.createCollection("uptime_change");
db.createCollection("sync_conf");
db.createCollection("mgo_txn");
db.createCollection("mgo_txn.stash");
db.createCollection("ex_tx_num_stat");



// create index
db.token_flow.createIndex({"block_height": -1});
db.account.createIndex({"address": 1}, {"unique": true});
db.block.createIndex({"height": -1}, {"unique": true});

db.stake_role_candidate.createIndex({"address": 1}, {"unique": true});
db.stake_role_candidate.createIndex({"pub_key": 1});

db.sync_task.createIndex({"start_height": 1, "end_height": 1}, {"unique": true});

db.tx_common.createIndex({"height": -1});
db.tx_common.createIndex({"time": -1});
db.tx_common.createIndex({"tx_hash": 1}, {"unique": true});
db.tx_common.createIndex({"from": 1});
db.tx_common.createIndex({"to": 1});
db.tx_common.createIndex({"type": 1});
db.tx_common.createIndex({"status": 1});

db.power_change.createIndex({"height": 1, "address": 1}, {"unique": true});


db.proposal.createIndex({"proposal_id": 1}, {"unique": true});
db.tx_msg.createIndex({"hash": 1}, {"unique": true});
db.ex_tx_num_stat.createIndex({"date": -1}, {"unique": true});

// init data
db.sync_conf.insert({"block_num_per_worker_handle": 50, "max_worker_sleep_time": 120});

// drop collection
// db.account.drop();
// db.block.drop();
// db.power_change.drop();
// db.proposal.drop();
// db.stake_role_candidate.drop();
// db.sync_task.drop();
// db.tx_common.drop();
// db.tx_msg.drop();
// db.mgo_txn.drop();
// db.mgo_txn.stash.drop();

// remove collection data
// db.account.remove({});
// db.block.remove({});
// db.power_change.remove({});
// db.proposal.remove({});
// db.stake_role_candidate.remove({});
// db.sync_task.remove({});
// db.tx_common.remove({});
// db.tx_msg.remove({});
// db.uptime_change.remove({});
// db.mgo_txn.remove({});
// db.mgo_txn.stash.remove({});
