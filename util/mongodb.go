package util

import (
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
		// PERF-7-22 [P2] 直接由扩展 JSON 解析为 bson.D，跳过 json→bson.M→bson.Marshal→bson.Unmarshal 的三重序列化往返（是否处于查询热路径需实测验证）
		var v bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &v)
		pipeline = append(pipeline, v)
	} else {
		var docs []bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &docs)
		pipeline = append(pipeline, docs...)
	}
	return pipeline
}

// ConvertBson 将 bson.E 元素追加到 bson.D 文档中，同键元素存在时进行合并。
//
// 同键元素已存在且值为 bson.D 时，将新值按子文档追加合并，否则用新值覆盖；
// 调用方须保证该场景下新旧值均为 bson.D，否则会触发类型断言 panic。
func ConvertBson(items bson.D, item bson.E) bson.D {
	// PERF-7-23 [P2] 入口一次遍历建立 键→下标 映射，单趟完成判重与定位，去掉先收集键切片再走反射型 InArray、随后二次全量循环定位的两次遍历
	indexOf := make(map[string]int, len(items))
	for i, v := range items {
		indexOf[v.Key] = i
	}
	if idx, ok := indexOf[item.Key]; !ok {
		items = append(items, item)
		return items
	} else {
		tp := indirect(items[idx].Value)
		switch tp.(type) {
		case bson.D:
			tran := items[idx].Value.(bson.D)
			appendVal := item.Value.(bson.D)
			tran = append(tran, appendVal...)
			items[idx].Value = tran
		default:
			items[idx].Value = item.Value
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
