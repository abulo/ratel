package singleton

import (
	"strconv"
	"sync"

	"github.com/abulo/ratel/v3/core/constant"
)

var singleton sync.Map

// PERF-4-16 [P2]：模块编号为整数，直接 strconv.Itoa 拼接，去掉 cast.ToString 的类型分支反射。
func genkey(module constant.Module, key string) string {
	return strconv.Itoa(int(module)) + key
}

// Load 读取指定模块与键下保存的单例；不存在时返回 nil 与 false。
func Load(module constant.Module, key string) (any, bool) {
	return singleton.Load(genkey(module, key))
}

// Store 将 val 保存到指定模块与键下；同模块同键重复调用会覆盖原有实例。
func Store(module constant.Module, key string, val any) {
	singleton.Store(genkey(module, key), val)
}
