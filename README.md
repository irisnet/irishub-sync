# How to use grpc client to query data of cosmoshub or irishub

## Install protoc-gen-gocosmos

1. checkout source code from [cosmos-proto](https://github.com/regen-network/cosmos-proto)
2. exec `make proto-gen`
3. `cd ./protoc-gen-gocosmos & go build`
4. copy `protoc-gen-gocosmos` into path-of-go-bin

## Add third party proto files

1. copy [third_party/proto](https://github.com/cosmos/cosmos-sdk/tree/master/third_party/proto) into current dir

## Gen client code

1. download [protocgen.sh](https://github.com/cosmos/cosmos-sdk/blob/master/scripts/protocgen.sh) into `scripts` dir
2. run `./scripts/protocgen.sh`

## Use GRPCClient connect to GRPCServer

- See detail in `grpc/client/client_test.go`