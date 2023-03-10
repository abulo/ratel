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

func Load(module constant.Module, key string) (any, bool) {
	return singleton.Load(genkey(module, key))
}

func Store(module constant.Module, key string, val any) {
	singleton.Store(genkey(module, key), val)
}
