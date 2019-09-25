package task

import "testing"

func Test_assertFastSyncFinished(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "assert fast sync finished",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := AssertFastSyncFinished()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		})
	}
}

//func TestMakeUpdateDelegatorTask(t *testing.T) {
//	updateDelegator()
//}
