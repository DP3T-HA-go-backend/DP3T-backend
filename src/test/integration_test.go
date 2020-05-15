package test

import (
	"dp3t-backend/api"

	"context"
	"log"
	"testing"
	"encoding/json"
	"net/http"
	"crypto/tls"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

const exposed_url string = "http://localhost:8080"
const authcode_url string = "http://localhost:8081"

var endpoints = []string{"0.0.0.0:2379"}
var etcd_conf *clientv3.Config

func init() {
	var err error
	var tls_conf *tls.Config

	tls_info := transport.TLSInfo{
		CertFile:      "config/etcd/server.pem",
		KeyFile:       "config/etcd/server-key.pem",
		TrustedCAFile: "config/etcd/ca.pem",
	}

	tls_conf, err = tls_info.ClientConfig()
	if err != nil {
		log.Fatal("Error initializing etcd's TLS configuration", err)
	}

	etcd_conf = &clientv3.Config{
		Endpoints:   endpoints,
		TLS:         tls_conf,
	}
}

func TestAuthcode(t *testing.T) {
	cli, err := clientv3.New(*etcd_conf)
	e0, err := cli.Get(context.TODO(), "/authcodes", clientv3.WithPrefix())
	num_codes_0 := len(e0.Kvs)

	r, err := http.Get(authcode_url)
	if err != nil {
		t.Errorf("Failed to fetch %s", authcode_url)
		return
	}

	code := &api.ProtoAuthData{}
	err = json.NewDecoder(r.Body).Decode(code)
	if err != nil {
		t.Errorf("Failed to decode JSON %s", r.Body)
		return
	}

	t.Log("Generate new authcode:", code.Value)

	e1, err := cli.Get(context.TODO(), "/authcodes", clientv3.WithPrefix())
	num_codes_1 := len(e1.Kvs)

	if num_codes_0 + 1 != num_codes_1 {
		t.Errorf("Mismatched number of authcodes")
		return
	}

	t.Log("Matching number of authcodes")

	e2, err := cli.Get(context.TODO(), "/authcodes/" + code.Value)
	if len(e2.Kvs) == 0 {
		t.Errorf("Key %s not present", code.Value)
		return
	}

	t.Log("Key", code.Value, "exists in etcd")
}
