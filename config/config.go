package config

import (
	"strings"

	"github.com/spf13/cast"
)

func Bool(key string, defVal ...bool) bool { return v.Bool(key, defVal...) }

func (v *Viper) Bool(key string, defVal ...bool) (value bool) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetBool(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func Int(key string, defVal ...int) int { return v.Int(key, defVal...) }

func (v *Viper) Int(key string, defVal ...int) (value int) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetInt(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func Uint(key string, defVal ...uint) uint { return v.Uint(key, defVal...) }

func (v *Viper) Uint(key string, defVal ...uint) (value uint) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetUint(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func Int64(key string, defVal ...int64) int64 { return v.Int64(key, defVal...) }

func (v *Viper) Int64(key string, defVal ...int64) (value int64) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetInt64(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func Ints(key string) []int { return v.Ints(key) }

func (v *Viper) Ints(key string) []int {
	return v.GetIntSlice(key)
}

func IntMap(key string) map[string]int { return v.IntMap(key) }

func (v *Viper) IntMap(key string) map[string]int {
	ret := GetStringMapString(key)
	data := make(map[string]int)
	for k, v := range ret {
		data[k] = cast.ToInt(v)
	}
	return data
}

func Float(key string, defVal ...float64) float64 { return v.Float(key, defVal...) }

func (v *Viper) Float(key string, defVal ...float64) (value float64) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetFloat64(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func String(key string, defVal ...string) string { return v.String(key, defVal...) }

func (v *Viper) String(key string, defVal ...string) (value string) {
	lcaseKey := strings.ToLower(key)
	if v.IsSet(lcaseKey) {
		return GetString(lcaseKey)
	}
	if len(defVal) > 0 {
		value = defVal[0]
	}
	return value
}

func Strings(key string) []string { return v.Strings(key) }

func (v *Viper) Strings(key string) []string {
	return v.GetStringSlice(key)
}

func StringMap(key string) map[string]string { return v.StringMap(key) }

func (v *Viper) StringMap(key string) map[string]string {
	return v.GetStringMapString(key)
}
