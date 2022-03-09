package state

import (
	"io/ioutil"
	"log"

	"github.com/hashicorp/go-memdb"
)

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{},
}

type StateStore struct {
	db *memdb.MemDB
}

func NewStateStore() (*StateStore, error) {
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		return nil, err
	}

	return &StateStore{
		db: db,
	}, nil
}

func (s *StateStore) SetLogger(logger *log.Logger) {
}

var defaultLogger = log.New(ioutil.Discard, "", 0)
