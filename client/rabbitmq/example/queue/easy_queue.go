package main

import (
	"log"

	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", rabbitmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queueBuilder := conn.QueueBuilder()
	queue := queueBuilder.
		SetQueueDeclareOpts(func(builder *rabbitmq.QueueDeclareOptsBuilder) *rabbitmq.QueueDeclareOpts {
			return builder.SetNowait(false).SetDurable(true).SetArgs(nil).Build()
		}).
		SetQueueBindOpts(func(builder *rabbitmq.QueueBindOptsBuilder) *rabbitmq.QueueBindOpts {
			return builder.SetNoWait(false).SetArgs(&amqp.Table{}).Build()
		}).
		SetRetryable(rabbitmq.DefaultTimesRetry()).
		Build()

	err = queue.DeclareAndBind("test", "key.test", "amq.direct")
	if err != nil {
		log.Fatal(err)
	}
}
