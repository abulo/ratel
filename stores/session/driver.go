package session

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/stores/redis"
)

// Session 表示一份存放在 Redis 中的会话，整份会话以 JSON 形式保存在 Name 对应的键下，
// 支持以点分路径读写嵌套字段。
//
// 零值的 Driver 为 nil，须先注入底层 Redis 客户端；会话内容读写依赖 Driver，
// 同一 Session 的并发写存在覆盖风险，调用方需自行保证访问时序。
type Session struct {
	// Driver 为底层 Redis 客户端，会话键均经由它读写，不可为 nil。
	Driver *redis.Client
	// Name 为会话在 Redis 中的键名。
	Name string
	// TTL 为会话过期时长，单位秒；0 表示不过期。
	TTL int64
}

// Put 将会话值写入指定键：键名含点号时按点分路径写入嵌套层级，否则作为顶层键写入。
//
// 写入时先读取现有会话 JSON，解析失败会返回错误；序列化后连同 TTL 回写会话键，
// Redis Set 返回的错误被忽略。ctx 为 nil 或已取消时替换为 context.TODO()。
func (session *Session) Put(ctx context.Context, key string, value any) error {
	var h = session.Driver
	var bytes []byte
	var err error
	var content string
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	m := make(map[string]any)

	if val, err := h.Exists(ctx, session.Name); err == nil && val {
		content, _ = h.Get(ctx, session.Name)
	} else {
		content = "{}"
	}

	err = json.Unmarshal([]byte(content), &m)
	if err != nil {
		return err
	}

	var keys = strings.Split(key, ".")
	var depth = len(keys)

	if depth < 2 {
		m[key] = value
	} else {
		m = setSliceMap(m, keys, value)
	}

	bytes, _ = json.Marshal(m)
	h.Set(ctx, session.Name, string(bytes), time.Duration(session.TTL)*time.Second)
	return nil
}

// Get 读取指定键的会话值：键名含点号时沿点分路径取嵌套层级，否则取顶层键；
// 键不存在或会话内容损坏时返回 nil。ctx 为 nil 或已取消时替换为 context.TODO()。
func (session *Session) Get(ctx context.Context, key string) any {
	var h = session.Driver
	var m map[string]any
	var content string

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	content, _ = h.Get(ctx, session.Name)
	_ = json.Unmarshal([]byte(content), &m)

	var keys = strings.Split(key, ".")
	var n = len(keys)
	if n < 2 {
		return m[key]
	}
	return getSliceMap(m, keys)
}

// Remove 删除指定键的会话值：键名含点号时删除路径末端的嵌套字段，否则删除顶层键，
// 删除结果连同 TTL 回写会话键。ctx 为 nil 或已取消时替换为 context.TODO()。
func (session *Session) Remove(ctx context.Context, key string) {
	var h = session.Driver
	var m map[string]any
	var content string
	var bytes []byte

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	content, _ = h.Get(ctx, session.Name)
	_ = json.Unmarshal([]byte(content), &m)

	var keys = strings.Split(key, ".")
	var n = len(keys)
	if n < 2 {
		delete(m, key)
		bytes, _ = json.Marshal(m)
		h.Set(ctx, session.Name, string(bytes), time.Duration(session.TTL)*time.Second)
		return
	}
	delSliceMap(m, keys)
	bytes, _ = json.Marshal(m)
	h.Set(ctx, session.Name, string(bytes), time.Duration(session.TTL)*time.Second)
}

// setSliceMap 沿点分路径 keys 在嵌套 map 中逐层补齐中间层级，并在末端键写入 value；
// 路径中段存在非 map 类型值时会触发类型断言 panic。
func setSliceMap(m map[string]any, keys []string, value any) map[string]any {
	var itMap = m
	var i int
	var limit = len(keys) - 1

	for i = range limit {
		_, ok := itMap[keys[i]]
		if !ok {
			itMap[keys[i]] = make(map[string]any)
		}
		itMap = itMap[keys[i]].(map[string]any)
	}
	itMap[keys[limit]] = value

	return m
}

// getSliceMap 沿点分路径 keys 逐级查找嵌套值；路径任一层级缺失或类型不符时返回 nil。
func getSliceMap(m map[string]any, keys []string) any {
	var itMap = m
	var i int
	var limit = len(keys) - 1
	var v any
	var ok bool

	for i = 0; i < limit; i++ {
		v, ok = itMap[keys[i]]
		if !ok {
			break
		}
		itMap = v.(map[string]any)
	}
	v, ok = itMap[keys[i]]
	if !ok {
		return nil
	}
	return v
}

// delSliceMap 沿点分路径 keys 逐级下钻并删除末端字段；路径中段缺失时直接返回。
func delSliceMap(m map[string]any, keys []string) {
	var itMap = m
	var i int
	var limit = len(keys) - 1
	var v any
	var ok bool

	for i = 0; i < limit; i++ {
		v, ok = itMap[keys[i]]
		if !ok {
			break
		}
		itMap = v.(map[string]any)
	}
	_, ok = itMap[keys[i]]
	if !ok {
		return
	}
	delete(itMap, keys[i])
}

// Destroy 删除整份会话键，返回被删除的键数量；Del 返回的错误被忽略。
func (session *Session) Destroy(ctx context.Context) int64 {
	var h = session.Driver

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	val, _ := h.Del(ctx, session.Name)
	return val
}
