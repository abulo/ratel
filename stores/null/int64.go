package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Int64 是可空的 int64，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Int64 struct {
	// Int64 为内部保存的 int64 取值。
	Int64 int64
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewInt64 使用给定取值与有效性标记创建 Int64；valid 为 false 时得到空值。
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		Int64: i,
		Valid: valid,
		Set:   true,
	}
}

// Int64From 创建一个始终有效的 Int64。
func Int64From(i int64) Int64 {
	return NewInt64(i, true)
}

// Int64FromPtr 根据指针创建 Int64；指针为 nil 时得到空值。
func Int64FromPtr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}
	return NewInt64(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (i Int64) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (i Int64) IsSet() bool {
	return i.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按 int64 解析，
// 解析失败时返回对应错误。
func (i *Int64) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int64 = 0
		return nil
	}

	if err := json.Unmarshal(data, &i.Int64); err != nil {
		return err
	}

	i.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 int64 解析，
// 解析失败时返回错误。
func (i *Int64) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (i Int64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (i *Int64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
	i.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (i Int64) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (i Int64) IsZero() bool {
	return !i.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (i *Int64) Scan(value any) error {
	if value == nil {
		i.Int64, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int64, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (i Int64) ValueOrDefault() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Int64) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatInt(a.Int64, 10)
}
