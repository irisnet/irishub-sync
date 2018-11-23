package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/upgrade"
)

var (
	Cdc *codec.Codec
)

func init() {
	Cdc = codec.New()

	bank.RegisterCodec(Cdc)
	stake.RegisterCodec(Cdc)
	slashing.RegisterCodec(Cdc)
	auth.RegisterCodec(Cdc)
	gov.RegisterCodec(Cdc)
	upgrade.RegisterCodec(Cdc)

	sdk.RegisterCodec(Cdc)

	codec.RegisterCrypto(Cdc)
}
