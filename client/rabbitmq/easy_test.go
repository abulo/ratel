package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func ExampleConsumer_Receive() {
	conn := getConnection()

	consumer := conn.Consumer()
	consumer.Receive(
		"queue.direct",
		NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&AbsReceiveListener{
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
					return
				}
				// defer xxx.close() // 关闭资源操作等
			},
		})

	time.Sleep(time.Minute) // 由于 Consumer.Receive() 内部采用了异步方式处理，因此 Receive 方法不会阻塞等待
}
