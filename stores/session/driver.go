package session

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/stores/redis"
	rdb "github.com/redis/go-redis/v9"
)

// sessionBytesPool 复用会话 JSON 解码所用的 []byte，避免每次请求都产生一次 string→[]byte 拷贝。
// PERF-1-11 [P2]（预期收益需实测验证）。
var sessionBytesPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

// acquireSessionBytes 从池中取出一个切片并写入会话内容；使用完毕后须调用 releaseSessionBytes 归还。
func acquireSessionBytes(content string) []byte {
	bp := sessionBytesPool.Get().(*[]byte)
	b := (*bp)[:0]
	b = append(b, content...)
	*bp = b
	return b
}

// releaseSessionBytes 将解码切片归还对象池；容量过大的切片不回池，避免长期驻留大块内存。
func releaseSessionBytes(b []byte) {
	if cap(b) > 64*1024 {
		return
	}
	sessionBytesPool.Put(&b)
}

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

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	// PERF-1-10 [P1] 去掉冗余的 Exists 往返：单次 Get 即可判断键是否存在，
	// Get 返回 redis.Nil 或空内容时按空会话处理，原 Exists 判定不存在的语义保持等价。
	m, err := session.decodeSessionMap(ctx)
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

// decodeSessionMap 单次 Get 读取会话 JSON 并反序列化为 map；键不存在或读取失败时按空会话处理。
// PERF-1-12 [P2] 合并 Put/Get/Remove 重复的「读取会话 → 解码」样板逻辑；
// 返回的解码错误由调用方按原有行为决定是否上抛。
func (session *Session) decodeSessionMap(ctx context.Context) (map[string]any, error) {
	var h = session.Driver
	var m = make(map[string]any)

	content, gErr := h.Get(ctx, session.Name)
	switch {
	case errors.Is(gErr, rdb.Nil) || content == "":
		// 键不存在：与原 Exists 判定不存在等价，按空会话处理。
		content = "{}"
	case gErr != nil:
		// 保留原行为：非 Nil 的读取错误被忽略，content 保持空串，交由 Unmarshal 报错。
		content = ""
	}

	buf := acquireSessionBytes(content)
	defer releaseSessionBytes(buf)
	err := json.Unmarshal(buf, &m)
	return m, err
}

// Get 读取指定键的会话值：键名含点号时沿点分路径取嵌套层级，否则取顶层键；
// 键不存在或会话内容损坏时返回 nil。ctx 为 nil 或已取消时替换为 context.TODO()。
func (session *Session) Get(ctx context.Context, key string) any {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	// PERF-1-11 [P2] 解码所用 []byte 由对象池复用，去掉每次读取的 string→[]byte 拷贝分配。
	// PERF-1-12 [P2] 复用与 Put/Remove 统一的会话读取解码逻辑。
	m, _ := session.decodeSessionMap(ctx)

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
	var bytes []byte

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	// PERF-1-12 [P2] 顶层与嵌套分支合并为统一的读-改-写流程：先复用解码逻辑读取会话，
	// 删除字段后统一序列化回写，消除两个分支重复的 Get/Unmarshal/Marshal/Set 样板。
	m, _ := session.decodeSessionMap(ctx)

	var keys = strings.Split(key, ".")
	var n = len(keys)
	if n < 2 {
		delete(m, key)
	} else {
		delSliceMap(m, keys)
	}

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
