module github.com/abulo/ratel

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/HdrHistogram/hdrhistogram-go v0.9.0 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/abulo/clickhouse-go v1.4.8
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/disintegration/imaging v1.6.2
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20210519083322-55daf7425ecb
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/sse v0.1.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-redis/redis/v8 v8.8.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/h2non/bimg v1.1.5
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/imdario/mergo v0.3.12
	github.com/jlaffaye/ftp v0.0.0-20210307004419-5d4190119067
	github.com/json-iterator/go v1.1.11
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.0
	github.com/prometheus/client_golang v1.10.0
	github.com/segmentio/kafka-go v0.4.16
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	github.com/tsuna/gohbase v0.0.0-20201125011725-348991136365
	github.com/uber/jaeger-client-go v2.29.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	github.com/ugorji/go/codec v1.2.6
	github.com/valyala/fasthttp v1.25.0
	go.mongodb.org/mongo-driver v1.5.2
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.3.1-0.20200828183125-ce943fd02449 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20210212180131-e7f2df4ecc2d
	google.golang.org/grpc v1.38.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0
