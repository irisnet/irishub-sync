package helper

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	stake "github.com/cosmos/cosmos-sdk/x/staking/types"
	sdktypes "github.com/irisnet/irishub/types"
)

const (
	delegatorDelegationsPath          = "custom/stake/delegatorDelegations"
	delegatorUnbondingDelegationsPath = "custom/stake/delegatorUnbondingDelegations"
)

// query delegator delegations from store
func GetDelegations(delegator string) []types.Delegation {
	var (
		delegations []types.Delegation
	)
	cdc := types.GetCodec()

	addr, err := types.AccAddressFromBech32(delegator)
	if err != nil {
		logger.Error("get addr from hex failed", logger.String("address", delegator),
			logger.String("err", err.Error()))
		return nil
	}

	params := stake.NewQueryDelegatorParams(addr)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		logger.Error("get query key fail", logger.String("delegatorAddr", delegator),
			logger.String("err", err.Error()))
		return nil
	}

	res, err := QueryWithPath(bz, delegatorDelegationsPath)
	if err != nil {
		logger.Error("query tm store fail", logger.String("err", err.Error()))
		return nil
	}

	if err := cdc.UnmarshalJSON(res, &delegations); err != nil {
		logger.Error("unmarshal json fail", logger.String("err", err.Error()))
		return nil
	} else {
		return delegations
	}
}

// query delegator unbondingDelegations from store
func GetUnbondingDelegations(delegator string) []types.UnbondingDelegation {
	var (
		unbondingDelegations []types.UnbondingDelegation
	)
	cdc := types.GetCodec()

	addr, err := types.AccAddressFromBech32(delegator)
	if err != nil {
		logger.Error("get addr from hex failed", logger.String("address", delegator),
			logger.String("err", err.Error()))
		return nil
	}

	params := stake.NewQueryDelegatorParams(addr)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		logger.Error("get query key fail", logger.String("delegatorAddr", delegator),
			logger.String("err", err.Error()))
		return nil
	}

	res, err := QueryWithPath(bz, delegatorUnbondingDelegationsPath)
	if err != nil {
		logger.Error("query tm store fail", logger.String("err", err.Error()))
		return nil
	}

	if err := cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
		logger.Error("unmarshal json fail", logger.String("err", err.Error()))
		return nil
	} else {
		return unbondingDelegations
	}
}

func CalculateDelegatorDelegationTokens(delegations []types.Delegation) float64 {
	var (
		token types.Dec
	)
	token = sdktypes.ZeroDec()
	if len(delegations) > 0 {
		for _, v := range delegations {
			validatorAddr := v.ValidatorAddress.String()
			if validator, err := GetValidator(validatorAddr); err != nil {
				logger.Error("get validator fail", logger.String("validatorAddr", validatorAddr),
					logger.String("err", err.Error()))
				continue
			} else {
				token = token.Add(validator.DelegatorShareExRate().Mul(v.Shares))
			}
		}
	}

	return ParseFloat(token.String())
}

func CalculateDelegatorUnbondingDelegationTokens(unbondingDelegations []types.UnbondingDelegation) float64 {
	var (
		token types.Int
	)
	token = sdktypes.ZeroInt()

	if len(unbondingDelegations) > 0 {
		for _, v := range unbondingDelegations {
			token = token.Add(v.InitialBalance.Amount)
		}
	}

	return ParseFloat(token.String())
}
