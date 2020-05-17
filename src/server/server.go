package server

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	Port           int    `ini:"port"`
	PrivateKeyFile string `ini:"private-key-file"`
	PublicKeyFile  string `ini:"public-key-file"`
	StoreType      string `ini:"store"`

	PublicKey  string // base64-encoded public key used in HTTP header
	PrivateKey *ecdsa.PrivateKey
	EtcdConfig *EtcdConfig
}

type EtcdConfig struct {
	Endpoints []string `ini:"endpoints" delim:","`
	CertFile  string   `ini:"cert-file"`
	KeyFile   string   `ini:"key-file"`
	CAFile    string   `ini:"ca-file"`
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

	keyfile, e = ioutil.ReadFile(conf.PublicKeyFile)
	if e != nil {
		return conf, fmt.Errorf("Failed to read public key: %s", e)
	}

	block, _ = pem.Decode(keyfile)
	if block == nil || block.Type != "PUBLIC KEY" {
		return conf, fmt.Errorf("Failed to decode PEM block containing public key")
	}

	_, e = x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		return conf, fmt.Errorf("Failed to parse EC public key: %s", e)
	}

	// TODO: Confirm whether this is standard base64 or base64-url
	conf.PublicKey = base64.StdEncoding.EncodeToString(keyfile)
	// conf.PublicKey = base64.RawURLEncoding.EncodeToString(keyfile)

	if conf.StoreType == "etcd" {
		conf.EtcdConfig = &EtcdConfig{}

		sec, e := i.GetSection("etcd")
		if e != nil {
			return conf, fmt.Errorf("Missing etcd section", e)
		}

		if e := sec.MapTo(conf.EtcdConfig); e != nil {
			return conf, fmt.Errorf("Failed to decode etcd config: %s", e)
		}

		if len(conf.EtcdConfig.Endpoints) < 1 {
			return conf, fmt.Errorf("No etcd endpoints")
		}

		if _, e := os.Stat(conf.EtcdConfig.CertFile); e != nil {
			return conf, fmt.Errorf("Failed to read etcd cert file: %s", e)
		}

		if _, e := os.Stat(conf.EtcdConfig.KeyFile); e != nil {
			return conf, fmt.Errorf("Failed to read etcd key file: %s", e)
		}

		if _, e := os.Stat(conf.EtcdConfig.CAFile); e != nil {
			return conf, fmt.Errorf("Failed to read etcd CA file: %s", e)
		}
	}

	return conf, nil
}
