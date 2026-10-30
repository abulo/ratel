package util

import (
	"encoding/json"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Pipeline 将 JSON 字符串解析为 MongoDB 聚合管道（mongo.Pipeline）。
//
// 字符串以 [ 开头时按数组解析为多个聚合阶段，否则按单个文档解析；
// JSON 解析失败的错误会被忽略，调用方需保证传入合法的 JSON。
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

// ConvertBson 将 bson.E 元素追加到 bson.D 文档中，同键元素存在时进行合并。
//
// 同键元素已存在且值为 bson.D 时，将新值按子文档追加合并，否则用新值覆盖；
// 调用方须保证该场景下新旧值均为 bson.D，否则会触发类型断言 panic。
func ConvertBson(items bson.D, item bson.E) bson.D {
	key := make([]string, 0)
	for _, v := range items {
		key = append(key, v.Key)
	}
	if !InArray(item.Key, key) {
		items = append(items, item)
		return items
	}
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
//
// indirect 反复解引用指针直到取得基础类型的值；输入为 nil 或非指针时原样返回。
func indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		// 避免为普通值创建 reflect.Value，减少不必要的开销。
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
