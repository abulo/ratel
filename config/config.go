package viper

import "github.com/spf13/cast"

func (v *Viper) Bool(key string, defVal ...bool) bool {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetBool(key)
}

func (v *Viper) Int(key string, defVal ...int) int {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetInt(key)
}

func (v *Viper) Uint(key string, defVal ...uint) uint {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetUInt(key)
}

func (v *Viper) GetUInt(key string) uint {
	return cast.ToUint(v.Get(key))
}

func (v *Viper) Int64(key string, defVal ...int64) int64 {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetInt64(key)
}

func (v *Viper) Ints(key string) []int {
	return v.GetIntSlice(key)
}

func IntMap(key string) map[string]int {
	mp := v.GetStringMap(key)

	data := make(map[string]int)
	for k, v := range mp {
		data[k] = cast.ToInt(v)
	}
	return data
}

func (v *Viper) Float(key string, defVal ...float64) float64 {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetFloat64(key)
}

func (v *Viper) String(key string, defVal ...string) string {
	if !v.IsSet(key) {
		if len(defVal) > 0 {
			return defVal[0]
		}
	}
	return v.GetString(key)
}

func (v *Viper) Strings(key string) []string {
	return v.GetStringSlice(key)
}

func (v *Viper) StringMap(key string) map[string]string {
	return v.GetStringMapString(key)
}
