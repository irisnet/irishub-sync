package document

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setup
	store.Start()

	code := m.Run()

	// shutdown
	os.Exit(code)
}

func TestBlock_CountValidatorPreCommits(t *testing.T) {
	type args struct {
		startBlock int64
		endBlock   int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test count validator precommits group by validator addr",
			args: args{
				startBlock: 90006,
				endBlock:   90106,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Block{}
			res, err := d.CalculateValidatorPreCommit(tt.args.startBlock, tt.args.endBlock)
			if err != nil {
				logger.Error(err.Error())
			}
			strRes, _ := json.Marshal(res)
			logger.Info(string(strRes))
		})
	}
}
