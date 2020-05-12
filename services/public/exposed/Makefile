URL=http://localhost:8080

exposed: main.go exposed.pb.go
	go build -o exposed

exposed.pb.go: exposed.proto
	PATH=$(PATH):$(HOME)/go/bin protoc --go_out=. exposed.proto

test:
	go test -v

post:
	@echo "POST $(URL)/"
	@curl -sS -X POST \
		--header "Content-Type: application/json" \
		--data '{"key":"SAMPLE_KEY","keyDate":1234}' \
		$(URL)/

get:
	@echo "GET $(URL)/$$(date +"%s")"
	@curl -sS -X GET $(URL)/$$(date +"%s") | \
		protoc --decode main.ProtoExposedList exposed.proto

ec256-key:
	ssh-keygen -t ecdsa -b 256 -f ec256-key -m pkcs8 -N ''

docker-image:
	docker build . -t dp3t.exposed.protobuf.api:latest

docker-run:
	docker run -p 8080:8080 --rm -it dp3t.exposed.protobuf.api

clean:
	rm -f exposed
	rm -f exposed.pb.go

.PHONY: post get clean
