package kafka

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type TopicBuilder func(string) string

func DefaultTopicBuilder(aggregateType string) string {
	return strings.ToLower(aggregateType)
}

func DefaultListenerConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Version = sarama.V2_6_0_0
	conf.Producer.Retry.Max = 10
	conf.Producer.Retry.Backoff = time.Second
	conf.Producer.MaxMessageBytes = 5000000
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	return conf
}

func DefaultClientConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Version = sarama.V2_6_0_0
	conf.Producer.Retry.Max = 10
	conf.Producer.Retry.Backoff = time.Second
	conf.Producer.MaxMessageBytes = 5000000
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	return conf
}
