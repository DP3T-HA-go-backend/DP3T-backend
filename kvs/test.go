package main

import (
	"context"
	"log"
	//"os"
	"time"
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"

	//"google.golang.org/grpc/grpclog"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"10.0.26.10:2379", "10.0.26.11:2379", "10.0.26.13:2379"}
)

func main() {
	ExampleConfig_withTLS()
	ExampleKV_get()
}


func ExampleConfig_withTLS() {
	tlsInfo := transport.TLSInfo{
		CertFile:      "/etc/ssl/etcd/ssl/node-node1.pem",
		KeyFile:       "/etc/ssl/etcd/ssl/node-node1-key.pem",
		TrustedCAFile: "/etc/ssl/etcd/ssl/ca.pem",
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, "foo", "bar")
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
}

func ExampleKV_get() {
	tlsInfo := transport.TLSInfo{
		CertFile:      "/etc/ssl/etcd/ssl/node-node1.pem",
		KeyFile:       "/etc/ssl/etcd/ssl/node-node1-key.pem",
		TrustedCAFile: "/etc/ssl/etcd/ssl/ca.pem",
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	_, err = cli.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// Output: foo : bar
}

