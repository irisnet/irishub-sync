package helper

import (
	"fmt"
	"github.com/irisnet/irishub-sync/types"
	"github.com/pkg/errors"
)

// Query from Tendermint with the provided storename and path
func Query(key types.HexBytes, storeName string, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	client := GetClient()
	defer client.Release()

	opts := types.ABCIQueryOptions{
		Height: 0,
		Prove:  false, //不需要验证prof
	}
	result, err := client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
	}
	return resp.Value, nil
}
//
//func QueryWithPath(key types.HexBytes, path string) (res []byte, err error) {
//	client := GetClient()
//	defer client.Release()
//
//	opts := types.ABCIQueryOptions{
//		Height: 0,
//		Prove:  false, //不需要验证prof
//	}
//	result, err := client.ABCIQueryWithOptions(path, key, opts)
//	if err != nil {
//		return res, err
//	}
//	resp := result.Response
//	if resp.Code != uint32(0) {
//		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
//	}
//	return resp.Value, nil
//}
//
//func QuerySubspace(subspace []byte, storeName string) (res []types.KVPair, err error) {
//	cdc := types.GetCodec()
//	resRaw, err := Query(subspace, storeName, "subspace")
//	if err != nil {
//		return res, err
//	}
//	cdc.MustUnmarshalBinaryBare(resRaw, &res)
//	return
//}
