package rabbitmq

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func ExampleDial() {
	conn, err := Dial("amqp://guest:guest@localhost:5672/", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("connected!")
	// Output: connected!
}

func TestConsumerTag(t *testing.T) {
	consumerTag := NewDefaultSAdder()
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(consumerTag())
		}()
	}
	time.Sleep(5 * time.Second)
}

func TestAbsReceiveListener_Remove(t *testing.T) {
	conn := getConnection()
	go func() {
		producer := conn.Producer()
		for i := 0; i < 10; i++ {
			err := producer.Send(
				"amq.direct",
				"key.direct",
				[]byte("test"+strconv.Itoa(i)),
				NewSendOptsBuilder().SetRetryable(DefaultTimesRetry()).Build(),
			)
			if err != nil {
				panic(err)
			}
		}
	}()
	consumer := conn.Consumer()
	var cnt1 = 0
	consumer.Receive(
		"queue.direct",
		NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&AbsReceiveListener{ConsumerMethod: func(d *amqp.Delivery) (brk bool) {
			if cnt1 == 5 {
				fmt.Println("receiver 1: break")
				return true
			}
			fmt.Println("receiver 1:", string(d.Body))
			cnt1++
			_ = d.Ack(false)
			return
		}},
	)
	var cnt2 = 0
	consumer.Receive(
		"queue.direct",
		NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&AbsReceiveListener{ConsumerMethod: func(d *amqp.Delivery) (brk bool) {
			if cnt2 == 1 {
				// 提前返回，不消费
				return
			}
			fmt.Println("receiver 2:", string(d.Body))
			_ = d.Ack(false)
			cnt2++
			return
		}},
	)
	var cnt3 = 0
	consumer.Receive(
		"queue.direct",
		NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&AbsReceiveListener{ConsumerMethod: func(d *amqp.Delivery) (brk bool) {
			if cnt3 == 2 {
				fmt.Println("receiver 3: break")
				return true
			}
			fmt.Println("receiver 3:", string(d.Body))
			cnt3++
			_ = d.Ack(false)
			return
		}},
	)
	time.Sleep(time.Second * 2)
	expectLen := 1
	actualLen := len(conn.operations)
	if actualLen != expectLen {
		log.Fatalf("Fail to remove! Expect: %v, actual: %v", expectLen, actualLen)
	}
}
