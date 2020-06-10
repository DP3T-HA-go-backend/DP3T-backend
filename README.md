dp3t-backend
============

A DP3T-compatible backend implemented in Go, using etcd as data store.
This backend provides the following services:

- `exposed`: Provides a service with the same interface as [DP3T
  Backend][dp3t-sdk-backend], which can be used to:
  - Get list of exposees for a particular batch.
  - Post new exposee.
- `authcode`: Generate authcodes, which are needed to validate new exposees.

## Build

This project requires:
- Go (>= 1.11)
- `protoc`: Protocol Buffers compiler
- `protoc-gen-go`: Plugin for the Protocol Buffers compiler to generate Go
  code

To build the services in your local environment, you can use the provided
Makefile:
```
make
```

It's also possible build the project as Docker images, without installing any
of the dependencies locally. The following command will build the
`dp3t.exposed.protobuf.api` and `dp3t.authcode.api` images:
```
make docker
```

## Run

To run the services, you'll first need to generate your own public/private EC
keys.
```
make config/ec256-key
```

And then simply run the servers as follows. By default, servers will read
the configuration file under `config/{authcode,exposed}.ini`, but the path
can be configured with the `--config` flag.
```
./server/{authcode,exposed}/main [--config path/to/config/file.ini]
```

Configurations are defined as INI files with the following mandatory fields:
```
port = NUMBER
private-key-file = /path/to/private/key/file
public-key-file = /path/to/public/key/file
store = inmem|etcd
```

The `store` field defines the data store to be used, and currently only
supports one of these:
- `inmem`: In-memory. Used mainly for developement purposes, and only
  partially implemented.
- `etcd`: Stores data on [etcd][etcd]. Requires the following additional
  section in the configuration file in order to configure etcd's entripoints
  and TLS certificates/keys:
```
[etcd]
endpoints = 0.0.0.0:2379
cert-file = /path/to/etcd/server.crt
key-file = /path/to/etcd/server.key
ca-file = /path/to/etcd/ca.crt
```

## Test

We provide a docker-compose environment for testing the services locally,
using etcd as data store. Docker images for each service are also available
and can be built with the `Makefile`, and the only configuration needed to run
the tests is generating the EC keys for the servers, as well as generating the
TLS certificates/keys for etcd's client-to-server connections:
```
make docker
make test-keys
```

Run the integration tests as follows:
```
docker-compose up -d
make test
docker-compose down
```

[dp3t-sdk-backend]: https://github.com/DP-3T/dp3t-sdk-backend
[etcd]: https://etcd.io/
