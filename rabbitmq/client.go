package rabbitmq

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	Ch *amqp.Channel
	ExName string
	QueueName string
}

type Config struct {
	User string
	Password string
	Host string
	Port int
	VHost string
}

type SendConfig struct {
	Config
	ExName string
}

type RecvConfig struct {
	Config
	ExName string
	QueueName string
	RouteKey string
}

var clientMap map[string]*RabbitMQClient = make(map[string]*RabbitMQClient);
var mu sync.Mutex;

func NewSendClient(c SendConfig) (*RabbitMQClient,error) {
	key := fmt.Sprintf("Send_%s_%s_%s_%d_%s",c.User,c.Password,c.Host,c.Port,c.ExName);
	if client,ok := clientMap[key]; ok {
		return client,nil;
	}

	mu.Lock();
	if client,ok := clientMap[key]; ok {
		return client,nil;
	}

	client,err := newClient(Config{c.User,c.Password,c.Host,c.Port,c.VHost});
	if err != nil {
		return nil,err;
	}

	client.ExName = c.ExName;

	err = client.Ch.ExchangeDeclare(
		c.ExName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	);

	if err != nil {
		return nil,err;
	}

	clientMap[key] = client;
	mu.Unlock();

	return client,nil;
}

func NewRecvClient(c RecvConfig) (*RabbitMQClient,error) {
	key := fmt.Sprintf("Recv_%s_%s_%s_%d_%s_%s",c.User,c.Password,c.Host,c.Port,c.ExName,c.VHost);
	if client,ok := clientMap[key]; ok {
		return client,nil;
	}

	mu.Lock();
	if client,ok := clientMap[key]; ok {
		return client,nil;
	}

	client,err := newClient(Config{c.User,c.Password,c.Host,c.Port,c.VHost});
	if err != nil {
		return nil,err;
	}

	client.QueueName = c.QueueName;

	err = client.Ch.ExchangeDeclare(
		c.ExName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	);

	if err != nil {
		return nil,err;
	}

	_,err = client.Ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		return nil,err;
	}

	err = client.Ch.QueueBind(
		c.QueueName,
		c.RouteKey,
		c.ExName,
		false,
		nil,
	);

	clientMap[key] = client;
	mu.Unlock();

	return client,nil;
}

func newClient(c Config) (*RabbitMQClient,error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",c.User,c.Password,c.Host,c.Port,c.VHost);
	conn,err := amqp.Dial(connStr);
	if err != nil {
		return nil,err;
	}

	ch,err := conn.Channel();
	if err != nil {
		return nil,err;
	}

	client := &RabbitMQClient{
		Ch:ch,
	};

	return client,nil;
}