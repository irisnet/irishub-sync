package task

import (
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"time"
)

func MakeValidatorHistoryTask() Task {
	return NewLockTaskFromEnv(server.CronSaveValidatorHistory, "save_validator_history_lock", func() {
		logger.Debug("========================task's trigger [CalculateAndSaveValidatorUpTime] begin===================")
		SaveValidatorHistory()
		logger.Debug("========================task's trigger [CalculateAndSaveValidatorUpTime] end===================")
	})
}

func SaveValidatorHistory() {

	var vHistory []document.ValidatorHistory
	var validatorsModel document.Candidate
	var historyModel document.ValidatorHistory

	validators := validatorsModel.QueryAll()

	updateTime := time.Now()
	for _, v := range validators {
		vHistory = append(vHistory, document.ValidatorHistory{
			Candidate:  v,
			UpdateTime: updateTime,
		})
	}

	if err := historyModel.RemoveAll(); err == nil {
		historyModel.SaveAll(vHistory)
	}
}
