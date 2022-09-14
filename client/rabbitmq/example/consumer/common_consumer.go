package main

import (
	"log"

	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", nil)
	onErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	onErr(err)
	defer ch.Close()

	// 如果想要消费一次后立刻退出，可以使用如下方式
	delivery, ok, err := ch.Get("queue.direct", false)
	onErr(err)
	if ok {
		log.Println("queue.direct-get ", delivery.DeliveryTag, " ", string(delivery.Body))
		err := delivery.Ack(false)
		onErr(err)
	}

	// 如果期望循环消费，可以使用如下方式
	err = ch.ReceiveOpts(
		"queue.direct",
		func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct ", delivery.DeliveryTag, " ", string(delivery.Body))
			err := delivery.Ack(false)
			onErr(err)
			return
		},
		rabbitmq.NewReceiveOptsBuilder().
			SetAutoAck(false).
			Build(),
	)

	// 此后语句将不可达，除非 ReceiveOpts 产生错误
	onErr(err)
}
