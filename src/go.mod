module dp3t-backend

require (
	github.com/coreos/etcd v3.3.20+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20200316104309-cb8b64719ae3 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.1
	github.com/google/uuid v1.1.1 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	go.etcd.io/etcd v3.3.20+incompatible
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/grpc v1.26.0 // indirect
	google.golang.org/protobuf v1.22.0
	gopkg.in/dgrijalva/jwt-go.v3 v3.2.0
	gopkg.in/ini.v1 v1.56.0

)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
