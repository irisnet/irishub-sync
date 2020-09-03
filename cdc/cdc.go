package cdc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/irisnet/irishub/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/oracle"
	"github.com/irisnet/irishub/modules/random"
	"github.com/irismod/record"
	"github.com/irismod/service"
	"github.com/irismod/token"
	"github.com/irismod/nft"
	"github.com/irismod/coinswap"
	"github.com/irismod/htlc"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/std"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/irisnet/irishub-sync/types"
)

var (
	encodecfg    params.EncodingConfig
	moduleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.AppModuleBasic{},
		slashing.AppModuleBasic{},
		staking.AppModuleBasic{},
		evidence.AppModuleBasic{},
		crisis.AppModuleBasic{},
		random.AppModuleBasic{},
		nft.AppModuleBasic{},
		token.AppModuleBasic{},
		service.AppModuleBasic{},
		record.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		guardian.AppModuleBasic{},
		oracle.AppModuleBasic{},
		htlc.AppModuleBasic{},
	)
)
// 初始化账户地址前缀
func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(address.Bech32PrefixAccAddr, address.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()
	types.Bech32AccountAddrPrefix = sdk.GetConfig().GetBech32AccountAddrPrefix()

	amino := codec.New()
	interfaceRegistry := ctypes.NewInterfaceRegistry()
	moduleBasics.RegisterInterfaces(interfaceRegistry)
	sdk.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, std.DefaultPublicKeyCodec{}, tx.DefaultSignModes)

	encodecfg = params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

func GetTxDecoder() sdk.TxDecoder {
	return encodecfg.TxConfig.TxDecoder()
}
