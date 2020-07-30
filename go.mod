module github.com/irisnet/irishub-sync

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200721190130-5d71020270ae
	github.com/go-kit/kit v0.10.0
	github.com/irismod/coinswap v0.0.0-20200722055706-deeded9d99b8
	github.com/irismod/htlc v0.0.0-20200722060015-b71f49c9b167
	github.com/irismod/nft v1.1.1-0.20200722060344-38fec5db63a2
	github.com/irismod/service v1.1.1-0.20200723031529-6abecb02ceb1
	github.com/irismod/token v1.1.1-0.20200723031618-028bdd6fb30a
	github.com/irisnet/irishub v0.16.3-0.20200723084819-68aaadaefc0d
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/robfig/cron v1.2.0
	github.com/tendermint/tendermint v0.33.6
	go.uber.org/zap v1.13.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.19.1-0.20200722022502-e2d6c76ae750
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
