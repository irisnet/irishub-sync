// init mongodb session and provide common functions

package store

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"testing"
	"time"
)

func TestInitWithAuth(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "tets initWithAuth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start()
		})
	}
}

type I interface {
	Test()
}

type Block struct {
	Height int64                  `bson:"height"`
	Hash   string                 `bson:"hash"`
	Time   time.Time              `bson:"time"`
	NumTxs int64                  `bson:"num_txs"`
	I      map[string]interface{} `bson:"i"`
}

type I1 struct {
	F1 string
	F2 string
}

func (I1) Test() {
}

func TestInterface(t *testing.T) {
	Start()
	c := session.Copy().DB("").C("user1")
	var r []Block
	c.Find(bson.M{"i.f1": "f1"}).All(&r)
	fmt.Println(r)
}

func TestTransaction(t *testing.T) {
	Start()
	c := session.Copy().DB("").C("transaction")
	//txn.SetLogger(logger.GetLogger())
	//txn.SetDebug(true)
	runner := txn.NewRunner(c)
	ops := []txn.Op{
		{
			C:  "user1",
			Id: bson.NewObjectId(),
			//Insert: Block{Height: 1, I: I1{F1: "f1"}},
		},
		{
			C:      "user1",
			Id:     bson.NewObjectId(),
			Insert: Block{Height: 2, Hash: "xxxxxx"},
		},
		{
			C:      "user1",
			Id:     bson.NewObjectId(),
			Insert: Block{Height: 2},
		},
		{
			C:      "user1",
			Id:     bson.NewObjectId(),
			Insert: Block{Height: 3},
		},
		//{
		//	C:      "user2",
		//	Id:     "5bfcbb90914120370830b8ad",
		//	Assert: bson.M{"balance": bson.M{"$gte": 100}},
		//	Update: bson.M{"$set": bson.M{"height": 4}},
		//},
	}
	id := bson.NewObjectId() // Optional
	err := runner.Run(ops, id, nil)
	fmt.Println(err)
}

type Task struct {
	From       int64  `bson:"from"`
	To         int64  `bson:"to"`
	Current    int64  `bson:"current"`
	Status     string `bson:"status"`
	WorkLocker string `bson:"work_locker"`
	Role       string `bson:"role"`
}

func TestApply(t *testing.T) {
	prepare()
	for i := 1; i <= 10; i++ {
		go func(i int) {
			worker := fmt.Sprintf("thread[%d]", i)
			task, err := fetchTask(worker)
			if err != nil {
				fmt.Println("worker:", worker, " do not take task:", err.Error())
			} else {
				fmt.Println("worker:", worker, " take task:", task)
			}
		}(i)
	}
	time.Sleep(5 * time.Second)

}

func fetchTask(worker string) (Task, error) {
	var task Task
	userDao := session.Copy().DB("").C("task")
	for j := 1; j <= 5; j++ {
		change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"status": "processing", "work_locker": worker}},
			ReturnNew: true,
		}

		_, err := userDao.Find(bson.M{"status": "ready", "work_locker": "", "role": "catch-up"}).Limit(1).Apply(change, &task)
		if err == nil {
			return task, nil
		}
	}
	return task, errors.New("fail")
}

func prepare() {
	Start()

	var current = 100000
	var syncHeight = 0

	var goroutine_num = 5

	syncBlockNumFastSync := (current - syncHeight) / goroutine_num

	var batch []txn.Op
	for i := 1; i <= goroutine_num; i++ {
		var (
			start = syncHeight + (i-1)*syncBlockNumFastSync + 1
			end   = syncHeight + i*syncBlockNumFastSync
		)

		task := Task{
			From:    int64(start),
			To:      int64(end),
			Current: -1,
			Status:  "ready",
			Role:    "catch-up",
		}

		batch = append(batch, txn.Op{
			C:  "task",
			Id: bson.NewObjectId(),
			//Assert: txn.DocExists,
			Insert: task,
		})
	}

	//添加一条follow记录
	batch = append(batch, txn.Op{
		C:  "task",
		Id: bson.NewObjectId(),
		//Assert: txn.DocExists,
		Insert: Task{
			From:    int64(current + 1),
			Current: int64(current),
			Status:  "ready",
			Role:    "follow",
		},
	})

	c := session.Copy().DB("").C("transaction")
	runner := txn.NewRunner(c)

	id := bson.NewObjectId() // Optional
	runner.Run(batch, id, nil)
}

//func get_external() string {
//	resp, err := http.Get("http://myexternalip.com/raw")
//	if err != nil {
//		return ""
//	}
//	defer resp.Body.Close()
//	content, _ := ioutil.ReadAll(resp.Body)
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(resp.Body) //s := buf.String()
//	return string(content)
//}
