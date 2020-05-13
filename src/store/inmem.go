package store

import (
	"dp3t-backend/api"
)

type InMem struct {
	exposed *api.ProtoExposedList
	codes map[string]codeValue
}

type codeValue struct {
	time uint64
	used bool
}

func (m *InMem) Init() error {
	m.exposed = &api.ProtoExposedList{
		BatchReleaseTime: 123456789,
		Exposed:          []*api.ProtoExposee{},
	}
	m.codes = make(map[string]codeValue)
	return nil
}

func (m *InMem) GetExposed(timestamp int64) (*api.ProtoExposedList, error) {
	return m.exposed, nil
}

func (m *InMem) AddExposee(exposee *api.ProtoExposee) error {
	m.exposed.Exposed = append(m.exposed.Exposed, exposee)
	return nil
}

func (m *InMem) AddAuthCode(code string) error {
	m.codes[code] = codeValue{0, false}
	return nil
}

func (m *InMem) ExpireExposees() error {
	return nil
}

func (m *InMem) ExpireAuthCodes() error {
	return nil
}
