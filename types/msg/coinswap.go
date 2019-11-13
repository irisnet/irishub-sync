package msg

import (
	"github.com/irisnet/irishub-sync/util/constant"
	"github.com/irisnet/irishub-sync/store"
	itypes "github.com/irisnet/irishub-sync/types"
)

type DocTxMsgAddLiquidity struct {
	MaxToken     store.Coin `bson:"max_token"`      // coin to be deposited as liquidity with an upper bound for its amount
	ExactIrisAmt string     `bson:"exact_iris_amt"` // exact amount of native asset being add to the liquidity pool
	MinLiquidity string     `bson:"min_liquidity"`  // lower bound UNI sender is willing to accept for deposited coins
	Deadline     int64      `bson:"deadline"`
	Sender       string     `bson:"sender"`
}

func (doctx *DocTxMsgAddLiquidity) Type() string {
	return constant.TxTypeAddLiquidity
}

func (doctx *DocTxMsgAddLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgAddLiquidity)
	doctx.Sender = msg.Sender.String()
	doctx.MinLiquidity = msg.MinLiquidity.String()
	doctx.ExactIrisAmt = msg.ExactIrisAmt.String()
	doctx.Deadline = msg.Deadline
	doctx.MaxToken = itypes.ParseCoin(msg.MaxToken.String())
}

type DocTxMsgRemoveLiquidity struct {
	MinToken          string     `bson:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity store.Coin `bson:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinIrisAmt        string     `bson:"min_iris_amt"`       // minimum amount of the native asset the sender is willing to accept
	Deadline          int64      `bson:"deadline"`
	Sender            string     `bson:"sender"`
}

func (doctx *DocTxMsgRemoveLiquidity) Type() string {
	return constant.TxTypeRemoveLiquidity
}

func (doctx *DocTxMsgRemoveLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgRemoveLiquidity)
	doctx.Sender = msg.Sender.String()
	doctx.MinIrisAmt = msg.MinIrisAmt.String()
	doctx.MinToken = msg.MinToken.String()
	doctx.Deadline = msg.Deadline
	doctx.WithdrawLiquidity = itypes.ParseCoin(msg.WithdrawLiquidity.String())
}

type DocTxMsgSwapOrder struct {
	Input      Input  `bson:"input"`        // the amount the sender is trading
	Output     Output `bson:"output"`       // the amount the sender is receiving
	Deadline   int64  `bson:"deadline"`     // deadline for the transaction to still be considered valid
	IsBuyOrder bool   `bson:"is_buy_order"` // boolean indicating whether the order should be treated as a buy or sell
}

type Input struct {
	Address string     `bson:"address"`
	Coin    store.Coin `bson:"coin"`
}

type Output struct {
	Address string     `bson:"address"`
	Coin    store.Coin `bson:"coin"`
}

func (doctx *DocTxMsgSwapOrder) Type() string {
	return constant.TxTypeSwapOrder
}

func (doctx *DocTxMsgSwapOrder) BuildMsg(txMsg interface{}) {
	msg := txMsg.(itypes.MsgSwapOrder)
	doctx.Deadline = msg.Deadline
	doctx.IsBuyOrder = msg.IsBuyOrder
	doctx.Input = Input{
		Address: msg.Input.Address.String(),
		Coin:    itypes.ParseCoin(msg.Input.Coin.String()),
	}
	doctx.Output = Output{
		Address: msg.Output.Address.String(),
		Coin:    itypes.ParseCoin(msg.Output.Coin.String()),
	}
}
