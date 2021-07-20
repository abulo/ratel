package wx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abulo/ratel/redis"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	programConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	officialConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

// Client --
type Client struct {
	*wechat.Wechat
	MiniProgram     *miniprogram.MiniProgram
	OfficialAccount *officialaccount.OfficialAccount
}

func NewClient() *Client {
	wc := wechat.NewWechat()
	newClient := &Client{}
	newClient.Wechat = wc
	return newClient
}

// NewMiniProgram 实例化小程序
func (c *Client) NewMiniProgram(cfg *programConfig.Config) *Client {
	clientMiniProgram := c.Wechat.GetMiniProgram(cfg)
	c.MiniProgram = clientMiniProgram
	return c
}

func (c *Client) NewOfficialAccount(cfg *officialConfig.Config) *Client {
	clientOfficialAccount := c.Wechat.GetOfficialAccount(cfg)
	c.OfficialAccount = clientOfficialAccount
	return c
}

type Cache struct {
	Driver *redis.Client
}

//Get 获取一个值
func (c *Cache) Get(key string) interface{} {
	ctx := context.TODO()
	var data interface{}
	content := c.Driver.Get(ctx, key).Val()
	json.Unmarshal([]byte(content), &data)
	return data
}

//Set 设置一个值
func (c *Cache) Set(key string, val interface{}, timeout time.Duration) (err error) {
	ctx := context.TODO()
	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}
	_, err = c.Driver.SetEX(ctx, key, string(data), timeout).Result()
	return
}

//IsExist 判断key是否存在
func (c *Cache) IsExist(key string) bool {
	ctx := context.TODO()
	a, _ := c.Driver.Exists(ctx, key).Result()
	return a > 0
}

//Delete 删除
func (c *Cache) Delete(key string) error {
	ctx := context.TODO()
	if _, err := c.Driver.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
