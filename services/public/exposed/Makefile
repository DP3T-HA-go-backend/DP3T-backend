URL=http://localhost:8080/v1

exposed: main.go exposed.pb.go
	go build

exposed.pb.go: exposed.proto
	PATH=$(PATH):$(HOME)/go/bin protoc --go_out=. exposed.proto

post:
	@echo "POST $(URL)/exposed"
	@echo 'key:"SAMPLE_KEY"\nkeyDate:1234' | \
		protoc --encode main.ProtoExposee exposed.proto | \
		curl -sS -X POST --data-binary @- $(URL)/exposed

get:
	@echo "GET $(URL)/exposed/$$(date +"%s")"
	@curl -sS -X GET $(URL)/exposed/$$(date +"%s") | \
		protoc --decode main.ProtoExposedList exposed.proto

clean:
	rm -f exposed
	rm -f exposed.pb.go

.PHONY: post get clean
