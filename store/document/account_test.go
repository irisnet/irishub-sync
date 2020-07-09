package document

import (
	"encoding/json"
	"testing"
)

func TestGetAccountByPK(t *testing.T) {
	doc := Account{
		Address: "123",
	}
	res, err := doc.getAccountByPK()
	if err != nil {
		t.Fatal(err)
	}
	resBytes, _ := json.MarshalIndent(res, "", "\t")
	t.Log(string(resBytes))
}

//
//func TestAccount_UpsertBalanceInfo(t *testing.T) {
//	type args struct {
//		address       string
//		balance       store.Coin
//		accountNumber uint64
//		height        int64
//		timestamp     int64
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "TestAccount_UpsertBalanceInfo",
//			args: args{
//				address: "12445",
//				balance: store.Coin{
//					Denom:  constant.IrisAttoUnit,
//					Amount: float64(1),
//				},
//				accountNumber: 1,
//				height:        121,
//				timestamp:     time.Now().Unix(),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := Account{}
//			if err := d.SaveBalanceInfo(tt.args.address, tt.args.balance, tt.args.accountNumber, tt.args.height, tt.args.timestamp); err != nil {
//				t.Fatal(err)
//			}
//		})
//	}
//}
//
//func TestAccount_SaveDelegationInfo(t *testing.T) {
//	type fields struct {
//	}
//	type args struct {
//		address    string
//		delegation store.Coin
//		height     int64
//		timestamp  int64
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := Account{}
//			if err := d.SaveDelegationInfo(tt.args.address, tt.args.delegation, tt.args.height, tt.args.timestamp); (err != nil) != tt.wantErr {
//				t.Errorf("Account.SaveDelegationInfo() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAccount_UpsertUnbondingDelegationInfo(t *testing.T) {
//	type fields struct {
//	}
//	type args struct {
//		address             string
//		unbondingDelegation store.Coin
//		height              int64
//		timestamp           int64
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := Account{}
//			if err := d.SaveUnbondingDelegationInfo(tt.args.address, tt.args.unbondingDelegation, tt.args.height, tt.args.timestamp); (err != nil) != tt.wantErr {
//				t.Errorf("Account.SaveUnbondingDelegationInfo() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
