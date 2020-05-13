package store

import (
	"dp3t-backend/api"
)

type Store interface {
	Init() error

	// Returns the list of exposees for a given day
	GetExposed(timestamp int64) (*api.ProtoExposedList, error)

	// Atomically add exposee to list if auth code is present
	AddExposee(exposee *api.ProtoExposee) error

	AddAuthCode(code string) error

	ExpireExposees() error
	ExpireAuthCodes() error
}
