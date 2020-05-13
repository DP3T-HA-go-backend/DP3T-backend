package store

import (
	"dp3t-backend/api"
	kvs "dp3t-backend/etcd"
	"errors"

	"crypto/tls"
	"time"

	"go.etcd.io/etcd/pkg/transport"
)

type Etcd struct {
	Endpoints      []string
	TLSInfo        transport.TLSInfo
	DialTimeout    time.Duration
	RequestTimeout time.Duration

	TLS *tls.Config
}

func (e *Etcd) Init() error {
	e.DialTimeout = 5 * time.Second
	e.RequestTimeout = 10 * time.Second
	e.Endpoints = []string{"10.0.26.10:2379", "10.0.26.11:2379", "10.0.26.13:2379"}
	e.TLSInfo = transport.TLSInfo{
		CertFile:      "/etc/ssl/etcd/ssl/node-node1.pem",
		KeyFile:       "/etc/ssl/etcd/ssl/node-node1-key.pem",
		TrustedCAFile: "/etc/ssl/etcd/ssl/ca.pem",
	}

	var err error
	e.TLS, err = e.TLSInfo.ClientConfig()
	if err != nil {
		return err
	}

	return nil
}

func (e *Etcd) GetExposed(timestamp int64) (*api.ProtoExposedList, error) {
	return nil, nil
}

func (e *Etcd) AddExposee(exposee *api.ProtoExposee) error {
	return nil
}

func (e *Etcd) AddAuthCode(code string) error {
	r1 := kvs.KVPutIfNotExists(code, "")
	if r1 != nil {
		return errors.New("Authcode already existed")
	}
	return nil
}

func (e *Etcd) ExpireExposees() error {
	return nil
}

func (e *Etcd) ExpireAuthCodes() error {
	return nil
}
