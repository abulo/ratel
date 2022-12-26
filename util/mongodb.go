package util

import (
	"encoding/json"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Pipeline gets aggregation pipeline from a string
func Pipeline(str string) mongo.Pipeline {
	var pipeline = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.M
		json.Unmarshal([]byte(str), &doc)
		var v bson.D
		b, _ := bson.Marshal(doc)
		bson.Unmarshal(b, &v)
		pipeline = append(pipeline, v)
	} else {
		var docs []bson.M
		json.Unmarshal([]byte(str), &docs)
		for _, doc := range docs {
			var v bson.D
			b, _ := bson.Marshal(doc)
			bson.Unmarshal(b, &v)
			pipeline = append(pipeline, v)
		}
	}
	return pipeline
}

// ConvertBson 构建 bson.D 查询数据
func ConvertBson(items bson.D, item bson.E) bson.D {
	key := make([]string, 0)
	for _, v := range items {
		key = append(key, v.Key)
	}
	if !InArray(item.Key, key) {
		items = append(items, item)
		return items
	}
	//再次赋值
	for i, newItem := range items {
		if newItem.Key == item.Key {
			tp := indirect(newItem.Value)
			switch tp.(type) {
			case bson.D:
				tran := items[i].Value.(bson.D)
				appendVal := item.Value.(bson.D)
				tran = append(tran, appendVal...)
				items[i].Value = tran
			default:
				items[i].Value = item.Value
			}
		}
	}
	return items
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
