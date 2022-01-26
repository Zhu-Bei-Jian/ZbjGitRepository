module sanguosha.com/sgs_herox

go 1.12

replace sanguosha.com/baselib => ./_depends/sanguosha.com/baselib

replace github.com/everalbum/redislock => ./_depends/github.com/everalbum/redislock

replace github.com/tidwall/buntdb => ./_depends/github.com/tidwall/buntdb

require (
	github.com/Shopify/sarama v1.24.1
	github.com/bshuster-repo/logrus-logstash-hook v0.4.1 // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/everalbum/redislock v0.0.0-00010101000000-000000000000
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.1
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f // indirect
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	github.com/nats-io/nats.go v1.13.0
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/wanghuiyt/ding v0.0.1
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e // indirect
	google.golang.org/protobuf v1.23.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	sanguosha.com/baselib v0.0.0-00010101000000-000000000000
)
