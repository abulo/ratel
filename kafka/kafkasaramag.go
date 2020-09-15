package kafka

import (
	"context"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/abulo/ratel/logger"
)

type saramaConsumerGroupHandler struct {
	*ConsumerConfig
	ConsumerMessageHandler
	context.Context
}

func (h *saramaConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	if h.ConsumerConfig.Offset != 0 {
		// 如果是OffsetOldest则不变, 否则自动前移, 因为位置从0开始.
		var realOffset int64 = h.ConsumerConfig.Offset
		// 修正CURRENT-OFFSET
		if realOffset > 0 {
			realOffset--
		} else if realOffset < 0 {
			realOffset = 0
		}
		for t, ps := range s.Claims() {
			for _, p := range ps {
				s.ResetOffset(t, p, realOffset, "")
			}
		}
	}
	return nil
}
func (h *saramaConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }
func (h *saramaConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) (err error) {
	// 使用context控制message与error的handler的退出时机
	for {
		select {
		case <-h.Context.Done():
			return
		case msg, ok := <-c.Messages():
			// Consume joins a cluster of consumers for a given list of topics and
			// starts a blocking ConsumerGroupSession through the ConsumerGroupHandler.
			//
			// The life-cycle of a session is represented by the following steps:
			//
			// 1. The consumers join the group (as explained in https://kafka.apache.org/documentation/#intro_consumers)
			//    and is assigned their "fair share" of partitions, aka 'claims'.
			// 2. Before processing starts, the handler's Setup() hook is called to notify the user
			//    of the claims and allow any necessary preparation or alteration of state.
			// 3. For each of the assigned claims the handler's ConsumeClaim() function is then called
			//    in a separate goroutine which requires it to be thread-safe. Any state must be carefully protected
			//    from concurrent reads/writes.
			// 4. The session will persist until one of the ConsumeClaim() functions exits. This can be either when the
			//    parent context is cancelled or when a server-side rebalance cycle is initiated.
			// 5. Once all the ConsumeClaim() loops have exited, the handler's Cleanup() hook is called
			//    to allow the user to perform any final tasks before a rebalance.
			// 6. Finally, marked offsets are committed one last time before claims are released.
			//
			// Please note, that once a rebalance is triggered, sessions must be completed within
			// Config.Consumer.Group.Rebalance.Timeout. This means that ConsumeClaim() functions must exit
			// as quickly as possible to allow time for Cleanup() and the final offset commit. If the timeout
			// is exceeded, the consumer will be removed from the group by Kafka, which will cause offset
			// commit failures.
			// This method should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims.
			if !ok {
				return
			}
			switch h.ConsumerConfig.Ack {
			case ACK_BEFORE_AUTO:
				s.MarkMessage(msg, "")
				err = h.ConsumerMessageHandler(msg)
			case ACK_AFTER_NOERROR:
				if err = h.ConsumerMessageHandler(msg); err == nil {
					s.MarkMessage(msg, "")
				}
			case ACK_AFTER_NOMATTER:
				err = h.ConsumerMessageHandler(msg)
				s.MarkMessage(msg, "")
			default:
				logger.Logger.Panic("invalid ack type: " + strconv.Itoa(h.ConsumerConfig.Ack))
			}
		}
	}

}

type saramaConsumerGroup struct {
	*ConsumerConfig
	sarama.ConsumerGroup
	context.CancelFunc
}

func newSaramaConsumerGroup(c *ConsumerConfig) (ret *saramaConsumerGroup, err error) {
	grp, err := sarama.NewConsumerGroup(c.Address, c.Group, consumerConfig(c))
	if err != nil {
		return
	}
	ret = &saramaConsumerGroup{
		ConsumerConfig: c,
		ConsumerGroup:  grp,
	}
	return
}

func (g *saramaConsumerGroup) Close() (err error) {
	if g.CancelFunc != nil {
		g.CancelFunc()
	}
	if g.ConsumerGroup != nil {
		err = g.ConsumerGroup.Close()
	}
	return
}

// 必须保证ConsumerMessageHandler, ConsumerErrorHandler没有panic
func (g *saramaConsumerGroup) Consume(topic string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) (err error) {
	return g.ConsumeM([]string{topic}, mh, eh)
}

// 必须保证ConsumerMessageHandler, ConsumerErrorHandler没有panic
func (g *saramaConsumerGroup) ConsumeM(topics []string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) (err error) {
	// 先关闭旧的消费过程
	if g.CancelFunc != nil {
		g.CancelFunc()
	}

	// 每次进入都会新起context. 因为ConsumeM()不能重复调用,否则会停掉之前的工作
	var ctx context.Context
	ctx, g.CancelFunc = context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-g.Errors():
				eh(e)
			}
		}
	}(ctx)
	for {
		err = g.ConsumerGroup.Consume(ctx, topics, &saramaConsumerGroupHandler{
			ConsumerConfig:         g.ConsumerConfig,
			ConsumerMessageHandler: mh,
			Context:                ctx,
		})
		if err != nil {
			return
		} else if err = ctx.Err(); err != nil {
			return
		}
	}
	return
}
