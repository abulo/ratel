package kafka

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

//Config 数据库配置
type Config struct {
	Username string //账号 root
	Password string //密码
	Host     string //host localhost
	Broker   []string
}

type Clientkafka struct {
	host   string
	writer *Writer
	reader *Reader
	dialer *kafka.Dialer
	broker []string
}
type Writer struct {
	topic map[string]*ClientWriter
	mu    sync.RWMutex
}
type ClientWriter struct {
	*kafka.Writer
}

type Reader struct {
	topic map[string]*ClientReader
	mu    sync.RWMutex
}
type ClientReader struct {
	*kafka.Reader
}

//New 新连接
func New(config *Config) *Clientkafka {

	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
	}

	if !util.Empty(config.Username) && !util.Empty(config.Password) {

		mechanism, err := scram.Mechanism(scram.SHA512, config.Username, config.Password)
		if err != nil {
			logger.Logger.Panic(err)
		}
		dialer.DualStack = true
		dialer.SASLMechanism = mechanism
	}
	writer := &Writer{
		topic: make(map[string]*ClientWriter),
	}
	reader := &Reader{
		topic: make(map[string]*ClientReader),
	}
	return &Clientkafka{
		host:   config.Host,
		writer: writer,
		reader: reader,
		dialer: dialer,
		broker: config.Broker,
	}
}

func (client *Clientkafka) Kafka() *kafka.Conn {
	conn, err := kafka.Dial("tcp", client.host)
	if err != nil {
		logger.Logger.Panic(err)
	}
	return conn
}

func (client *Clientkafka) Host() string {
	return client.host
}

//ListTopics To list topics
func (client *Clientkafka) ListTopics() []string {
	partitions, err := client.Kafka().ReadPartitions()
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	list := make([]string, 0)
	for _, p := range partitions {
		list = append(list, p.Topic)
	}
	return util.ArrayStringUniq(list)
}

//CreateTopic create
func (client *Clientkafka) CreateTopic(topic string) error {
	controller, err := client.Kafka().Controller()
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

//NameSpace 获取分组
func (proxy *Writer) NameSpace(topic string) *ClientWriter {
	proxy.mu.RLock()
	res := proxy.topic[topic]
	proxy.mu.RUnlock()
	return res
}

//SetNameSpace 设置组
func (proxy *Writer) SetNameSpace(group string, client *ClientWriter) {
	proxy.mu.Lock()
	proxy.topic[group] = client
	proxy.mu.Unlock()
}

//NameSpace 获取分组
func (proxy *Reader) NameSpace(topic string) *ClientReader {
	proxy.mu.RLock()
	res := proxy.topic[topic]
	proxy.mu.RUnlock()
	return res
}

//SetNameSpace 设置组
func (proxy *Reader) SetNameSpace(group string, client *ClientReader) {
	proxy.mu.Lock()
	proxy.topic[group] = client
	proxy.mu.Unlock()
}

// NewWriter ..
func (client *Clientkafka) NewWriter(topic string) *ClientWriter {
	if conn, ok := client.writer.topic[topic]; ok {
		return conn
	}
	newWrite := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  client.broker,
		Topic:    topic,
		Balancer: &kafka.Hash{},
		Dialer:   client.dialer,
	})
	newClient := &ClientWriter{
		newWrite,
	}
	client.writer.SetNameSpace(topic, newClient)
	return newClient
}

func (cw *ClientWriter) SendMessage(ctx context.Context, msg ...kafka.Message) error {

	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("kafka", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "kafka")
			span.SetTag("method", "SendMessage")
			span.LogFields(log.String("topic", cw.Stats().Topic))
			span.LogFields(log.Object("message", msg))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	return cw.Writer.WriteMessages(ctx, msg...)
}

func (client *Clientkafka) NewReader(topic string) *ClientReader {
	if conn, ok := client.reader.topic[topic]; ok {
		return conn
	}
	newReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   client.broker,
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	newClient := &ClientReader{
		newReader,
	}

	client.reader.SetNameSpace(topic, newClient)
	return newClient
}

func (cr *ClientReader) ReadMessage(ctx context.Context) (kafka.Message, error) {

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("kafka", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "kafka")
			span.SetTag("method", "ReadMessage")
			span.LogFields(log.String("topic", cr.Config().Topic))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return cr.Reader.ReadMessage(ctx)
}
