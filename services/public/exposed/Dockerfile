FROM golang:buster AS builder

RUN mkdir /build
COPY . /build

RUN apt-get update && \
    apt-get -y install protobuf-compiler && \
    go get github.com/julienschmidt/httprouter && \
    go get github.com/dgrijalva/jwt-go && \
    #go get github.com/gbrlsnchs/jwt/v3 && \
    go get google.golang.org/protobuf/proto && \
    go get github.com/golang/protobuf/proto && \
    go get google.golang.org/protobuf/encoding/protojson && \
    go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    cd /build && \
    make

FROM alpine

RUN apk add --no-cache libc6-compat

COPY --from=builder /build/exposed /service/exposed
COPY ec256-key /ec256-key

EXPOSE 8080/tcp
ENTRYPOINT ["/service/exposed"]
