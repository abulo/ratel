module github.com/abulo/ratel

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.27.0
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20200914085616-93a5e1c059ad
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/sse v0.1.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis/v8 v8.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.11
	github.com/jlaffaye/ftp v0.0.0-20200812143550-39e3779af0db
	github.com/json-iterator/go v1.1.10
	github.com/mattn/go-isatty v0.0.12
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/shirou/gopsutil v2.20.8+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/streadway/amqp v1.0.0
	github.com/tsuna/gohbase v0.0.0-20200831170559-79db14850535
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/ugorji/go/codec v1.1.8
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.32.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0
