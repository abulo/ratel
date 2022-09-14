package main

import (
	"log"
	"time"

	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", rabbitmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer := conn.Consumer()
	consumer.Receive(
		"queue.direct",
		rabbitmq.NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&rabbitmq.AbsReceiveListener{
			ConsumerMethod: func(delivery *amqp.Delivery) (brk bool) {
				log.Println("queue.direct ", delivery.DeliveryTag, " ", string(delivery.Body))
				err := delivery.Ack(false)
				if err != nil {
					log.Println(err)
				}
				return
			},
			FinishMethod: func(err error) {
				if err != nil {
					// 处理错误
					log.Fatal(err)
				}
				// defer xxx.close() // 关闭资源操作等
			},
		})

	time.Sleep(time.Minute)
}
