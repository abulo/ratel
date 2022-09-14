package rabbitmq

import (
	"context"
	"log"

	"github.com/abulo/ratel/v3/logger"
)

var defaultURL = "amqp://guest:guest@localhost:5672/"

func init() {
	log.Println("init test...")
	logger.Logger.Info("debug")

	c, err := Dial(defaultURL, DefaultTimesRetry())
	onErr(err)
	defer c.Close()

	ch, err := c.Channel()
	onErr(err)

	// declare queue
	_, err = ch.QueueDeclare("queue.direct", true, false, false, false, nil)
	onErr(err)
	err = ch.QueueBind("queue.direct", "key.direct", "amq.direct", false, nil)
	onErr(err)
}

// ---- utils ----

func onErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getConnection() *Connection {
	c, err := Dial(defaultURL, DefaultTimesRetry())
	onErr(err)
	return c
}

func getChannel() (*Channel, *Connection) {
	var err error

	conn := getConnection()
	channel, err := conn.Channel()
	onErr(err)

	return channel, conn
}

func getChannelWithContext() (*Channel, context.CancelFunc, *Connection) {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := Dial(defaultURL, DefaultCtxRetry(ctx))
	onErr(err)
	channel, err := conn.Channel()
	onErr(err)

	return channel, cancel, conn
}
