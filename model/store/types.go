// interface for a document

package store

import (
	"gopkg.in/mgo.v2"
)

type Docs interface {
	// collection name
	Name() string
	// primary key pair(used to find a unique record)
	PkKvPair() map[string]interface{}
	// add index for collection
	Index() []mgo.Index
}
