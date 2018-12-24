package task

import (
	"github.com/irisnet/irishub-sync/store/document"
)

func AssertFastSyncFinished() (bool, error) {
	var (
		syncTaskModel document.SyncTask
		syncConfModel document.SyncConf
	)

	status := []string{document.SyncTaskStatusUnderway}
	tasks, err := syncTaskModel.QueryAll(status, document.SyncTaskTypeFollow)

	if err != nil {
		return false, err
	}

	if len(tasks) != 0 {
		blockChainLatestHeight, err := getBlockChainLatestHeight()
		if err != nil {
			return false, err
		}
		syncConf, err := syncConfModel.GetConf()
		if err != nil {
			return false, err
		}

		task := tasks[0]
		if task.CurrentHeight+syncConf.BlockNumPerWorkerHandle > blockChainLatestHeight {
			return true, err
		}
	}

	return false, nil
}
