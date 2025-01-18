package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		// 只处理文件
		if !info.IsDir() {
			// 获取文件扩展名
			ext := filepath.Ext(path)
			if ext == "" {
				return nil
			}

			// 设置配置文件类型
			c.SetConfigType(strings.TrimPrefix(ext, "."))

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

func (c *Config) String(key string) string {
	return c.GetString(key)
}
