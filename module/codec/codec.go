package codec

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/modules/upgrade"
)

var (
	Cdc *wire.Codec
)

func init() {
	Cdc = wire.NewCodec()

	ibc.RegisterWire(Cdc)
	bank.RegisterWire(Cdc)
	stake.RegisterWire(Cdc)
	slashing.RegisterWire(Cdc)
	auth.RegisterWire(Cdc)
	gov.RegisterWire(Cdc)
	upgrade.RegisterWire(Cdc)

	sdktypes.RegisterWire(Cdc)

	wire.RegisterCrypto(Cdc)
}
