package main

import (
	"log"
	"time"

	"github.com/abulo/ratel/v3/client/rabbitmq"
)

func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error

	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", nil)
	onErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	onErr(err)
	defer ch.Close()

	// 设置为 Plain Persistent （无格式/持久化）的形式发送消息
	err = ch.SendOpts(
		"amq.direct", "key.direct", []byte("Send() | "+time.Now().Format("2006-01-02 15:04:05")),
		rabbitmq.NewSendOptsBuilder().SetMessageFactory(rabbitmq.MessagePlainPersistent).Build(),
	)
	onErr(err)
}
