// Package rabbitmq 提供基于 github.com/streadway/amqp 的 RabbitMQ 客户端封装。
//
// 该包提供可断线重连的 Connection、复用连接的 Channel、队列声明与绑定的 Builder、
// 消费者与生产者、收发选项 Builder、可配置的重试策略以及常用消息工厂；注册的消费操作会在
// 每次断线重连后自动重新执行。
package rabbitmq
