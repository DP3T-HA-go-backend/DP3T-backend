URL=http://localhost:8080

exposed: main.go 
	go build -o authcode

exposed.pb.go: exposed.proto
	PATH=$(PATH):$(HOME)/go/bin protoc --go_out=. exposed.proto


get:
	@echo "GET $(URL)/authcode"
	@curl -sS -X GET $(URL)/authcode 

docker-image:
	docker build . -t dp3t.authcode.api:latest

docker-run:
	docker run -p 8080:8080 --rm -it dp3t.authcode.api

clean:
	rm -f authcode

.PHONY: post get clean
