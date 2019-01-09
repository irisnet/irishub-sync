package document

import (
	"encoding/json"
	"testing"
)

func TestSyncConf_GetConf(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test get sync config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SyncConf{}
			res, err := d.GetConf()
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(resBytes))
		})
	}
}
