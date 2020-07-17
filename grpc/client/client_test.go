package client

import (
	"testing"
	"google.golang.org/grpc"
	authTypes "github.com/irisnet/irishub-sync/x/auth/types"
	"context"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init()  {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("iaa", "iva")
	//config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	//config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()
}

func TestGRPCClient(t *testing.T)  {
	conn, err := grpc.Dial("127.0.0.1:32781", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		t.Fatal(err)
	} else {
		client := authTypes.NewQueryClient(conn)

		address := "iaa15eh2t9hux3katpcp328dr30adt27yenglhc9mn"
		acc, err := github_com_cosmos_cosmos_sdk_types.AccAddressFromBech32(address)
		if err != nil {
			t.Fatal(err)
		}
		req := authTypes.QueryAccountRequest{
			Address: acc,
		}
		if res, err := client.Account(context.TODO(), &req); err != nil {
			t.Fatal(err)
		} else {

			resBytes, _ := json.Marshal(res)
			t.Log(string(resBytes))
		}
	}
}


