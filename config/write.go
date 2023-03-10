package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/imdario/mergo"
)

var (
	errReadonly   = errors.New("the config instance in 'readonly' mode")
	errKeyIsEmpty = errors.New("the config key is cannot be empty")
)

// SetData for override the Config.Data
func SetData(data map[string]any) {
	dc.SetData(data)
}

// SetData for override the Config.Data
func (c *Config) SetData(data map[string]any) {
	c.lock.Lock()
	c.data = data
	c.lock.Unlock()

	c.fireHook(OnSetData)
}

// Set val by key
func Set(key string, val any, setByPath ...bool) error {
	return dc.Set(key, val, setByPath...)
}

// Set a value by key string.
func (c *Config) Set(key string, val any, setByPath ...bool) (err error) {
	if c.opts.Readonly {
		return errReadonly
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return errKeyIsEmpty
	}

	defer c.fireHook(OnSetValue)
	if strings.IndexByte(key, sep) == -1 {
		c.data[key] = val
		return
	}

	// disable set by path.
	if len(setByPath) > 0 && !setByPath[0] {
		c.data[key] = val
		return
	}

	keys := strings.Split(key, string(sep))
	topK := keys[0]
	paths := keys[1:]

	var ok bool
	var item any

	// find top item data based on top key
	if item, ok = c.data[topK]; !ok {
		// not found, is new add
		c.data[topK] = buildValueByPath(paths, val)
		return
	}

	switch typeData := item.(type) {
	case map[any]any: // from yaml
		dstItem := make(map[string]any)
		for k, v := range typeData {
			sk := fmt.Sprintf("%v", k)
			dstItem[sk] = v
		}

		// create a new item for the topK
		newItem := buildValueByPath(paths, val)
		// merge new item to old item
		err = mergo.Merge(&dstItem, newItem, mergo.WithOverride)
		if err != nil {
			return
		}

		c.data[topK] = dstItem
	case map[string]any: // from json,toml
		// create a new item for the topK
		newItem := buildValueByPath(paths, val)
		// merge new item to old item
		err = mergo.Merge(&typeData, &newItem, mergo.WithOverride)
		if err != nil {
			return
		}

		c.data[topK] = typeData
	case []any: // is array
		index, err := strconv.Atoi(keys[1])
		if len(keys) == 2 && err == nil {
			if index <= len(typeData) {
				typeData[index] = val
			}

			c.data[topK] = typeData
		} else {
			err = errors.New("max allow 1 level for setting array value, current key: " + key)
			return err
		}
	default:
		// as a top key
		c.data[key] = val
		// err = errors.New("not supported value type, cannot setting value for the key: " + key)
	}
	return
}

/**
more setter: SetIntArr, SetIntMap, SetString, SetStringArr, SetStringMap
*/

// build new value by key paths
// "site.info" -> map[string]map[string]val
func buildValueByPath(paths []string, val any) (newItem map[string]any) {
	if len(paths) == 1 {
		return map[string]any{paths[0]: val}
	}

	sliceReverse(paths)

	// multi nodes
	for _, p := range paths {
		if newItem == nil {
			newItem = map[string]any{p: val}
		} else {
			newItem = map[string]any{p: newItem}
		}
	}
	return
}

// reverse a slice. (slice 是引用，所以可以直接改变)
func sliceReverse(ss []string) {
	ln := len(ss)
	for i := 0; i < ln/2; i++ {
		li := ln - i - 1
		ss[i], ss[li] = ss[li], ss[i]
	}
}
