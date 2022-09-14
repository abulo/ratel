package main

import (
	"log"
	"time"

	"github.com/abulo/ratel/v3/client/rabbitmq"
)

func main() {
	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", rabbitmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	producer := conn.Producer()
	err = producer.Send("amq.direct", "key.direct",
		[]byte("producer.Send() | "+time.Now().Format("2006-01-02 15:04:05")),
		rabbitmq.DefaultSendOpts())
	if err != nil {
		log.Fatal(err)
	}
}
