package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// KVPut to put key and value
func KVPut(cliconfig *clientv3.Config, key string, value string) {
	cli, err := clientv3.New(*cliconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	_, err = cli.Put(context.TODO(), key, value)
	if err != nil {
		log.Fatal(err)
	}
}

// KVPutTTL to put key and value with a Time To Live in days
func KVPutTTL(cliconfig *clientv3.Config, key string, value string, days int64) {
	cli, err := clientv3.New(*cliconfig)
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
func KVGet(cliconfig *clientv3.Config, key string, requestTimeout time.Duration) *clientv3.GetResponse {
	cli, err := clientv3.New(*cliconfig)
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

// KVDelete to delete a key
func KVDelete(cliconfig *clientv3.Config, key string, requestTimeout time.Duration) {
	cli, err := clientv3.New(*cliconfig)
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

// KVPut only if the key did not exist
func KVPutIfNotExists(cliconfig *clientv3.Config, putNamespace string, KeyToPut string, ValueToPut string, requestTimeout time.Duration) error {
	cli, err := clientv3.New(*cliconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	NotExistsKeyToPut := clientv3.Compare(clientv3.CreateRevision(putNamespace+KeyToPut), "=", 0)
	r, err := cli.Txn(ctx).If(NotExistsKeyToPut).Then(clientv3.OpPut(putNamespace+KeyToPut, ValueToPut)).Commit()

	if r.Succeeded {
		return nil
	}

	return errors.New("Key already existed")
}

// KVDelete only if the key did existed
func KVDeleteIfExists(cliconfig *clientv3.Config, KeyToDelete string, requestTimeout time.Duration) error {
	cli, err := clientv3.New(*cliconfig)
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

// KVDelete one existing Key and KVPut another one only if the first existed and was deleted
func KVPutAndDelete(cliconfig *clientv3.Config, deleteNamespace string, KeyToDelete string, putNamespace string,
	KeyToPut string, ValueToPut string, TTL int64, requestTimeout time.Duration) error {
	cli, err := clientv3.New(*cliconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// minimum lease TTL is in seconds
	resp, err := cli.Grant(ctx, TTL)
	if err != nil {
		log.Fatal(err)
	}

	NotExistsKeyToPut := clientv3.Compare(clientv3.CreateRevision(putNamespace+KeyToPut), "=", 0)
	ExistsKeyToDelete := clientv3.Compare(clientv3.CreateRevision(deleteNamespace+KeyToDelete), ">", 0)
	r, err := cli.Txn(ctx).If(NotExistsKeyToPut, ExistsKeyToDelete).
		Then(clientv3.OpDelete(deleteNamespace+KeyToDelete), clientv3.OpPut(putNamespace+KeyToPut, ValueToPut, clientv3.WithLease(resp.ID))).Commit()

	if r.Succeeded {
		return nil
	}

	return errors.New("Put/Delete could not be completed")

}

// KVDeleteWithPrefix to delete all the keys with the prefix key
func KVDeleteWithPrefix(cliconfig *clientv3.Config, key string, requestTimeout time.Duration) {
	cli, err := clientv3.New(*cliconfig)
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

// KVGetWithPrefix to get all the keys with prefix key
func KVGetWithPrefix(cliconfig *clientv3.Config, key string, requestTimeout time.Duration) *clientv3.GetResponse {
	cli, err := clientv3.New(*cliconfig)
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

func KVGetAllKeys(cliconfig *clientv3.Config, keyNamespace string, requestTimeout time.Duration) *clientv3.GetResponse {
	cli, err := clientv3.New(*cliconfig)

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	//resp, err := cli.Get(ctx, key, clientv3.WithRange("\x00"))
	resp, err := cli.Get(ctx, keyNamespace, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil
	}

	return resp
}

// KVGetWithRange to get all the keys within a range  [key, end).
func KVDeleteAllKeys(cliconfig *clientv3.Config, requestTimeout time.Duration) error {
	cli, err := clientv3.New(*cliconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	resp, err := cli.Get(ctx, "\x00", clientv3.WithRange("\x00"))
	cancel()

	if err != nil {
		return errors.New("couldn't delete all keys")
	}

	for _, key := range resp.Kvs {
		cli.Delete(context.TODO(), key.String())
	}

	return nil
}
