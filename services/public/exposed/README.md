## Simple Go implementation of DP3T's `/v1/exposed`

This implementation is just meant to be a demo for testing purposes.
Data is stored in memory and not persisted.

A 'ec256-key' ECS256 PRIVATE key needs to be placed in the folder. It needs to match the PUBKEY in the DP3T App module.

Building the service requires Go and Protocol Buffers:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go
make
```

Run:

```
./exposed
make post
make get
```
