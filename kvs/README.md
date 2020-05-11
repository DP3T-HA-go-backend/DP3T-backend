[Back to Index](../README.md)

# KVS schemas, structure, configuration parameters


## Go and etcd

- ### `etcd/clientv3` is the Go official etcd client
- ### Install

```bash
> go get go.etcd.io/etcd/clientv3
```

## Command line to test etcd keys
```bash
ETCDCTL_API=3 etcdctl --endpoints=https://10.0.26.10:2379,https://10.0.26.11:2379,https://10.0.26.13:2379 --cacert=/etc/ssl/etcd/ssl/ca.pem --cert=/etc/ssl/etcd/ssl/node-node1.pem --key=/etc/ssl/etcd/ssl/node-node1-key.pem
```

## Functions

```go
//KVPut to put key and value
func KVPut(key string, value string)
```

```go
//KVPutTTL to put key and value with a Time To Live in days
func KVPutTTL(key string, value string, days int64) 
```


```go
// KVGet to Get a key
func KVGet(key string) *clientv3.GetResponse
```


```go
//KVDelete to delete a key
func KVDelete(key string)
```


```go
//KVDeleteWithPrefix to delete all the keys with the prefix key
func KVDeleteWithPrefix(key string)
```


```go
//KVGetWithPrefix to get all the keys with prefix key
func KVGetWithPrefix(key string) *clientv3.GetResponse
```


```go
//TestConfigWithTLS function to test connection with TLS config
func TestConfigWithTLS() 
```