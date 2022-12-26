package session

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/abulo/ratel/v2/store/redis"
)

// Session ...
type Session struct {
	Driver *redis.Client
	Name   string
	TTL    int64 // seconds
}

// Put https://laravel.com/docs/5.8/session
// Pushing To Array Session Values
// session.Put('user.teams', 'developers');  => {user: {teams: "developer"}}
func (session *Session) Put(ctx context.Context, key string, value interface{}) error {
	var h = session.Driver
	var bytes []byte
	var err error
	var content string
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	m := make(map[string]interface{})

	// fmt.Printf("Exists %s?\n", session.Name)
	if h.Exists(ctx, session.Name).Val() == 1 {
		content = h.Get(ctx, session.Name).Val()
	} else {
		content = "{}"
	}
	// fmt.Printf("%s => %s", session.Name, content)

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

// Get s.Get
func (session *Session) Get(ctx context.Context, key string) interface{} {
	var h = session.Driver
	var m map[string]interface{}
	var content string

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	content = h.Get(ctx, session.Name).Val()
	_ = json.Unmarshal([]byte(content), &m)

	var keys = strings.Split(key, ".")
	var n = len(keys)
	if n < 2 {
		return m[key]
	}
	return getSliceMap(m, keys)
}

// Remove ...
func (session *Session) Remove(ctx context.Context, key string) {
	var h = session.Driver
	var m map[string]interface{}
	var content string
	var bytes []byte

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	content = h.Get(ctx, session.Name).Val()
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

// sess.Put("info.age", "28")
func setSliceMap(m map[string]interface{}, keys []string, value interface{}) map[string]interface{} {
	var itMap = m
	var i int
	var limit = len(keys) - 1

	for i = 0; i < limit; i++ {
		_, ok := itMap[keys[i]]
		if !ok {
			// } else {
			itMap[keys[i]] = make(map[string]interface{})
		}
		itMap = itMap[keys[i]].(map[string]interface{})
	}
	itMap[keys[limit]] = value

	return m
}

func getSliceMap(m map[string]interface{}, keys []string) interface{} {
	var itMap = m
	var i int
	var limit = len(keys) - 1
	var v interface{}
	var ok bool

	for i = 0; i < limit; i++ {
		v, ok = itMap[keys[i]]
		if !ok {
			break
		}
		itMap = v.(map[string]interface{})
	}
	v, ok = itMap[keys[i]]
	if !ok {
		return nil
	}
	return v
}

func delSliceMap(m map[string]interface{}, keys []string) {
	var itMap = m
	var i int
	var limit = len(keys) - 1
	var v interface{}
	var ok bool

	for i = 0; i < limit; i++ {
		v, ok = itMap[keys[i]]
		if !ok {
			break
		}
		itMap = v.(map[string]interface{})
	}
	_, ok = itMap[keys[i]]
	if !ok {
		return
	}
	delete(itMap, keys[i])
}

/*
// TO BE DONE
func (session *Session) Has(key string) bool {
	var h = session.getDriver()
	return h.Exists(session.Name).Val() == 1
}
// TO BE DONE
func (session *Session) Push(key string, e interface{}) {
}
*/

// Destroy 释放
func (session *Session) Destroy(ctx context.Context) int64 {
	var h = session.Driver

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return h.Del(ctx, session.Name).Val()
}

// func TestSession_Put(t *testing.T) {
// 	var sess = Session{
// 		Name:   "sess_1",
// 		Driver: client,
// 		TTL:    0,
// 	}
// 	var userID int64
// 	userID = 2147483647
// 	sess.Put("id", strconv.FormatInt(userID, 10))
// 	// sess.Put("id", userID)
// 	sess.Put("name", "mingzhanghui")
// 	sess.Put("info.age", "29")
// 	sess.Put("info.gender", "male")
// 	sess.Put("a.b.c.d", "1")
// }
// func TestSession_Get(t *testing.T) {
// 	var sess = Session{
// 		Name:   "sess_1",
// 		Driver: client,
// 	}
// 	userId := sess.Get("id")
// 	// t.Logf("%d\n", userId.(int64))
// 	assert.Equal(t, userId, "2147483647")

// 	gender := sess.Get("info.gender")
// 	t.Logf(gender.(string))
// 	assert.Equal(t, gender, "male")

// 	foo := sess.Get("a.b.c.d")
// 	t.Logf(foo.(string))
// 	assert.Equal(t, foo, "1")

// }
