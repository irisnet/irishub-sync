package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/util/helper"
)

// calculate and save validator upTime
// latest x blocks, calculate how much precommit which validator had execute(n).
// so upTime is n / x
// note: this method is not goroutine safety, it should be execute during watch block.
func calculateAndSaveValidatorUpTime() {
	var (
		methodName    = "AnalyzeValidatorUpTime"
		intervalBlock = constant.IntervalBlockNumCalculateValidatorUpTime
		blockModel    document.Block
		syncTaskModel document.SyncTask
		model         document.ValidatorUpTime
		valUpTimes    []document.ValidatorUpTime
	)
	logger.Info("Start", logger.String("method", methodName))
	// query synced latest height
	tasks, err := syncTaskModel.QueryAll([]string{document.SyncTaskStatusUnderway},
		document.SyncTaskTypeFollow)
	if err != nil {
		logger.Error("Query follow task failed", logger.String("err", err.Error()))
		return
	}
	if len(tasks) == 0 {
		logger.Error("There is no follow task")
		return
	}
	latestHeight := tasks[0].CurrentHeight

	// get validator precommit
	res, err := blockModel.CalculateValidatorPreCommit(latestHeight-intervalBlock, latestHeight)

	if err != nil {
		logger.Error("blockModel.CalculateValidatorPreCommit fail", logger.String("err", err.Error()))
		return
	}

	if len(res) > 0 {
		for _, v := range res {
			tmp := float64(v.PreCommitsNum) / float64(intervalBlock) //注意必须是浮点数相除
			valUpTime := document.ValidatorUpTime{
				ValAddress: v.Address,
				UpTime:     helper.RoundFloat(tmp*100, 0),
			}
			valUpTimes = append(valUpTimes, valUpTime)
		}

		// remove all data
		err := model.RemoveAll()
		if err != nil {
			logger.Error("RemoveAll fail", logger.String("err", err.Error()))
			return
		}

		// save latest data
		err2 := model.SaveAll(valUpTimes)
		if err2 != nil {
			logger.Error("SaveAll fail", logger.String("err", err2.Error()))
			return
		}
	}

	logger.Info("End", logger.String("method", methodName))
}

func MakeCalculateAndSaveValidatorUpTimeTask() Task {
	return NewLockTaskFromEnv(conf.CronCalculateUpTime, "calculate_and_save_validator_uptime_lock", func() {
		logger.Debug("========================task's trigger [CalculateAndSaveValidatorUpTime] begin===================")
		calculateAndSaveValidatorUpTime()
		logger.Debug("========================task's trigger [CalculateAndSaveValidatorUpTime] end===================")
	})
}
