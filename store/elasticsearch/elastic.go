package elasticsearch

import (
	"github.com/abulo/ratel/v2/logger"
	"github.com/olivere/elastic/v7"
)

// Client --
type Client struct {
	*elastic.Client
}

// NewClient --
func NewClient(options ...elastic.ClientOptionFunc) *Client {
	client, err := elastic.NewClient(options...)
	if err != nil {
		logger.Logger.Panic(err)
	}
	newClient := &Client{}
	newClient.Client = client
	return newClient
}
