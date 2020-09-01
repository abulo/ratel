package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/abulo/ratel/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

// Client --
type Client struct {
	esClient *elasticsearch.Client
}

// NewClient --
func NewClient(config elasticsearch.Config) *Client {
	esClient, err := elasticsearch.NewClient(config)
	if err != nil {
		logger.Fatal(err)
	}
	return &Client{esClient: esClient}
}

// Index --
func (client *Client) Index(ctx context.Context, index, id string, doc interface{}) (*IndexResponse, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := json.NewEncoder(buf).Encode(doc); err != nil {
		return nil, err
	}
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       buf,
		Refresh:    "true",
	}
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "IndexRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.String("DocumentID", id))
			span.LogFields(log.Object("Body", doc))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	indexResponse := &IndexResponse{}
	if err := json.NewDecoder(res.Body).Decode(indexResponse); err != nil {
		return nil, err
	}
	indexResponse.HTTPStatusCode = res.StatusCode
	return indexResponse, nil
}

// Create --
func (client *Client) Create(ctx context.Context, index, id string, doc interface{}) (*IndexResponse, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := json.NewEncoder(buf).Encode(doc); err != nil {
		return nil, err
	}

	req := esapi.CreateRequest{
		Index:      index,
		DocumentID: id,
		Body:       buf,
		Refresh:    "true",
	}

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "CreateRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.String("DocumentID", id))
			span.LogFields(log.Object("Body", doc))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	indexResponse := &IndexResponse{}
	if err := json.NewDecoder(res.Body).Decode(indexResponse); err != nil {
		return nil, err
	}
	indexResponse.HTTPStatusCode = res.StatusCode
	return indexResponse, nil
}

// Update --
func (client *Client) Update(ctx context.Context, index, id string, doc interface{}) (*IndexResponse, error) {
	query := make(map[string]interface{})
	query["doc"] = doc
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       buf,
		Refresh:    "true",
	}

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "UpdateRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.String("DocumentID", id))
			span.LogFields(log.Object("Body", doc))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	indexResponse := &IndexResponse{}
	if err := json.NewDecoder(res.Body).Decode(indexResponse); err != nil {
		return nil, err
	}
	indexResponse.HTTPStatusCode = res.StatusCode
	return indexResponse, nil
}

// Get --
func (client *Client) Get(ctx context.Context, index, id string, response interface{}) error {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "GetRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.String("DocumentID", id))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return err
	}
	return nil
}

// Delete --
func (client *Client) Delete(ctx context.Context, index, id string) (*IndexResponse, error) {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "DeleteRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.String("DocumentID", id))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	indexResponse := &IndexResponse{}
	if err := json.NewDecoder(res.Body).Decode(indexResponse); err != nil {
		return nil, err
	}
	indexResponse.HTTPStatusCode = res.StatusCode
	return indexResponse, nil
}

// Search --
func (client *Client) Search(ctx context.Context, index string, query interface{}, response interface{}) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return err
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  buf,
	}

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("elasticsearch", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "elasticsearch")
			span.SetTag("method", "SearchRequest")
			span.LogFields(log.String("Index", index))
			span.LogFields(log.Object("Body", query))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	res, err := req.Do(ctx, client.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return err
	}
	return nil
}
