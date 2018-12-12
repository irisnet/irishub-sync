package service

const (
	goroutineNumCreateTask = 2
)

var (
	blockNumPerWorker int64
)

func init() {
	// TODO: query sync_config to get value of block_num_per_worker
	blockNumPerWorker = int64(100)
}

func CreateSyncTask() {
	// check max block number in sync_task collection
	// TODO: query sync_task collection
	maxBlockNum := int64(1000)

	// get current block height
	// TODO: use http client get from tendermint
	currentBlockHeight := int64(20000)

	// create sync task
	for maxBlockNum+blockNumPerWorker <= currentBlockHeight {

	}

}
