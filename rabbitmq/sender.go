package rabbitmq

import (
	"github.com/streadway/amqp"
)

func (c *RabbitMQClient) Send(routeKey string, msg []byte) error {
	err := c.Ch.Publish(
		c.ExName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)

	return err
}
