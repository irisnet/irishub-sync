package cron

import (
	"time"
	"os"
	"os/signal"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/util/helper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CronService struct{}

func (s *CronService) StartCronService() {

	logger.Info("Start Update Txs CronService ...")
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	fnUpdate := func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("CronService have error", logger.Any("err", r))
			}
		}()

		runValue := true
		skip := 0
		for runValue {
			total, err := UpdateUnknownOrEmptyTypeTxsByPage(skip, 20)
			if err != nil {
				logger.Error("getCommonTx have error", logger.String("err", err.Error()))
			}
			if total < 20 {
				runValue = false
				logger.Info("Finish UpdateUnknownTxsByPage.", logger.Int("total", total))
			} else {
				skip = skip + total
				logger.Info("Continue UpdateUnknownTxsByPage", logger.Int("skip", skip))
			}
		}

		logger.Info("Finish Update Txs.")
	}
	fnUpdate()
	for {
		select {
		case <-ticker.C:
			fnUpdate()
		case <-stop:
			close(stop)
			logger.Info("Update Txs CronService Quit...")
			return
		}

	}

}

func UpdateUnknownOrEmptyTypeTxsByPage(skip, limit int) (int, error) {

	res, err := new(document.CommonTx).GetUnknownOrEmptyTypeTxs(skip, limit)
	if err != nil {
		return 0, err
	}

	if len(res) > 0 {
		doWork(res, UpdateUnknowOrEmptyTypeTxs)
	}

	return len(res), nil
}

func doWork(commonTxs []document.CommonTx, fn func([]*document.CommonTx) error) {
	client := helper.GetClient()
	defer func() {
		client.Release()
	}()

	for _, val := range commonTxs {
		txs, err := ParseUnknownOrEmptyTypeTxs(val.Height, client)
		if err != nil {
			logger.Error("ParseUnknownTxs have error", logger.String("error", err.Error()))
			continue
		}
		if err := fn(txs); err != nil {
			logger.Warn("UpdateUnknowTxs have error", logger.String("error", err.Error()))
		}
	}

}

func ParseUnknownOrEmptyTypeTxs(b int64, client *helper.Client) (commontx []*document.CommonTx, err error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse block fail", logger.Int64("blockHeight", b),
				logger.Any("err", err))
		}
	}()

	block, err := client.Block(&b)
	if err != nil {
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := helper.GetClient()
		block, err2 = client2.Block(&b)
		client2.Release()
		if err2 != nil {
			return nil, err2
		}
	}

	commontx = make([]*document.CommonTx, 0, len(block.Block.Txs))

	for _, txByte := range block.Block.Txs {
		tx := helper.ParseTx(txByte, block.Block)
		if tx.Status != document.Unknow_Status {
			commontx = append(commontx, &tx)
		}

	}
	return
}

func UpdateUnknowOrEmptyTypeTxs(commontx []*document.CommonTx) error {

	update_fn := func(tx *document.CommonTx) error {
		fn := func(c *mgo.Collection) error {
			return c.Update(bson.M{"tx_hash": tx.TxHash},
				bson.M{"$set": bson.M{"from": tx.From, "to": tx.To, "type": tx.Type, "amount": tx.Amount,
					"actual_fee": tx.ActualFee, "status": tx.Status, "tags": tx.Tags, "msgs": tx.Msgs,
					"code": tx.Code, "log": tx.Log, "gas_wanted": tx.GasWanted}})
		}

		if err := store.ExecCollection(document.CollectionNmCommonTx, fn); err != nil {
			return err
		}
		return nil
	}

	for _, dbval := range commontx {
		update_fn(dbval)
	}

	return nil
}