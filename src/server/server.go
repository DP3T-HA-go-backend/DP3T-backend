package server

import (
	"dp3t-backend/store"

	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/ini.v1"
)

const PUBLIC_KEY string = "" +
	"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlL" +
	"b1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhy" +
	"Sjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90" +
	"clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"

type Config struct {
	Port           int    `ini:"port"`
	PrivateKeyFile string `ini:"private-key-file"`
	StoreType      string `ini:"store"`
	PrivateKey     *ecdsa.PrivateKey
}

func InitConfig(conf_file string) (*Config, error) {
	conf := &Config{}

	i, e := ini.Load(conf_file)
	if e != nil {
		return conf, fmt.Errorf("Failed to read config file: %s", e)
	}

	if e := i.MapTo(conf); e != nil {
		return conf, fmt.Errorf("Failed to decode config: %s", e)
	}

	if _, e := os.Stat(conf.PrivateKeyFile); e != nil {
		return conf, fmt.Errorf("Failed to read private key: %s", e)
	}

	keyfile, e := ioutil.ReadFile(conf.PrivateKeyFile)
	if e != nil {
		return conf, fmt.Errorf("Failed to read private key: %s", e)
	}

	block, _ := pem.Decode(keyfile)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return conf, fmt.Errorf("Failed to decode PEM block containing EC private key")
	}

	conf.PrivateKey, e = x509.ParseECPrivateKey(block.Bytes)
	if e != nil {
		return conf, fmt.Errorf("Failed to parse EC private key: %s", e)
	}

	return conf, nil
}

func InitStore(conf *Config) (store.Store, error) {
	switch conf.StoreType {
	case "inmem":
		return &store.InMem{}, nil
	case "etcd":
		return &store.Etcd{}, nil
	default:
		return nil, fmt.Errorf("Unknown store kind %s", conf.StoreType)
	}
}
