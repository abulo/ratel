/*
Package config is a go config management implement. support YAML,TOML,JSON,INI,HCL format.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/config

JSON format content example:

	{
		"name": "app",
		"debug": false,
		"baseKey": "value",
		"age": 123,
		"envKey": "${SHELL}",
		"envKey1": "${NotExist|defValue}",
		"map1": {
			"key": "val",
			"key1": "val1",
			"key2": "val2"
		},
		"arr1": [
			"val",
			"val1",
			"val2"
		],
		"lang": {
			"dir": "res/lang",
			"defLang": "en",
			"allowed": {
				"en": "val",
				"zh-CN": "val2"
			}
		}
	}

Usage please see example(more example please see examples folder in the lib):
*/
package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Ini There are supported config format
const (
	Ini  = "ini"
	Hcl  = "hcl"
	Yml  = "yml"
	JSON = "json"
	Yaml = "yaml"
	Toml = "toml"

	// default delimiter
	defaultDelimiter byte = '.'
	// default struct tag name for binding data to struct
	defaultStructTag = "mapstructure"
)

// internal vars
// type intArr []int
type strArr []string

// type intMap map[string]int
type strMap map[string]string

// This is a default config manager instance
var dc = New("default")

// Config structure definition
type Config struct {
	// save latest error, will clear after read.
	err error
	// config instance name
	name string
	lock sync.RWMutex

	// config options
	opts *Options
	// all config data
	data map[string]any

	// loaded config files records
	loadedFiles []string
	driverNames []string

	// TODO Deprecated decoder and encoder, use driver instead
	// drivers map[string]Driver

	// decoders["toml"] = func(blob []byte, v any) (err error){}
	// decoders["yaml"] = func(blob []byte, v any) (err error){}
	decoders map[string]Decoder
	encoders map[string]Encoder

	// cache on got config data
	intCache map[string]int
	strCache map[string]string
	// iArrCache map[string]intArr TODO cache it
	// iMapCache map[string]intMap
	sArrCache map[string]strArr
	sMapCache map[string]strMap

	onConfigChange func(fsnotify.Event)
}

// New config instance
func New(name string) *Config {
	return &Config{
		name: name,
		opts: newDefaultOption(),
		data: make(map[string]any),

		// default add JSON driver
		encoders: map[string]Encoder{JSON: JSONEncoder},
		decoders: map[string]Decoder{JSON: JSONDecoder},
	}
}

// NewEmpty config instance
func NewEmpty(name string) *Config {
	return &Config{
		name: name,
		opts: newDefaultOption(),
		data: make(map[string]any),

		// don't add any drivers
		encoders: map[string]Encoder{},
		decoders: map[string]Decoder{},
	}
}

// NewWith create config instance, and you can call some init func
func NewWith(name string, fn func(c *Config)) *Config {
	return New(name).With(fn)
}

// NewWithOptions config instance
func NewWithOptions(name string, opts ...func(*Options)) *Config {
	return New(name).WithOptions(opts...)
}

// Default get the default instance
func Default() *Config {
	return dc
}

/*************************************************************
 * config drivers
 *************************************************************/

// AddDriver set a decoder and encoder driver for a format.
func AddDriver(driver Driver) { dc.AddDriver(driver) }

// AddDriver set a decoder and encoder driver for a format.
func (c *Config) AddDriver(driver Driver) {
	format := driver.Name()

	c.driverNames = append(c.driverNames, format)
	c.decoders[format] = driver.GetDecoder()
	c.encoders[format] = driver.GetEncoder()
}

// HasDecoder has decoder
func (c *Config) HasDecoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.decoders[format]
	return ok
}

// HasEncoder has encoder
func (c *Config) HasEncoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.encoders[format]
	return ok
}

// DelDriver delete driver of the format
func (c *Config) DelDriver(format string) {
	format = fixFormat(format)
	delete(c.decoders, format)
	delete(c.encoders, format)
}

/*************************************************************
 * helper methods
 *************************************************************/

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Error get last error, will clear after read.
func (c *Config) Error() error {
	err := c.err
	c.err = nil
	return err
}

// IsEmpty of the config
func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}

// LoadedFiles get loaded files name
func (c *Config) LoadedFiles() []string {
	return c.loadedFiles
}

// DriverNames get loaded driver names
func (c *Config) DriverNames() []string {
	return c.driverNames
}

// ClearAll data and caches
func ClearAll() { dc.ClearAll() }

// ClearAll data and caches
func (c *Config) ClearAll() {
	c.ClearData()
	c.ClearCaches()

	c.loadedFiles = []string{}
	c.opts.Readonly = false
}

// ClearData clear data
func (c *Config) ClearData() {
	c.fireHook(OnCleanData)

	c.data = make(map[string]any)
	c.loadedFiles = []string{}
}

// ClearCaches clear caches
func (c *Config) ClearCaches() {
	if c.opts.EnableCache {
		c.intCache = nil
		c.strCache = nil
		c.sMapCache = nil
		c.sArrCache = nil
	}
}

/*************************************************************
 * helper methods
 *************************************************************/

// fire hook
func (c *Config) fireHook(name string) {
	if c.opts.HookFunc != nil {
		c.opts.HookFunc(name, c)
	}
}

// record error
func (c *Config) addError(err error) {
	c.err = err
}

// format and record error
func (c *Config) addErrorf(format string, a ...any) {
	c.err = fmt.Errorf(format, a...)
}
