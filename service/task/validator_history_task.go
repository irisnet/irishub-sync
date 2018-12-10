package task

import (
	"github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/service/handler"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
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
	validators := helper.GetValidators()

	var vHistory []document.ValidatorHistory
	var historyModel document.ValidatorHistory

	updateTime := time.Now()
	for _, v := range validators {
		candidate := handler.BuildValidatorDocument(v)
		vHistory = append(vHistory, document.ValidatorHistory{
			Candidate:  candidate,
			UpdateTime: updateTime,
		})
	}

	if err := historyModel.RemoveAll(); err == nil {
		historyModel.SaveAll(vHistory)
	}
}
