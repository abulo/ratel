package singleton

import (
	"sync"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/spf13/cast"
)

var singleton sync.Map

func genkey(module constant.Module, key string) string {
	return cast.ToString(int(module)) + key
}

// Load 读取指定模块与键下保存的单例；不存在时返回 nil 与 false。
func Load(module constant.Module, key string) (any, bool) {
	return singleton.Load(genkey(module, key))
}

// Store 将 val 保存到指定模块与键下；同模块同键重复调用会覆盖原有实例。
func Store(module constant.Module, key string, val any) {
	singleton.Store(genkey(module, key), val)
}
