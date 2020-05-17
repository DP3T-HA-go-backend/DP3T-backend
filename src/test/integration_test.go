package test

import (
	"dp3t-backend/api"

	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const exposed_url string = "http://localhost:8080"
const authcode_url string = "http://localhost:8081"

var endpoints = []string{"0.0.0.0:2379"}
var etcd_conf *clientv3.Config

func init() {
	var err error
	var tls_conf *tls.Config

	tls_info := transport.TLSInfo{
		CertFile:      "../config/test/etcd/server.crt",
		KeyFile:       "../config/test/etcd/server.key",
		TrustedCAFile: "../config/test/etcd/ca.crt",
	}

	tls_conf, err = tls_info.ClientConfig()
	if err != nil {
		log.Fatal("Error initializing etcd's TLS configuration", err)
	}

	etcd_conf = &clientv3.Config{
		Endpoints: endpoints,
		TLS:       tls_conf,
	}
}

func TestAuthcode(t *testing.T) {
	cli, err := clientv3.New(*etcd_conf)
	e0, err := cli.Get(context.TODO(), "/authcodes", clientv3.WithPrefix())
	num_codes_0 := len(e0.Kvs)

	code, err := getAuthCode(t)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

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

	t.Log("Authcode", code.Value, "exists in etcd")
}

func TestExposed(t *testing.T) {
	cli, err := clientv3.New(*etcd_conf)

	e0, err := cli.Get(context.TODO(), "/exposed", clientv3.WithPrefix())
	num_exposed_0 := len(e0.Kvs)
	t.Log("Number of exposees:", num_exposed_0)

	code, err := getAuthCode(t)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	ts_ms := time.Now().UnixNano() / int64(time.Millisecond)

	exposee := &api.ProtoExposee{
		Key:      []byte("exposee0"),
		KeyDate:  ts_ms,
		AuthData: code,
	}

	t.Logf("Creating exposee (%s)", exposee)

	exposee_json, err := protojson.Marshal(exposee)
	if err != nil {
		t.Log("Failed to encode JSON", exposee)
	}

	exposee_buf := bytes.NewBuffer(exposee_json)
	resp, err := http.Post(exposed_url, "application/json", exposee_buf)
	if err != nil {
		t.Errorf("Failed to create request %s", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Failed to add exposee: %s", resp.Status)
	}

	t.Logf("Created exposee (%s)", exposee.Key)

	e1, err := cli.Get(context.TODO(), "/exposed", clientv3.WithPrefix())
	num_exposed_1 := len(e1.Kvs)
	t.Log("Number of exposees:", num_exposed_1)

	// Clean-up used keys
	cli.Delete(context.TODO(), "/exposed/"+string(exposee.Key))

	// Try to retrieve the last batch
	batch_len := int64(7200 * 1000)
	last_batch := int64(ts_ms / batch_len) * batch_len
	exposed_batch := fmt.Sprintf("%s/%d", exposed_url, last_batch)

	t.Log("Retrieving last batch:", exposed_batch)

	r, err := http.Get(exposed_batch)
	if err != nil {
		t.Log("Failed GET request", exposed_batch)
	}

	in, _ := ioutil.ReadAll(r.Body)
	exposed := &api.ProtoExposedList{}
	err = proto.Unmarshal(in, exposed)
	if err != nil {
		t.Errorf("Failed to decode protobuf %s", r.Body)
		return
	}

	num_exposees := len(exposed.Exposed)
	if num_exposees > 0 {
		t.Errorf("Unexpected number of exposees in batch: %d (0 expected)", num_exposees)
		return
	}

	t.Logf("Matching number of exposees in batch (%d)", num_exposees)
}

func getAuthCode(t *testing.T) (*api.ProtoAuthData, error) {
	r, err := http.Get(authcode_url)
	if err != nil {
		return nil, fmt.Errorf("Failed GET request %s", authcode_url)
	}

	code := &api.ProtoAuthData{}
	err = json.NewDecoder(r.Body).Decode(code)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode JSON %s", r.Body)
	}

	t.Log("Generated new authcode:", code.Value)
	return code, nil
}
