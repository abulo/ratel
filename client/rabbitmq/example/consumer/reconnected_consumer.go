package main

import (
	"log"
	"time"

	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

// 断线重连，循环消费
func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", rabbitmq.DefaultTimesRetry())
	onErr(err)
	defer conn.Close()

	consumerTag := rabbitmq.NewDefaultSAdder()
	// amqp 原生方法消费
	conn.RegisterAndExec(func(key string, ch *rabbitmq.Channel) {
		// NOTE: 可在此使用 defer 语句关闭某些资源
		deliverys, e := ch.Consume("queue.direct", consumerTag(), true,
			false, false, false, nil)
		if e != nil {
			log.Fatal(e)
		}
		for delivery := range deliverys {
			log.Println("queue.direct-1 ", delivery.DeliveryTag, " ", string(delivery.Body))
			//if 满足某种条件 {
			//	ch.RemoveOperation(key) // 主动退出时需要移除当前监听器
			//	break
			//}
		}
	})

	// rabbitmq 方法消费
	conn.RegisterAndExec(func(key string, ch *rabbitmq.Channel) {
		e := ch.ReceiveOpts(
			"queue.direct",
			func(delivery *amqp.Delivery) (brk bool) {
				log.Println("queue.direct-2 ", delivery.DeliveryTag, " ", string(delivery.Body))
				e := delivery.Ack(false)
				if e != nil {
					log.Fatal(e)
				}
				return
			},
			rabbitmq.NewReceiveOptsBuilder().
				SetAutoAck(false).
				Build(),
		)
		if e != nil {
			log.Fatal(e)
		}
	})

	time.Sleep(time.Minute)
}
