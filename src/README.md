dp3t-backend
============

A DP3T-compatible backend implemented in Go, using etcd as data store.

## Services

- `exposed`: Provides a service with the same interface as [DP3T
  Backend][dp3t-sdk-backend], which can be used to:
  - Get list of exposees for a particular batch.
  - Post new exposee.
- `authcode`: Generate authcodes, which are needed to validate new exposees.

## Testing

We provide a docker-compose environment for testing the services
locally, using etcd as data store. The only configuration needed to run
etcd is the generation of certificates/keys for client-to-server
connections:

```
cd config/test/etcd/
make
```

Run the integration tests as follows:

```
docker-compose up -d
make test
docker-compose down
```

[dp3t-sdk-backend]: https://github.com/DP-3T/dp3t-sdk-backend
