BUILD_SRC = $(shell find server -type f -name 'main.go')
BUILD_BIN = $(BUILD_SRC:main.go=main)

PROTO = $(shell find api -type f -name '*.proto')
PB_GO = $(PROTO:.proto=.pb.go)

default: proto build

build: $(BUILD_BIN)

%main: %main.go
	go build -o $@ ./$<

proto: $(PB_GO)

%.pb.go: %.proto
	PATH=$(PATH):$(HOME)/go/bin protoc --go_out=. $<


# Test -----------------------------------------------------------------

test:
	go test ./test/ -count=1 -v


# Docker images --------------------------------------------------------

docker: docker-exposed docker-authcode

docker-exposed:
	docker build -f docker/Dockerfile.exposed -t dp3t.exposed.protobuf.api:latest .

docker-authcode:
	docker build -f docker/Dockerfile.authcode -t dp3t.authcode.api:latest .


# Keys for development & testing  --------------------------------------

config/ec256-key:
	openssl ecparam -genkey -name prime256v1 -noout -out $@
	openssl ec -in $@ -pubout -out $@.pub

config/test/ec256-key:
	openssl ecparam -genkey -name prime256v1 -noout -out $@
	openssl ec -in $@ -pubout -out $@.pub

test-keys: config/test/ec256-key
	cd config/test/etcd/ && make && cd -


# ----------------------------------------------------------------------

clean:
	rm -f $(BUILD_BIN)
	rm -f $(PB_GO)

.PHONY: test docker clean
