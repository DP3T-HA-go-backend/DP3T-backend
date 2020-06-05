package store

import (
	"dp3t-backend/api"
	"dp3t-backend/server"

	"fmt"
)

type Store interface {
	Init(conf *server.Config) error

	// Returns the list of exposees for a given day
	GetExposed(timestamp int64) (*api.ProtoExposedList, error)

	// Atomically add exposee to list if auth code is present
	AddExposee(exposee *api.ProtoExposee) error

	AddAuthCode(code string) error

	ExpireExposees() error
	ExpireAuthCodes() error
}

func InitStore(conf *server.Config) (Store, error) {
	switch conf.StoreType {
	case "inmem":
		return &InMem{}, nil
	case "etcd":
		return &Etcd{}, nil
	default:
		return nil, fmt.Errorf("Unknown store kind %s", conf.StoreType)
	}
}
