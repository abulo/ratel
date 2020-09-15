package kafka

import (
	"errors"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/abulo/ratel/logger"
)

const (
	ACK_BEFORE_AUTO    = 0
	ACK_AFTER_NOERROR  = 1
	ACK_AFTER_NOMATTER = 2
)

type ProducerMessage = sarama.ProducerMessage
type ProducerError = sarama.ProducerError
type ConsumerMessage = sarama.ConsumerMessage
type ConsumerError = sarama.ConsumerError
type StringEncoder = sarama.StringEncoder
type ByteEncoder = sarama.ByteEncoder

type ProducerConfig struct {
	Key           string
	Address       []string // kafka地址
	Async         bool
	ReturnSuccess bool
	ReturnError   bool
	User          string //username and password for SASL/PLAIN  or SASL/SCRAM authentication
	Password      string
	Version       *sarama.KafkaVersion // kafka version
}

type ConsumerConfig struct {
	Key          string
	Address      []string // kafka地址
	Group        string   // groupId
	Offset       int64
	Ack          int    // ack类型
	User         string //username and password for SASL/PLAIN  or SASL/SCRAM authentication
	Password     string
	DialTimeout  time.Duration // How long to wait for the initial connection.
	ReadTimeout  time.Duration // How long to wait for a response.
	WriteTimeout time.Duration // How long to wait for a transmit.
	KeepAlive    time.Duration
	Version      *sarama.KafkaVersion // kafka version
}

type ProducerMessageHandler func(msg *ProducerMessage)

type ProducerErrorHandler func(err *ProducerError)

type Producer interface {
	Close() error
	Produce(msgs ...*ProducerMessage) error
	AsyncHandle(mh ProducerMessageHandler, eh ProducerErrorHandler) // 必须设置 asyncReturnSuccess 或 asyncReturnError
}

type ConsumerMessageHandler func(msg *ConsumerMessage) error

type ConsumerErrorHandler func(err error)

type Consumer interface {
	Close() error
	// blocking to consume the messages
	Consume(topics string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) error
	ConsumeM(topics []string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) error
}

func producerConfig(opt *ProducerConfig) (config *sarama.Config) {
	config = sarama.NewConfig()
	if opt.Version != nil {
		config.Version = *opt.Version
	} else {
		config.Version = sarama.V2_6_0_0 // consumer groups require Version to be >= V0_10_2_0
	}
	if opt.User != "" { // only plain
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = opt.User
		config.Net.SASL.Password = opt.Password
	}
	config.Producer.Return.Successes = opt.ReturnSuccess
	config.Producer.Return.Errors = opt.ReturnError
	return
}

func consumerConfig(opt *ConsumerConfig) (config *sarama.Config) {
	config = sarama.NewConfig()
	if opt.Version != nil {
		config.Version = *opt.Version
	} else {
		config.Version = sarama.V2_6_0_0 // consumer groups require Version to be >= V0_10_2_0
	}
	if opt.User != "" { // only plain
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = opt.User
		config.Net.SASL.Password = opt.Password
	}
	if opt.KeepAlive > 0 {
		config.Net.KeepAlive = opt.KeepAlive
	}
	if opt.DialTimeout > 0 {
		config.Net.DialTimeout = opt.DialTimeout
	}
	if opt.ReadTimeout > 0 {
		config.Net.ReadTimeout = opt.ReadTimeout
	}
	if opt.WriteTimeout > 0 {
		config.Net.WriteTimeout = opt.WriteTimeout
	}
	if opt.Offset < 0 { // only -1(OffsetNewest), -2(OffsetOldest)
		if opt.Offset == sarama.OffsetNewest {
			config.Consumer.Offsets.Initial = sarama.OffsetNewest
		} else if opt.Offset == sarama.OffsetOldest {
			config.Consumer.Offsets.Initial = sarama.OffsetOldest
		} else {
			logger.Logger.Panic("invalid initial offset")
		}
	}
	return
}

var producers map[string]Producer = make(map[string]Producer)
var consumers map[string]Consumer = make(map[string]Consumer)

func GetProducer(key string) Producer {
	if rt, ok := producers[key]; ok {
		return rt
	}
	return nil
}

func GetConsumer(key string) Consumer {
	if rt, ok := consumers[key]; ok {
		return rt
	}
	return nil
}

func SetupProducer(opt *ProducerConfig) (err error) {

	keys := strings.Split(opt.Key, ",")
	for _, k := range keys {
		if _, ok := producers[k]; ok {
			err = errors.New("duplicate kafka producer key " + k)
			return
		}
	}
	var p Producer
	if opt.Async {
		p, err = newSaramaAsyncProducer(opt)
	} else {
		p, err = newSaramaSyncProducer(opt)
	}
	if err != nil {
		return
	}
	for _, k := range keys {
		producers[k] = p
	}

	return
}

func SetupConsumer(opt *ConsumerConfig) (err error) {

	keys := strings.Split(opt.Key, ",")
	for _, k := range keys {
		if _, ok := consumers[k]; ok {
			err = errors.New("duplicate kafka consumer key " + k)
			return
		}
	}

	c, err := newSaramaConsumerGroup(opt)
	if err != nil {
		return
	}
	for _, k := range keys {
		consumers[k] = c
	}
	return
}

func Close() {
	for _, p := range producers {
		p.Close()
	}
	for _, c := range consumers {
		c.Close()
	}
}
