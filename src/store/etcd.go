package store

import (
	"dp3t-backend/api"
	kvs "dp3t-backend/etcd"
	"dp3t-backend/server"
	"encoding/base64"
	"log"
	"math"
	"strconv"

	"crypto/tls"
	"errors"
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

	ClientConfig *clientv3.Config
}

func (e *Etcd) Init(conf *server.Config) error {
	e.DialTimeout = 5 * time.Second
	e.RequestTimeout = 10 * time.Second
	e.Endpoints = []string(conf.EtcdConfig.Endpoints)
	e.TLSInfo = transport.TLSInfo{
		CertFile:      conf.EtcdConfig.CertFile,
		KeyFile:       conf.EtcdConfig.KeyFile,
		TrustedCAFile: conf.EtcdConfig.CAFile,
	}

	var err error
	e.TLS, err = e.TLSInfo.ClientConfig()
	if err != nil {
		return err
	}

	e.ClientConfig = &clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: e.DialTimeout,
		TLS:         e.TLS,
	}

	return nil
}

func (e *Etcd) GetExposed(timestamp int64) (*api.ProtoExposedList, error) {

	r1 := kvs.KVGetAllKeys(e.ClientConfig, strconv.FormatInt(timestamp, 10), e.RequestTimeout)
	if r1 != nil {
		exposees := make([]*api.ProtoExposee, 0, len(r1.Kvs))
		for _, exposee := range r1.Kvs {
			key, _ := strconv.ParseInt(string(exposee.Key), 10, 64)

			log.Printf("RAW: ExposeeKey: %s - KeyDate: %s\n", base64.StdEncoding.EncodeToString(exposee.Value), string(exposee.Key))
			exposees = append(exposees, &api.ProtoExposee{
				Key:     exposee.Value,
				KeyDate: key,
			})
		}
		data := &api.ProtoExposedList{
			BatchReleaseTime: timestamp,
			Exposed:          exposees,
		}

		return data, nil
	}
	return nil, errors.New("Could not retrieve exposees")

}

func (e *Etcd) AddExposee(exposee *api.ProtoExposee) error {
	expirationTTL := int64(math.Ceil(float64(((time.Now().UnixNano() / int64(time.Millisecond)) - exposee.KeyDate) / 1000 * 3660 * 24)))
	log.Printf("Storing new Exposee: Date: %s, Key %s", strconv.FormatInt(exposee.KeyDate, 10), base64.StdEncoding.EncodeToString(exposee.Key))
	r1 := kvs.KVPutAndDelete(e.ClientConfig, exposee.AuthData.Value, strconv.FormatInt(exposee.KeyDate, 10), string(exposee.Key), expirationTTL, e.RequestTimeout)
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
	_ = kvs.KVDeleteAllKeys(e.ClientConfig, e.RequestTimeout)
	return nil
}

func (e *Etcd) ExpireAuthCodes() error {
	return nil
}
