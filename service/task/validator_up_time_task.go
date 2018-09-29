package task

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/constant"
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
		model         document.ValidatorUpTime
		valUpTimes    []document.ValidatorUpTime
	)
	logger.Info.Printf("%v: Start\n", methodName)

	// query synced latest height
	syncTask, _ := document.QuerySyncTask()
	latestHeight := syncTask.Height

	// get validator precommit
	res, err := blockModel.CalculateValidatorPreCommit(latestHeight-intervalBlock, latestHeight)

	if err != nil {
		logger.Error.Printf("%v: Query fail, err is %v\n", methodName, err)
		return
	}

	if len(res) > 0 {
		for _, v := range res {
			valUpTime := document.ValidatorUpTime{
				ValAddress: v.Address,
				UpTime:     float64(v.PreCommitsNum) / float64(intervalBlock) * 100,
			}
			valUpTimes = append(valUpTimes, valUpTime)
		}

		// remove all data
		err := model.RemoveAll()
		if err != nil {
			logger.Error.Printf("%v: Remove all data fail, err is %v\n",
				methodName, err)
			return
		}

		// save latest data
		err2 := model.SaveAll(valUpTimes)
		if err2 != nil {
			logger.Error.Printf("%v: Save latest data fail, err is %v\n",
				methodName, err2)
			return
		}
	}

	logger.Info.Printf("%v: End\n", methodName)
}

func MakeCalculateAndSaveValidatorUpTimeTask() Task {
	return NewLockTaskFromEnv(conf.CronCalculateUpTime, "calculate_and_save_validator_uptime_lock", func() {
		logger.Info.Printf("========================task's trigger [%s] begin===================", "CalculateAndSaveValidatorUpTime")
		calculateAndSaveValidatorUpTime()
		logger.Info.Printf("========================task's trigger [%s] end===================", "CalculateAndSaveValidatorUpTime")
	})
}
