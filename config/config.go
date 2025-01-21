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

type Config struct {
	*viper.Viper
}

func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}

// LoadFile 加载单个配置文件
func (c *Config) LoadFile(filePath string) error {
	// filePath 是完整路径
	if !util.FileExists(filePath) {
		return fmt.Errorf("file is not exist: %s", filePath)
	}
	// 获取文件扩展名
	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file extension is required")
	}
	// 设置配置文件类型
	c.SetConfigType(strings.TrimPrefix(ext, "."))
	// 设置配置文件路径
	c.SetConfigFile(filePath)
	// 读取配置文件
	if err := c.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	return nil
}

// LoadDir 加载目录中的所有配置文件
func (c *Config) LoadDir(dirPath string) error {
	// 遍历目录中的所有文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := make([]string, 0)
		ext = append(ext, "json", "toml", "yaml")

		// 只处理文件
		if !info.IsDir() {
			// 获取文件扩展名
			fileExt := filepath.Ext(path)
			if fileExt == "" {
				return nil
			}
			// 判断 fileExt 是否在 ext 中
			if len(ext) > 0 && !util.InArray(strings.TrimPrefix(fileExt, "."), ext) {
				return nil
			}
			// 设置配置文件类型
			c.SetConfigType(strings.TrimPrefix(fileExt, "."))

			// 设置配置文件路径
			c.SetConfigFile(path)

			// 读取配置文件
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

func (c *Config) Exists(key string) bool {
	return c.IsSet(key)
}

func (c *Config) BindStruct(key string, data any) error {
	return c.UnmarshalKey(key, data)
}

// 获取 String 类型配置, 如果不存在则返回默认值
func (c *Config) String(key string, defVal ...string) string {
	if c.Exists(key) {
		return c.GetString(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return ""
}

// 获取 Int 类型配置, 如果不存在则返回默认值
func (c *Config) Int(key string, defVal ...int) int {
	if c.Exists(key) {
		return c.GetInt(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// 获取 Uint 类型配置, 如果不存在则返回默认值
func (c *Config) Uint(key string, defVal ...uint) uint {
	if c.Exists(key) {
		return c.GetUint(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// 获取 Int64 类型配置, 如果不存在则返回默认值
func (c *Config) Int64(key string, defVal ...int64) int64 {
	if c.Exists(key) {
		return c.GetInt64(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// 获取  time.Duration 类型配置, 如果不存在则返回默认值
func (c *Config) Duration(key string, defVal ...time.Duration) time.Duration {
	if c.Exists(key) {
		return c.GetDuration(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// 获取 float64 类型配置, 如果不存在则返回默认值
func (c *Config) Float(key string, defVal ...float64) float64 {
	if c.Exists(key) {
		return c.GetFloat64(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

// 获取 bool 类型配置, 如果不存在则返回默认值
func (c *Config) Bool(key string, defVal ...bool) bool {
	if c.Exists(key) {
		return c.GetBool(key)
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return false
}
