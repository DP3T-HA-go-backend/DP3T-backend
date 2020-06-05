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

Docker containers for each service are also available and can be built as
follows:
```
make docker
```

## Test

We provide a docker-compose environment for testing the services locally,
using etcd as data store. The only configuration needed to run the tests is
generating the EC keys for the servers, as well as generating the TLS
certificates/keys for etcd's client-to-server connections:

```
make test-keys
```

Run the integration tests as follows:

```
docker-compose up -d
make test
docker-compose down
```

[dp3t-sdk-backend]: https://github.com/DP-3T/dp3t-sdk-backend
