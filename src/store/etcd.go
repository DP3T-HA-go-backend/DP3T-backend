package store

import (
    "fmt"
    "dp3t-backend/api"
	"dp3t-backend/server"
    kvs "dp3t-backend/etcd"

    "errors"
    "crypto/tls"
    "time"

    "go.etcd.io/etcd/clientv3"
    "go.etcd.io/etcd/pkg/transport"
)

type Etcd struct {
    Endpoints      []string
    TLSInfo        transport.TLSInfo
    DialTimeout    time.Duration
    RequestTimeout time.Duration

    TLS *tls.Config

    ClientConfig clientv3.Config
}

func (e *Etcd) Init(conf *server.Config) error {
    fmt.Println("Inicializamos todo")
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

    e.ClientConfig = clientv3.Config{
        Endpoints:   e.Endpoints,
        DialTimeout: e.DialTimeout,
        TLS:         e.TLS,
    }

    return nil
}

func (e *Etcd) GetExposed(timestamp int64) (*api.ProtoExposedList, error) {
    return nil, nil
}

func (e *Etcd) AddExposee(exposee *api.ProtoExposee) error {
    //KVPutAndDelete(KeyToDelete string, KeyToPut string, ValueToPut string)
    r1 := kvs.KVPutAndDelete(e.ClientConfig, exposee.AuthData.Value, string(exposee.Key), string(exposee.KeyDate), e.RequestTimeout)
    if r1 != nil {
        return errors.New("Exposee could not be added")
    }
    return nil
}

func (e *Etcd) AddAuthCode(code string) error {
    r1 := kvs.KVPutIfNotExists(e.ClientConfig, code, "", e.RequestTimeout)
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
