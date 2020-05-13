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

docker: docker-exposed docker-authcode

docker-exposed:
	docker build -f docker/Dockerfile.exposed -t dp3t.public.exposed:latest .

docker-authcode:
	docker build -f docker/Dockerfile.authcode -t dp3t.public.authcode:latest .

clean:
	rm -f $(BUILD_BIN)
	rm -f $(PB_GO)

.PHONY: docker clean
