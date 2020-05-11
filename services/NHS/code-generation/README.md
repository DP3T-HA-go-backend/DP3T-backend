## Simple Go implementation of DP3T's `/v1/authcode`

This implementation is just meant to be a demo for testing purposes.
Data is stored in memory and not persisted.

A 'ec256-key' ECS256 PRIVATE key needs to be placed in the folder. It needs to match the PUBKEY in the DP3T App module.

Building the service requires Go:

```
make
```

Run:

```
./authcode
make get
```

### Docker

Build the Docker image and run the container as follows:

```
make docker-image
make docker-run
```
