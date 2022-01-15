package rabbitmq

import (
	"runtime"

	"github.com/streadway/amqp"
)

type RecvHander func(msg []byte) bool

func (c *RabbitMQClient) Consume(f RecvHander, isPallal bool) error {
	msgs, err := c.Ch.Consume(
		c.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	if isPallal {
		cups := runtime.NumCPU()
		for i := 0; i < cups; i++ {
			go recvData(f, msgs)
		}
	} else {
		go recvData(f, msgs)
	}

	return nil
}

func recvData(f RecvHander, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		msg := d.Body
		res := f(msg)
		d.Ack(res)
	}
}
