module github.com/irisnet/irishub-sync

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200825201020-d9fd4d2ca9a3
	github.com/go-kit/kit v0.10.0
	github.com/irismod/coinswap v0.0.0-20200828021228-f86e9736690e
	github.com/irismod/htlc v0.0.0-20200830023142-e3da67f16b03
	github.com/irismod/nft v1.1.1-0.20200827095318-d16861212579
	github.com/irismod/record v1.1.1-0.20200827095301-3e27fc43ae73
	github.com/irismod/service v1.1.1-0.20200831103113-6c6a5a77d22f
	github.com/irismod/token v1.1.1-0.20200831162307-f05f6766f5b7
	github.com/irisnet/irishub v0.16.3-0.20200901032435-08bf001a1c1b
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/robfig/cron v1.2.0
	github.com/tendermint/tendermint v0.34.0-rc3
	go.uber.org/zap v1.13.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20200827085823-7b1a0adbfd1e
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
