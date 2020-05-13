package kvs

import (
	"context"
	"crypto/tls"
	"errors"
	"log"

	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"go.etcd.io/etcd/pkg/transport"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"10.0.26.10:2379", "10.0.26.11:2379", "10.0.26.13:2379"}
)

func tlsConfig() *tls.Config {
	tlsInfo := transport.TLSInfo{
		CertFile:      "/Users/dcarrera/Desktop/DP3T/DP3T-backend/kvs/keys/node-node1.pem",
		KeyFile:       "/Users/dcarrera/Desktop/DP3T/DP3T-backend/kvs/keys/node-node1-key.pem",
		TrustedCAFile: "/Users/dcarrera/Desktop/DP3T/DP3T-backend/kvs/keys/ca.pem",
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	return tlsConfig
}

//KVPut to put key and value
func KVPut(key string, value string) {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	_, err = cli.Put(context.TODO(), key, value)
	if err != nil {
		log.Fatal(err)
	}
}

//KVPutTTL to put key and value with a Time To Live in days
func KVPutTTL(key string, value string, days int64) {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// minimum lease TTL is in seconds
	resp, err := cli.Grant(context.TODO(), days*24*3600)
	if err != nil {
		log.Fatal(err)
	}

	// after grant seconds, the key will be removed
	_, err = cli.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
}

// KVGet to Get a key
func KVGet(key string) *clientv3.GetResponse {
	tlsConfig := tlsConfig()
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
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

//KVDelete to delete a key
func KVDelete(key string) {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// delete the keys
	_, err = cli.Delete(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
}

//KVPut only if the key did not exist
func KVPutIfNotExists(KeyToPut string, ValueToPut string) error {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	NotExistsKeyToPut := clientv3.Compare(clientv3.CreateRevision(KeyToPut), "=", 0)
	r, err := cli.Txn(ctx).If(NotExistsKeyToPut).Then(clientv3.OpPut(KeyToPut, ValueToPut)).Commit()

	if r.Succeeded {
		return nil
	}

	return errors.New("Key already existed")

}

//KVDelete only if the key did existed
func KVDeleteIfExists(KeyToDelete string) error {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	ExistsKeyToDelete := clientv3.Compare(clientv3.CreateRevision(KeyToDelete), ">", 0)
	r, err := cli.Txn(ctx).If(ExistsKeyToDelete).Then(clientv3.OpDelete(KeyToDelete)).Commit()

	if r.Succeeded {
		return nil
	}

	return errors.New("Key already existed")

}

//KVDelete one existing Key and KVPut another one only if the first existed and was deleted
func KVPutAndDelete(KeyToDelete string, KeyToPut string, ValueToPut string) error {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	NotExistsKeyToPut := clientv3.Compare(clientv3.CreateRevision(KeyToPut), "=", 0)
	ExistsKeyToDelete := clientv3.Compare(clientv3.CreateRevision(KeyToDelete), ">", 0)
	r, err := cli.Txn(ctx).If(NotExistsKeyToPut, ExistsKeyToDelete).
		Then(clientv3.OpDelete(KeyToDelete), clientv3.OpPut(KeyToPut, ValueToPut)).Commit()

	if r.Succeeded {
		return nil
	}

	return errors.New("Put/Delete could not be completed")

}

//KVDeleteWithPrefix to delete all the keys with the prefix key
func KVDeleteWithPrefix(key string) {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// count keys about to be deleted
	gresp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	// delete the keys
	dresp, err := cli.Delete(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all keys:", int64(len(gresp.Kvs)) == dresp.Deleted)
	// Output:
	// Deleted all keys: true
}

//KVGetWithPrefix to get all the keys with prefix key
func KVGetWithPrefix(key string) *clientv3.GetResponse {
	tlsConfig := tlsConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

//TestConfigWithTLS function to test connection with TLS config
func TestConfigWithTLS() {
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
