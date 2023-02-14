# rabbitmq: An easy golang amqp client

Introduction
---

rabbitmq 扩展了 github.com/streadway/amqp 包，在 amqp 包的基础上实现了网络断线重连、消息发送失败重试的能力。

由于消息的发送和接收，依赖 Channel 的创建，Channel 的创建依赖与 Connection 的连接。如果断开网络连接，Connection 就会断开。 
但是 Channel 并不知道是谁创建了自己。所以通常情况下，我们有两种断线重连的方案：

* 一种是 Connection 断开后，Channel 自己重新获取 Connection。此种方式是最简单的实现方案。
  但此法的问题是，如果 Channel 数量极其庞大，
  每个 Channel 都会创建 Connection 并连接。当网络状况不大好的时候，可能会有数量极其庞大的 Channel 反复尝试重连，
  导致服务器资源占用会暴增，甚至加剧网络的阻塞。但其实我们只需要一个 Connection 去判断网络是否能连接、已连接。
* 一种是 Connection 断开后，Connection 自己重连。如果重连不成功，Channel 不做任何操作。如果重连成功，自动重跑注册的操作。
  而如果使用断线重连发送消息，将在 Connection 连接成功后继续发送需要发送的消息。此种方式实现比较复杂，
  但避免了服务器资源占用以及加剧网络阻塞的问题。
  
rabbitmq 采用的是后者。

Feature
---

目前实现了

* 断线重连
* 消息重发
* 同步的消息发送确认

Usage
---

```go
package main

import (
  "log"
  "time"
  "github.com/abulo/ratel/v2/client/rabbitmq"
  "github.com/streadway/amqp"
)

func main() {
  // -- 创建 Connection 并连接服务器 --
  conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/", rabbitmq.DefaultTimesRetry())
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  // -- 发送消息 --
  go func() {
    producer := conn.Producer()
    err = producer.Send("amq.direct", "key.direct",
      []byte("producer.Send() | "+time.Now().Format("2006-01-02 15:04:05")),
      rabbitmq.DefaultSendOpts())
    if err != nil {
      log.Fatal(err)
    }
  }()

  // -- 接收消息 --
  go func() {
    consumer := conn.Consumer()
    consumer.Receive(
      "queue.direct",
      rabbitmq.NewReceiveOptsBuilder().SetAutoAck(false).Build(),
      &rabbitmq.AbsReceiveListener{
        ConsumerMethod: func(d *amqp.Delivery) (brk bool) {
          log.Println("queue.direct ", d.DeliveryTag, " ", string(d.Body))
          err := d.Ack(false)
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
  }()

  time.Sleep(time.Second * 10)
}
```

> 小建议：
> 
> 如无特殊需求，我们在全局只需创建一个 Connection，所有消息的发送和接收都使用 Producer 和 Consumer 处理。

Examples
---

消费

* [消费者极简用法](example/consumer/easy_consumer.go)
* [消费者常规用法](example/consumer/common_consumer.go)
* [断线重连，循环消费](example/consumer/reconnected_consumer.go)

生产

* [生产者极简用法](example/producer/easy_producer.go)
* [生产者常规用法](example/producer/common_producer.go)
* [断线重连，消息确认，失败重发](example/producer/re-send_producer.go)

Unresolved
---

* 消息异步发送与确认
* 重复消费
* 网络分区状态下，找不到队列时消息重发（ReturnListener）
* ~~消费者 ack 重试~~

License
---

LGPL