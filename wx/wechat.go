package wx

import (
	"sync"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	programConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	officialConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/openplatform"
	openConfig "github.com/silenceper/wechat/v2/openplatform/config"
)

// Client --
type Client struct {
	*wechat.Wechat
	miniProgram     map[string]*miniprogram.MiniProgram
	officialAccount map[string]*officialaccount.OfficialAccount
	openPlatform    map[string]*openplatform.OpenPlatform
	mu              sync.RWMutex
}

func NewClient() *Client {
	wc := wechat.NewWechat()
	newClient := &Client{
		miniProgram:     make(map[string]*miniprogram.MiniProgram),
		officialAccount: make(map[string]*officialaccount.OfficialAccount),
		openPlatform:    make(map[string]*openplatform.OpenPlatform),
	}
	newClient.Wechat = wc
	return newClient
}

// SetMiniProgram 实例化
func (c *Client) SetMiniProgram(group string, cfg *programConfig.Config) *Client {
	clientMiniProgram := c.Wechat.GetMiniProgram(cfg)
	c.mu.Lock()
	c.miniProgram[group] = clientMiniProgram
	c.mu.Unlock()
	return c
}

// SetOfficialAccount 实例化
func (c *Client) SetOfficialAccount(group string, cfg *officialConfig.Config) *Client {
	clientOfficialAccount := c.Wechat.GetOfficialAccount(cfg)
	c.mu.Lock()
	c.officialAccount[group] = clientOfficialAccount
	c.mu.Unlock()
	return c
}

// SetOpenPlatform 实例化
func (c *Client) SetOpenPlatform(group string, cfg *openConfig.Config) *Client {
	clientOpenPlatform := c.Wechat.GetOpenPlatform(cfg)
	c.mu.Lock()
	c.openPlatform[group] = clientOpenPlatform
	c.mu.Unlock()
	return c
}

// MiniProgram 获取实例
func (c *Client) MiniProgram(group string) *miniprogram.MiniProgram {
	c.mu.RLock()
	res := c.miniProgram[group]
	c.mu.RUnlock()
	return res
}

// OfficialAccount 获取实例
func (c *Client) OfficialAccount(group string) *officialaccount.OfficialAccount {
	c.mu.RLock()
	res := c.officialAccount[group]
	c.mu.RUnlock()
	return res
}

// OfficialAccount 获取实例
func (c *Client) OpenPlatform(group string) *openplatform.OpenPlatform {
	c.mu.RLock()
	res := c.openPlatform[group]
	c.mu.RUnlock()
	return res
}
