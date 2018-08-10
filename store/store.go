// init mongodb session and provide common functions

package store

import (
	"time"

	conf "github.com/irisnet/irishub-sync/conf/db"
	"github.com/irisnet/irishub-sync/module/logger"

	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session *mgo.Session
	docs    []Docs
)

func RegisterDocs(d Docs) {
	docs = append(docs, d)
}

func InitWithAuth() {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	addrs := []string{addr}

	dialInfo := &mgo.DialInfo{
		Addrs:     addrs, // []string{"192.168.6.122"}
		Database:  conf.Database,
		Username:  conf.User,
		Password:  conf.Passwd,
		Direct:    false,
		Timeout:   time.Second * 10,
		PoolLimit: 4096, // Session.SetPoolLimit
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		logger.Error.Panicln(err)
	}
}

func getSession() *mgo.Session {
	// max session num is 4096
	return session.Clone()
}

// get collection object
func ExecCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(conf.Database).C(collection)
	return s(c)
}

func Find(collection string, query interface{}) *mgo.Query {
	session := getSession()
	defer session.Close()
	c := session.DB(conf.Database).C(collection)
	return c.Find(query)
}

func Save(h Docs) error {
	save := func(c *mgo.Collection) error {
		pk := h.PkKvPair()
		n, _ := c.Find(pk).Count()
		if n >= 1 {
			errMsg := fmt.Sprintf("Record exists")
			return errors.New(errMsg)
		}
		// logger.Info.Printf("insert %s  %+v\n", h.Name(), h)
		return c.Insert(h)
	}
	return ExecCollection(h.Name(), save)
}

func SaveAll(collectionName string, docs []interface{}) error {
	session := getSession()
	defer session.Close()

	c := session.DB(conf.Database).C(collectionName)
	return c.Insert(docs...)
}

func SaveOrUpdate(h Docs) error {
	save := func(c *mgo.Collection) error {
		n, err := c.Find(h.PkKvPair()).Count()
		if err != nil {
			logger.Error.Printf("Count:%d err:%+v\n", n, err)
		}

		if n >= 1 {
			return Update(h)
		}
		// logger.Trace.Printf("insert %s  %+v\n", h.Name(), h)
		return c.Insert(h)
	}

	return ExecCollection(h.Name(), save)
}

func Update(h Docs) error {
	update := func(c *mgo.Collection) error {
		key := h.PkKvPair()
		// logger.Trace.Printf("update %s set %+v where %+v\n", h.Name(), h, key)
		return c.Update(key, h)
	}
	return ExecCollection(h.Name(), update)
}

func Delete(h Docs) error {
	remove := func(c *mgo.Collection) error {
		key := h.PkKvPair()
		return c.Remove(key)
	}
	return ExecCollection(h.Name(), remove)
}

func Query(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	return results, ExecCollection(collectionName, exop)
}
