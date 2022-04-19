package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

var (
	errInvalidKey = errors.New("invalid config key string")
	errNotFound   = errors.New("this key does not exist in the config")
)

// Exists key exists check
func Exists(key string, findByPath ...bool) bool { return dc.Exists(key, findByPath...) }

// Exists key exists check
func (c *Config) Exists(key string, findByPath ...bool) (ok bool) {
	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return
	}

	if _, ok = c.data[key]; ok {
		return
	}

	// disable find by path.
	if len(findByPath) > 0 && !findByPath[0] {
		return
	}

	// has sub key? eg. "lang.dir"
	if strings.IndexByte(key, sep) == -1 {
		return
	}

	keys := strings.Split(key, string(sep))
	topK := keys[0]

	// find top item data based on top key
	var item interface{}
	if item, ok = c.data[topK]; !ok {
		return
	}
	for _, k := range keys[1:] {
		switch typeData := item.(type) {
		case map[string]int: // is map(from Set)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[string]string: // is map(from Set)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[string]interface{}: // is map(decode from toml/json)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[interface{}]interface{}: // is map(decode from yaml)
			if item, ok = typeData[k]; !ok {
				return
			}
		case []int: // is array(is from Set)
			i, err := strconv.Atoi(k)

			// check slice index
			if err != nil || len(typeData) < i {
				return false
			}
		case []string: // is array(is from Set)
			i, err := strconv.Atoi(k)
			if err != nil || len(typeData) < i {
				return false
			}
		case []interface{}: // is array(load from file)
			i, err := strconv.Atoi(k)
			if err != nil || len(typeData) < i {
				return false
			}
		default: // error
			return false
		}
	}
	return true
}

/*************************************************************
 * read config data
 *************************************************************/

// Data return all config data
func Data() map[string]interface{} { return dc.Data() }

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
}

// Get config value by key string, support get sub-value by key path(eg. 'map.key'),
// ok is true, find value from config
// ok is false, not found or error
func Get(key string, findByPath ...bool) interface{} { return dc.Get(key, findByPath...) }

// Get config value by key
func (c *Config) Get(key string, findByPath ...bool) interface{} {
	val, _ := c.GetValue(key, findByPath...)
	return val
}

// GetValue get value by given key string.
func GetValue(key string, findByPath ...bool) (interface{}, bool) {
	return dc.GetValue(key, findByPath...)
}

// GetValue get value by given key string.
func (c *Config) GetValue(key string, findByPath ...bool) (value interface{}, ok bool) {
	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		c.addError(errInvalidKey)
		return
	}

	// if not is readonly
	if !c.opts.Readonly {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}

	// is top key
	if value, ok = c.data[key]; ok {
		return
	}

	// disable find by path.
	if len(findByPath) > 0 && !findByPath[0] {
		// c.addError(errNotFound)
		return
	}

	// has sub key? eg. "lang.dir"
	if strings.IndexByte(key, sep) == -1 {
		// c.addError(errNotFound)
		return
	}

	keys := strings.Split(key, string(sep))
	topK := keys[0]

	// find top item data based on top key
	var item interface{}
	if item, ok = c.data[topK]; !ok {
		// c.addError(errNotFound)
		return
	}

	// find child
	// NOTICE: don't merge case, will result in an error.
	// e.g. case []int, []string
	// OR
	// case []int:
	// case []string:
	for _, k := range keys[1:] {
		switch typeData := item.(type) {
		case map[string]int: // is map(from Set)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[string]string: // is map(from Set)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[string]interface{}: // is map(decode from toml/json)
			if item, ok = typeData[k]; !ok {
				return
			}
		case map[interface{}]interface{}: // is map(decode from yaml)
			if item, ok = typeData[k]; !ok {
				return
			}
		case []int: // is array(is from Set)
			i, err := strconv.Atoi(k)

			// check slice index
			if err != nil || len(typeData) < i {
				ok = false
				c.addError(err)
				return
			}

			item = typeData[i]
		case []string: // is array(is from Set)
			i, err := strconv.Atoi(k)
			if err != nil || len(typeData) < i {
				ok = false
				c.addError(err)
				return
			}

			item = typeData[i]
		case []interface{}: // is array(load from file)
			i, err := strconv.Atoi(k)
			if err != nil || len(typeData) < i {
				ok = false
				c.addError(err)
				return
			}

			item = typeData[i]
		default: // error
			ok = false
			c.addErrorf("cannot get value of the key '%s'", key)
			return
		}
	}

	return item, true
}

/*************************************************************
 * read config (basic data type)
 *************************************************************/

// String get a string by key
func String(key string, defVal ...string) string { return dc.String(key, defVal...) }

// String get a string by key, if not found return default value
func (c *Config) String(key string, defVal ...string) string {
	value, ok := c.getString(key)

	if !ok && len(defVal) > 0 { // give default value
		value = defVal[0]
	}
	return value
}

func (c *Config) getString(key string) (value string, ok bool) {
	// find from cache
	if c.opts.EnableCache && len(c.strCache) > 0 {
		value, ok = c.strCache[key]
		if ok {
			return
		}
	}

	val, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typVal := val.(type) {
	// from json `int` always is float64
	case string:
		value = typVal
		if c.opts.ParseEnv {
			value = envutil.ParseEnvValue(value)
		}
	default:
		// value = fmt.Sprintf("%v", val)
		value,_ = strutil.AnyToString(val, false)
	}

	// add cache
	if ok && c.opts.EnableCache {
		if c.strCache == nil {
			c.strCache = make(map[string]string)
		}
		c.strCache[key] = value
	}
	return
}

// Int get a int by key
func Int(key string, defVal ...int) int { return dc.Int(key, defVal...) }

// Int get a int value, if not found return default value
func (c *Config) Int(key string, defVal ...int) (value int) {
	i64, exist := c.tryInt64(key)

	if exist {
		value = int(i64)
	} else if len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// Uint get a uint value, if not found return default value
func Uint(key string, defVal ...uint) uint { return dc.Uint(key, defVal...) }

// Uint get a int value, if not found return default value
func (c *Config) Uint(key string, defVal ...uint) (value uint) {
	i64, exist := c.tryInt64(key)

	if exist {
		value = uint(i64)
	} else if len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// Int64 get a int value, if not found return default value
func Int64(key string, defVal ...int64) int64 { return dc.Int64(key, defVal...) }

// Int64 get a int value, if not found return default value
func (c *Config) Int64(key string, defVal ...int64) (value int64) {
	value, exist := c.tryInt64(key)

	if !exist && len(defVal) > 0 {
		value = defVal[0]
	}
	return
}

// try get a int64 value by given key
func (c *Config) tryInt64(key string) (value int64, ok bool) {
	strVal, ok := c.getString(key)
	if !ok {
		return
	}

	value, err := strconv.ParseInt(strVal, 10, 0)
	if err != nil {
		c.addError(err)
	}
	return
}

// Float get a float64 value, if not found return default value
func Float(key string, defVal ...float64) float64 { return dc.Float(key, defVal...) }

// Float get a float64 by key
func (c *Config) Float(key string, defVal ...float64) (value float64) {
	str, ok := c.getString(key)
	if !ok {
		if len(defVal) > 0 {
			value = defVal[0]
		}
		return
	}

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		c.addError(err)
	}
	return
}

// Bool get a bool value, if not found return default value
func Bool(key string, defVal ...bool) bool { return dc.Bool(key, defVal...) }

// Bool looks up a value for a key in this section and attempts to parse that value as a boolean,
// along with a boolean result similar to a map lookup.
// of following(case insensitive):
//  - true
//  - yes
//  - false
//  - no
//  - 1
//  - 0
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (c *Config) Bool(key string, defVal ...bool) (value bool) {
	rawVal, ok := c.getString(key)
	if !ok {
		if len(defVal) > 0 {
			return defVal[0]
		}
		return
	}

	lowerCase := strings.ToLower(rawVal)
	switch lowerCase {
	case "", "0", "false", "no":
		value = false
	case "1", "true", "yes":
		value = true
	default:
		c.addErrorf("the value '%s' cannot be convert to bool", lowerCase)
	}
	return
}

/*************************************************************
 * read config (complex data type)
 *************************************************************/

// Ints get config data as a int slice/array
func Ints(key string) []int { return dc.Ints(key) }

// Ints get config data as a int slice/array
func (c *Config) Ints(key string) (arr []int) {
	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case []int:
		arr = typeData
	case []interface{}:
		for _, v := range typeData {
			iv, err := mathutil.ToInt(v)
			// iv, err := strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				c.addError(err)
				arr = arr[0:0] // reset
				return
			}

			arr = append(arr, iv)
		}
	default:
		c.addErrorf("value cannot be convert to []int, key is '%s'", key)
	}
	return
}

// IntMap get config data as a map[string]int
func IntMap(key string) map[string]int { return dc.IntMap(key) }

// IntMap get config data as a map[string]int
func (c *Config) IntMap(key string) (mp map[string]int) {
	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case map[string]int: // from Set
		mp = typeData
	case map[string]interface{}: // decode from json,toml
		mp = make(map[string]int)
		for k, v := range typeData {
			// iv, err := strconv.Atoi(fmt.Sprintf("%v", v))
			iv, err := mathutil.ToInt(v)
			if err != nil {
				c.addError(err)
				mp = map[string]int{} // reset
				return
			}
			mp[k] = iv
		}
	case map[interface{}]interface{}: // if decode from yaml
		mp = make(map[string]int)
		for k, v := range typeData {
			// iv, err := strconv.Atoi(fmt.Sprintf( "%v", v))
			iv, err := mathutil.ToInt(v)
			if err != nil {
				c.addError(err)
				mp = map[string]int{} // reset
				return
			}

			// sk := fmt.Sprintf("%v", k)
			sk, _ := strutil.AnyToString(k, false)
			mp[sk] = iv
		}
	default:
		c.addErrorf("value cannot be convert to map[string]int, key is '%s'", key)
	}
	return
}

// Strings get strings by key
func Strings(key string) []string { return dc.Strings(key) }

// Strings get config data as a string slice/array
func (c *Config) Strings(key string) (arr []string) {
	var ok bool
	// find from cache
	if c.opts.EnableCache && len(c.sArrCache) > 0 {
		arr, ok = c.sArrCache[key]
		if ok {
			return
		}
	}

	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case []string:
		arr = typeData
	case []interface{}:
		for _, v := range typeData {
			// arr = append(arr, fmt.Sprintf("%v", v))
			arr = append(arr, strutil.MustString(v))
		}
	default:
		c.addErrorf("value cannot be convert to []string, key is '%s'", key)
		return
	}

	// add cache
	if c.opts.EnableCache {
		if c.sArrCache == nil {
			c.sArrCache = make(map[string]strArr)
		}
		c.sArrCache[key] = arr
	}
	return
}

// StringMap get config data as a map[string]string
func StringMap(key string) map[string]string { return dc.StringMap(key) }

// StringMap get config data as a map[string]string
func (c *Config) StringMap(key string) (mp map[string]string) {
	var ok bool

	// find from cache
	if c.opts.EnableCache && len(c.sMapCache) > 0 {
		mp, ok = c.sMapCache[key]
		if ok {
			return
		}
	}

	rawVal, ok := c.GetValue(key)
	if !ok {
		return
	}

	switch typeData := rawVal.(type) {
	case map[string]string: // from Set
		mp = typeData
	case map[string]interface{}: // decode from json,toml
		mp = make(map[string]string)

		for k, v := range typeData {
			switch v.(type) {
			case string:
				if c.opts.ParseEnv {
					mp[k] = envutil.ParseEnvValue(v.(string))
				} else {
					mp[k] = v.(string)
				}
			default:
				// mp[k] = fmt.Sprintf("%v", v)
				mp[k], _ = strutil.AnyToString(v, false)
			}
		}
	case map[interface{}]interface{}: // decode from yaml
		mp = make(map[string]string)

		for k, v := range typeData {
			sk, _ := strutil.AnyToString(k, false)

			switch typVal := v.(type) {
			case string:
				if c.opts.ParseEnv {
					mp[sk] = envutil.ParseEnvValue(typVal)
				} else {
					mp[sk] = typVal
				}
			default:
				// mp[sk] = fmt.Sprintf("%v", v)
				mp[sk], _ = strutil.AnyToString(v, false)
			}
		}
	default:
		c.addErrorf("value cannot be convert to map[string]string, key is '%s'", key)
		return
	}

	// add cache
	if c.opts.EnableCache {
		if c.sMapCache == nil {
			c.sMapCache = make(map[string]strMap)
		}
		c.sMapCache[key] = mp
	}
	return
}
