// init mongodb session and provide common functions

package store

import "testing"

func TestInitWithAuth(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "tets initWithAuth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start()
		})
	}
}
