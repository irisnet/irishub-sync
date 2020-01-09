module github.com/irisnet/irishub-sync

go 1.13

require (
	github.com/go-kit/kit v0.9.0
	github.com/irisnet/irishub v0.16.0-rc0
	github.com/jolestar/go-commons-pool v2.0.0+incompatible
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.7
	go.uber.org/zap v1.12.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
)

replace (
	github.com/tendermint/iavl => github.com/irisnet/iavl v0.12.3
	github.com/tendermint/tendermint => github.com/irisnet/tendermint v0.32.0
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)
