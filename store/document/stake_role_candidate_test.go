package document

import (
	"encoding/json"
	"testing"

	"github.com/irisnet/irishub-sync/module/logger"
	"github.com/irisnet/irishub-sync/store"
)

func InitDB() {
	store.InitWithAuth()
}

func TestCandidate_GetUnRevokeValidators(t *testing.T) {
	InitDB()

	tests := []struct {
		name string
	}{
		{
			name: "test get unRevoke validators",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Candidate{}
			res, err := d.GetUnRevokeValidators()

			if err != nil {
				logger.Error.Fatalln(err)
			}

			strRes, _ := json.Marshal(res)
			logger.Info.Println(string(strRes))
		})
	}
}

func TestCandidate_RemoveCandidates(t *testing.T) {
	InitDB()

	tests := []struct {
		name string
	}{
		{
			name: "test remove candidates",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Candidate{}
			err := d.RemoveCandidates()
			if err != nil {
				logger.Error.Fatalln(err)
			}
		})
	}
}

func TestCandidate_SaveAll(t *testing.T) {
	InitDB()

	candidates := []Candidate{
		{
			Address:        "1",
			PubKey:         "1",
			PubKeyAddr:     "1",
			Revoked:        true,
			Tokens:         float64(9.0),
			OriginalTokens: "9.0",
			VotingPower:    float64(9.0),
			Description: Description{
				Moniker: "1",
			},
			BondHeight: 1,
		},

		{
			Address:        "2",
			PubKey:         "2",
			PubKeyAddr:     "2",
			Revoked:        false,
			Tokens:         float64(9.0),
			OriginalTokens: "9.0",
			VotingPower:    float64(9.0),
			Description: Description{
				Moniker: "2",
			},
			BondHeight: 2,
		},
	}

	type args struct {
		candidates []Candidate
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test save all candidates",
			args: args{
				candidates: candidates,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Candidate{}
			err := d.SaveAll(tt.args.candidates)
			if err != nil {
				logger.Error.Fatalln(err)
			}
		})
	}
}
