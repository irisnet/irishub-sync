module github.com/irisnet/irishub-sync

        go 1.13

        require (
        github.com/cosmos/cosmos-sdk v0.34.4-0.20200716143106-10783f27d0cf
        github.com/go-kit/kit v0.10.0
        github.com/gogo/protobuf v1.3.1 // indirect
        github.com/irismod/coinswap v0.0.0-20200717084559-d162bdf94677
        github.com/irismod/htlc v0.0.0-20200717084245-6c78f425eb6b
        github.com/irismod/nft v1.1.1-0.20200717090223-19ae27993d05
        github.com/irismod/service v1.1.1-0.20200717083211-da28297d9e73
        github.com/irismod/token v1.1.1-0.20200717083658-a6d44d130830
        github.com/irisnet/irishub v0.16.3-0.20200720082623-67145905e2aa
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
        github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.19.1-0.20200720031859-08f888c08df5
        github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
        )
