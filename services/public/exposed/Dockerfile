FROM golang:buster AS builder

RUN mkdir /build
COPY . /build

RUN apt-get update && \
    apt-get -y install protobuf-compiler && \
    go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    cd /build && \
    make clean && \
    make

FROM alpine

RUN apk add --no-cache libc6-compat

COPY ec256-key /service/etc/exposed/ec256-key
COPY config-production.ini /service/etc/exposed/config.ini

COPY --from=builder /build/exposed /service/bin/exposed

EXPOSE 8080/tcp
ENTRYPOINT ["/service/bin/exposed"]
CMD ["--config", "/service/etc/exposed/config.ini"]
