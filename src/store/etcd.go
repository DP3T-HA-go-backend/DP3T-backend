package store

import (
	"dp3t-backend/api"
	"dp3t-backend/server"

	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

const authcodesNamespace string = "/authcodes/"
const exposedNamespace string = "/exposed/"

type Etcd struct {
	ClientConfig *clientv3.Config
	Timeout      time.Duration
}

func (e *Etcd) Init(conf *server.Config) error {
	TLSInfo := transport.TLSInfo{
		CertFile:      conf.EtcdConfig.CertFile,
		KeyFile:       conf.EtcdConfig.KeyFile,
		TrustedCAFile: conf.EtcdConfig.CAFile,
	}

	TLS, err := TLSInfo.ClientConfig()
	if err != nil {
		return err
	}

	e.ClientConfig = &clientv3.Config{
		Endpoints:   []string(conf.EtcdConfig.Endpoints),
		DialTimeout: 5 * time.Second,
		TLS:         TLS,
	}

	e.Timeout = 10 * time.Second

	return nil
}

func (e *Etcd) GetExposed(timestamp int64) (*api.ProtoExposedList, error) {

	r1 := KVGetAllKeys(e.ClientConfig, exposedNamespace, e.Timeout)
	if r1 != nil {
		exposees := make([]*api.ProtoExposee, 0, len(r1.Kvs))
		for _, exposee := range r1.Kvs {
			splits := strings.Split(string(exposee.Key), "/")
			strkey := splits[len(splits)-1]
			key, _ := strconv.ParseInt(strkey, 10, 64)

			log.Printf("RAW: ExposeeKey: %s - KeyDate: %s\n", base64.StdEncoding.EncodeToString(exposee.Value), strkey)
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
	ts_ms := time.Now().UnixNano() / int64(time.Millisecond)
	expirationTTL := int64((3600 * 24 * 21) - (ts_ms - exposee.KeyDate) / 1000)
	log.Printf("Storing new Exposee: Date: %s, Key %s (expiration %ds)", strconv.FormatInt(exposee.KeyDate, 10), base64.StdEncoding.EncodeToString(exposee.Key), expirationTTL)

	r1 := KVPutAndDelete(e.ClientConfig, authcodesNamespace, exposee.AuthData.Value, exposedNamespace, string(exposee.Key), strconv.FormatInt(exposee.KeyDate, 10), expirationTTL, e.Timeout)
	if r1 != nil {
		return errors.New("Exposee could not be added")
	}
	return nil
}

func (e *Etcd) AddAuthCode(code string) error {
	r1 := KVPutIfNotExists(e.ClientConfig, authcodesNamespace, code, "", e.Timeout)
	if r1 != nil {
		return errors.New("Authcode already existed")
	}
	return nil
}

func (e *Etcd) ExpireExposees() error {
	_ = KVDeleteAllKeys(e.ClientConfig, e.Timeout)
	return nil
}

func (e *Etcd) ExpireAuthCodes() error {
	return nil
}
