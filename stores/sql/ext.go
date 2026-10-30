package sql

import (
	"errors"

	"github.com/spf13/cast"
)

// Pagination 承载分页参数，Offset 为偏移量，Limit 为每页条数；两者均为指针，
// 未设置时对应取值方法返回 nil 错误。
type Pagination struct {
	// Offset 为跳过的记录条数，nil 表示未设置。
	Offset *int64
	// Limit 为返回的记录条数，nil 表示未设置。
	Limit *int64
}

// GetOffset 返回偏移量；Offset 为 nil 时返回错误，调用方须先完成参数绑定与校验。
func (obj *Pagination) GetOffset() (int64, error) {
	if obj.Offset == nil {
		return 0, errors.New("offset is nil")
	}
	return cast.ToInt64(obj.Offset), nil
}

// GetLimit 返回每页条数；Limit 为 nil 时返回错误，调用方须先完成参数绑定与校验。
func (obj *Pagination) GetLimit() (int64, error) {
	if obj.Limit == nil {
		return 0, errors.New("limit is nil")
	}
	return cast.ToInt64(obj.Limit), nil
}

// Scope 承载起止范围参数，Start 与 End 为指针，未设置时对应取值方法返回 nil 错误；
// 取值方法按字符串、整数、浮点三种类型分别提供转换。
type Scope struct {
	// Start 为范围起始值，nil 表示未设置。
	Start *string
	// End 为范围结束值，nil 表示未设置。
	End *string
}

// GetStartString 以字符串形式返回起始值；Start 为 nil 时返回错误。
func (obj *Scope) GetStartString() (string, error) {
	if obj.Start == nil {
		return "", errors.New("start is nil")
	}
	return cast.ToString(obj.Start), nil
}

// GetEndString 以字符串形式返回结束值；End 为 nil 时返回错误。
func (obj *Scope) GetEndString() (string, error) {
	if obj.End == nil {
		return "", errors.New("end is nil")
	}
	return cast.ToString(obj.End), nil
}

// GetStartInt64 以 int64 形式返回起始值；Start 为 nil 时返回错误。
func (obj *Scope) GetStartInt64() (int64, error) {
	if obj.Start == nil {
		return 0, errors.New("start is nil")
	}
	return cast.ToInt64(obj.Start), nil
}

// GetEndInt64 以 int64 形式返回结束值；End 为 nil 时返回错误。
func (obj *Scope) GetEndInt64() (int64, error) {
	if obj.End == nil {
		return 0, errors.New("end is nil")
	}
	return cast.ToInt64(obj.End), nil
}

// GetStartFloat64 以 float64 形式返回起始值；Start 为 nil 时返回错误。
func (obj *Scope) GetStartFloat64() (float64, error) {
	if obj.Start == nil {
		return 0, errors.New("start is nil")
	}
	return cast.ToFloat64(obj.Start), nil
}

// GetEndFloat64 以 float64 形式返回结束值；End 为 nil 时返回错误。
func (obj *Scope) GetEndFloat64() (float64, error) {
	if obj.End == nil {
		return 0, errors.New("end is nil")
	}
	return cast.ToFloat64(obj.End), nil
}

// Sort 承载排序参数，Name 为排序字段，Desc 为是否降序；两者均为指针，
// 未设置时对应取值方法返回 nil 错误。
type Sort struct {
	// Name 为排序字段名，nil 表示未设置。
	Name *string
	// Desc 指定是否降序，nil 表示未设置。
	Desc *bool
}

// GetName 返回排序字段名；Name 为 nil 时返回错误。
func (obj *Sort) GetName() (string, error) {
	if obj.Name == nil {
		return "", errors.New("name is nil")
	}
	return cast.ToString(obj.Name), nil
}

// GetDesc 返回是否降序；Desc 为 nil 时返回错误。
func (obj *Sort) GetDesc() (bool, error) {
	if obj.Desc == nil {
		return false, errors.New("desc is nil")
	}
	return cast.ToBool(obj.Desc), nil
}
