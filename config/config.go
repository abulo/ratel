package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/viper"
)

// Config 在 Viper 之上封装配置读取能力，支持从文件加载配置并按类型读取配置项。
//
// 零值不可直接使用，须通过 New 构造；Config 内嵌 *viper.Viper，
// 可直接调用 Viper 原有方法。
type Config struct {
	*viper.Viper
}

// New 创建一个空的 Config 实例，尚未加载任何配置文件。
func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}

// LoadFile 加载指定路径的单个配置文件，文件格式由扩展名推断（json、toml、yaml 等）。
//
// filePath 为完整文件路径；文件不存在、缺少扩展名或读取失败时返回错误。
func (c *Config) LoadFile(filePath string) error {
	if !util.FileExists(filePath) {
		return fmt.Errorf("file is not exist: %s", filePath)
	}
	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file extension is required")
	}
	c.SetConfigType(strings.TrimPrefix(ext, "."))
	c.SetConfigFile(filePath)
	if err := c.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	return nil
}

// LoadDir 递归遍历目录，将其中 json、toml、yaml 配置文件依次合并到当前配置中。
//
// 目录遍历失败或任一配置文件合并失败时返回错误；非上述扩展名的文件会被跳过。
func (c *Config) LoadDir(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := make([]string, 0)
		ext = append(ext, "json", "toml", "yaml")

		if !info.IsDir() {
			fileExt := filepath.Ext(path)
			if fileExt == "" {
				return nil
			}
			if len(ext) > 0 && !util.InArray(strings.TrimPrefix(fileExt, "."), ext) {
				return nil
			}
			c.SetConfigType(strings.TrimPrefix(fileExt, "."))
			c.SetConfigFile(path)
			if err := c.MergeInConfig(); err != nil {
				return fmt.Errorf("failed to merge config file: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}
	return nil
}

// Exists 判断指定配置键是否已被设置。
func (c *Config) Exists(key string) bool {
	return c.IsSet(key)
}

// BindStruct 将指定键下的配置反序列化到 data 指向的结构体，data 必须为非 nil 指针。
func (c *Config) BindStruct(key string, data any) error {
	return c.UnmarshalKey(key, data)
}

// String 读取 key 对应的字符串配置；键不存在时返回可选默认值，未提供默认值则返回空字符串。
func (c *Config) String(key string, defVal ...string) string {
	if c.Exists(key) {
		return c.GetString(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return ""
}

// Int 读取 key 对应的整数配置；键不存在时返回可选默认值，未提供默认值则返回 0。
func (c *Config) Int(key string, defVal ...int) int {
	if c.Exists(key) {
		return c.GetInt(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Uint 读取 key 对应的无符号整数配置；键不存在时返回可选默认值，未提供默认值则返回 0。
func (c *Config) Uint(key string, defVal ...uint) uint {
	if c.Exists(key) {
		return c.GetUint(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Int64 读取 key 对应的 int64 配置；键不存在时返回可选默认值，未提供默认值则返回 0。
func (c *Config) Int64(key string, defVal ...int64) int64 {
	if c.Exists(key) {
		return c.GetInt64(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Duration 读取 key 对应的时长配置；键不存在时返回可选默认值，未提供默认值则返回 0。
func (c *Config) Duration(key string, defVal ...time.Duration) time.Duration {
	if c.Exists(key) {
		return c.GetDuration(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Float 读取 key 对应的浮点数配置；键不存在时返回可选默认值，未提供默认值则返回 0。
func (c *Config) Float(key string, defVal ...float64) float64 {
	if c.Exists(key) {
		return c.GetFloat64(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// Bool 读取 key 对应的布尔配置；键不存在时返回可选默认值，未提供默认值则返回 false。
func (c *Config) Bool(key string, defVal ...bool) bool {
	if c.Exists(key) {
		return c.GetBool(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return false
}
