package rabbitmq

import (
	"sync"

	"github.com/streadway/amqp"
)

// Consumer 绑定所属连接，提供持续订阅与单条拉取两种消费方式。
type Consumer struct {
	c *Connection
}

// Receive 订阅指定队列并持续消费，订阅操作会在每次断线重连后自动重新注册执行。
//
// 方法内部通过后台 goroutine 运行接收循环，没有消息时不会阻塞调用方；消费细节与 autoAck 语义
// 见 Channel.ReceiveOpts。消费结束（正常退出或出错）时会回调监听器的 Finish，成功退出时先调用
// Remove 注销订阅操作。
func (c *Consumer) Receive(queue string, opts *ReceiveOpts, lis ReceiveListener) {
	c.c.RegisterAndExec(func(key string, ch *Channel) {
		err := ch.ReceiveOpts(queue, lis.Consumer, opts)
		if err == nil {
			lis.Remove(key, ch)
		}
		lis.Finish(err)
	})
}

// Get 拉取队列中的单条消息；返回消息体、是否取到消息以及错误。
//
// autoAck 为 true 时服务器自动确认该消息，若通道或连接在消息处理完成前关闭，消息不会重新入队；
// autoAck 为 false 时由调用方手动确认。
func (c *Consumer) Get(queue string, autoAck bool) (*amqp.Delivery, bool, error) {
	// PERF-5-13 [P1]：拉取消费复用同一条互斥保护的长驻通道，不再每条消息 open/close；
	// autoAck 参数与错误返回语义保持不变，通道异常时槽位自动清空并于下次调用重建。
	ch, err := c.c.trackChannel(&c.c.getMu, &c.c.getCh)
	if err != nil {
		return nil, false, err
	}
	msg, ok, err := ch.Get(queue, autoAck)
	if err != nil {
		c.c.dropChannelOnErr(&c.c.getMu, &c.c.getCh, ch, err)
		return nil, false, err
	}
	return &msg, ok, err
}

// Producer 绑定所属连接，提供消息发布入口。
type Producer struct {
	c *Connection
}

// Send 向指定交换器与路由键发布消息；普通发布与确认模式分别复用连接上的常驻通道。
//
// opts 为 nil 时使用默认发送选项；配置 retryable 后进入同步发布确认流程，确认前阻塞后续发送，
// 会显著降低 QPS，详见 Channel.SendOpts。
func (p *Producer) Send(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	// PERF-5-13 [P1]：复用连接上互斥保护的长驻通道，消除每条消息 channel.open/close 两次 broker
	// 往返；普通发布与确认模式分槽复用，避免共用一条通道时 Confirm 状态与确认监听串扰；错误语义
	// 保持原样透传，可恢复错误后槽位自动重建。（预期收益需实测验证）
	resolved := opts
	if resolved == nil {
		resolved = DefaultSendOpts()
	}

	var mu *sync.Mutex
	var slot **Channel
	if resolved.retryable == nil {
		mu, slot = &p.c.pubMu, &p.c.pubCh
	} else {
		mu, slot = &p.c.confirmMu, &p.c.confirmCh
	}

	ch, err := p.c.trackChannel(mu, slot)
	if err != nil {
		return err
	}
	if err = ch.SendOpts(exchange, routingKey, body, opts); err != nil {
		p.c.dropChannelOnErr(mu, slot, ch, err)
		return err
	}
	return nil
}
