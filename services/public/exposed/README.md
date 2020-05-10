## Simple Go implementation of DP3T's `/v1/exposed`

This implementation is just meant to be a demo for testing purposes.
Data is stored in memory and not persisted.

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
