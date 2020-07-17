module github.com/irisnet/irishub-sync

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200713224032-8a62e1ab8127
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/jolestar/go-commons-pool v2.0.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/robfig/cron v1.2.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tendermint/tendermint v0.33.6
	go.uber.org/zap v1.13.0 // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
