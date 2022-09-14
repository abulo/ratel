package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

// --- EXAMPLE ---

// NOTE: 执行 EXAMPLE 前，本地必须已经启动了一个 amqp 服务器

func ExampleConnection_Dial_1_common() {
	conn := NewConnection("amqp://guest:guest@localhost:5672/", nil)
	defer conn.Close()
	if err := conn.Dial(); err != nil {
		panic(err)
	}
	fmt.Println("connected!")
	// Output: connected!
}

func ExampleConnection_Dial_2_reconnected() {
	conn := NewConnection("amqp://guest:guest@localhost:5672/", DefaultTimesRetry())
	defer conn.Close()

	if err := conn.Dial(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("reconnected!")
	// Output: reconnected!
}

func ExampleConnection_RegisterAndExec() {
	var err error
	conn := NewConnection("amqp://guest:guest@localhost:5672/", &TimesRetry{Interval: 3 * time.Second, RetryTimes: 10})
	defer conn.Close()

	err = conn.Dial()
	if err != nil {
		log.Fatal(err)
	}

	consumerTag := NewDefaultSAdder()
	conn.RegisterAndExec(func(key string, ch *Channel) {
		deliverys, e := ch.Consume("queue.direct", consumerTag(), true, false, false, false, nil)
		if e != nil {
			log.Fatal(e)
		}
		for delivery := range deliverys {
			log.Println("queue.direct-1 ", delivery.DeliveryTag, " ", string(delivery.Body))
		}
	})
	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Receive("queue.direct", func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct-2 ", delivery.DeliveryTag, " ", string(delivery.Body))
			return true
		})
		if e != nil {
			log.Fatal(e)
		}
	})

	//注意：conn.RegisterAndExec 会导致重连时再次执行 Operation 中的操作。如果不希望发生此情况，应该使用 Connect.Channel 消费消息
	conn.RegisterAndExec(func(key string, ch *Channel) {
		delivery, _, e := ch.Get("queue.direct", true)
		if e != nil {
			log.Fatal(e)
		}
		log.Println("queue.direct ", delivery.DeliveryTag, " ", string(delivery.Body))
	})

	//注意：conn.RegisterAndExec 会导致重连时再次执行 Operation 中的操作。如果不希望发生此情况，应该使用 Connect.Channel 发送消息
	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Send("amq.direct", "key.direct", []byte("rabbitmq addOperation() send test!"))
		if e != nil {
			log.Fatal(e)
		}
	})
}

// --- TEST ---

func TestConnection(t *testing.T) {
	conn, err := Dial(defaultURL, nil)
	onErr(err)
	defer conn.Close()
}

func TestConnection_Channel(t *testing.T) {
	var err error

	conn, err := Dial(defaultURL, nil)
	onErr(err)
	defer conn.Close()

	t.Log("connected...")
	time.Sleep(10 * time.Second)
	t.Log("sleep end...")
	_, err = conn.Channel()
	onErr(err)
}

// 测试总是重连
func TestConnection_Dial_always(t *testing.T) {
	var err error

	conn, err := Dial(defaultURL, &TimesRetry{Always: true, Interval: 3 * time.Second})
	onErr(err)
	defer conn.Close()

	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Receive("queue.direct", func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct-1 ", delivery.DeliveryTag, " ", string(delivery.Body))
			return
		})
		onErr(e)
	})
	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Receive("queue.direct", func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct-2 ", delivery.DeliveryTag, " ", string(delivery.Body))
			return
		})
		onErr(e)
	})

	time.Sleep(time.Minute)
}

// 测试按指定次数重连
func TestConnection_Dial_byTimes(t *testing.T) {
	var err error

	conn, err := Dial(defaultURL, &TimesRetry{Interval: 3 * time.Second, RetryTimes: 10})
	onErr(err)
	defer conn.Close()

	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Receive("queue.direct", func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct-1 ", delivery.DeliveryTag, " ", string(delivery.Body))
			return
		})
		onErr(e)
	})
	conn.RegisterAndExec(func(key string, ch *Channel) {
		e := ch.Receive("queue.direct", func(delivery *amqp.Delivery) (brk bool) {
			log.Println("queue.direct-2 ", delivery.DeliveryTag, " ", string(delivery.Body))
			return
		})
		onErr(e)
	})

	time.Sleep(time.Minute)
}
