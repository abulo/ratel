package initial

import (
	"time"

	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/abulo/ratel/v3/util"
)

// Initial ...
type Initial struct {
	Path       string         // 应用程序执行路径
	Config     *config.Config // 配置文件
	Store      *proxy.Proxy   // 数据库链接
	LaunchTime time.Time      // 时间设置
}

// Core 系统
var Core *Initial

// New Default returns an Initial instance.
func New() *Initial {
	engine := NewInitial()
	return engine
}

// NewInitial ...
func NewInitial() *Initial {
	Core = &Initial{
		Store: proxy.NewProxy(),
	}
	Core.InitPath(util.GetAppRootPath())
	Core.InitLaunchTime(util.Now())
	return Core
}
