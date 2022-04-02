package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/abulo/ratel/v2/config"
)

func main() {
	Config := config.New()
	Config.SetConfigName("config")
	// config.SetConfigName("mysql")
	Config.SetConfigType("toml")
	Config.AddConfigPath("/Users/abulo/WorkSpace/golang/src/ratel/example/env")
	// config.AddConfigPath("/Users/abulo/WorkSpace/golang/src/ratel/example/development")

	err := Config.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	configBytes, err := ioutil.ReadFile("/Users/abulo/WorkSpace/golang/src/ratel/example/development/mysql.toml")
	if err != nil {
		panic(fmt.Errorf("Could not read config file: %s \n", err))
	}
	Config.MergeConfig(bytes.NewBuffer(configBytes))

	cfg := Config.Get("configDir")

	fmt.Println(cfg)

	cfgmysql := Config.Get("mysql")

	fmt.Println(cfgmysql)

}
