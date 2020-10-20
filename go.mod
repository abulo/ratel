module github.com/abulo/ratel

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/HdrHistogram/hdrhistogram-go v0.9.0 // indirect
	github.com/Shopify/sarama v1.27.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/disintegration/imaging v1.6.2
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20201007143536-4b4020669208
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/sse v0.1.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.3.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/h2non/bimg v1.1.4
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.11
	github.com/jlaffaye/ftp v0.0.0-20201019174721-a50ad6ffd6d5
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.11.1 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/mapstructure v1.3.3
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.8.0
	github.com/samuel/go-zookeeper v0.0.0-20200724154423-2164a8ac840e // indirect
	github.com/shirou/gopsutil v2.20.9+incompatible
	github.com/sirupsen/logrus v1.7.0
	github.com/streadway/amqp v1.0.0
	github.com/tsuna/gohbase v0.0.0-20201006203713-f1ffe9f66b83
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	github.com/ugorji/go/codec v1.1.13
	go.mongodb.org/mongo-driver v1.4.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.3.1-0.20200828183125-ce943fd02449 // indirect
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	golang.org/x/tools v0.0.0-20201014231627-1610a49f37af // indirect
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	modernc.org/mathutil v1.1.1 // indirect
	modernc.org/strutil v1.1.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0
