package cdc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/irisnet/irishub/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//token "github.com/irismod/token/types"
	//auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	//bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	//distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	//dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	//gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	//rand "github.com/irisnet/irishub/modules/random/types"
	//oracle "github.com/irisnet/irishub/modules/oracle/types"
	//slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	//stake "github.com/cosmos/cosmos-sdk/x/staking/types"
	//staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	//crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	//coinswap "github.com/irismod/coinswap/types"
	//htlc "github.com/irismod/htlc/types"
	//nft "github.com/irismod/nft/types"
	//record "github.com/irismod/record/types"
	//service "github.com/irismod/service/types"
	//guardian "github.com/irisnet/irishub/modules/guardian/types"
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
)

var (
	cdc *codec.LegacyAmino

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
	)
)
// 初始化账户地址前缀
func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(address.Bech32PrefixAccAddr, address.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()
	cdc = codec.New()
	moduleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterEvidences(cdc)
}

func GetCodec() *codec.LegacyAmino {
	return cdc
}
