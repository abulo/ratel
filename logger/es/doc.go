package es

import (
	"context"

	"github.com/abulo/ratel/v3/store/elasticsearch"
	"github.com/pkg/errors"
)

func IndexString() string {
	body := `{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
			}
		},
		"mappings": {
			"dynamic": false,
			"properties": {
				"file": {
					"type": "keyword"
				},
				"func": {
					"type": "keyword"
				},
				"message": {
					"type": "keyword"
				},
				"level": {
					"type": "keyword"
				},
				"data": {
					"type": "object"
				},
				"host": {
					"type": "keyword"
				},
				"timestamp": {
					"type": "date"
				}
			}
		}
	}`
	return body
}

func CreateIndex(client *elasticsearch.Client) error {
	ctx := context.Background()
	//check es index
	exists, err := client.IndexExists("logger").Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("logger").BodyJson(IndexString()).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Not acknowledged")
		}
	}
	return nil
}
