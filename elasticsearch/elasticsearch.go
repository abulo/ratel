package elasticsearch

import (
	"github.com/abulo/ratel/logger"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

// Client --
type Client struct {
	*elastic.Client
}

// NewClient --
func NewClient(options *config.Config) *Client {
	client, err := elastic.NewClientFromConfig(options)
	if err != nil {
		logger.Logger.Panic(err)
	}
	newClient := &Client{}
	newClient.Client = client
	return newClient
}
