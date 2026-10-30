package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Int8 是可空的 int8，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Int8 struct {
	// Int8 为内部保存的 int8 取值。
	Int8 int8
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewInt8 使用给定取值与有效性标记创建 Int8；valid 为 false 时得到空值。
func NewInt8(i int8, valid bool) Int8 {
	return Int8{
		Int8:  i,
		Valid: valid,
		Set:   true,
	}
}

// Int8From 创建一个始终有效的 Int8。
func Int8From(i int8) Int8 {
	return NewInt8(i, true)
}

// Int8FromPtr 根据指针创建 Int8；指针为 nil 时得到空值。
func Int8FromPtr(i *int8) Int8 {
	if i == nil {
		return NewInt8(0, false)
	}
	return NewInt8(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (i Int8) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (i Int8) IsSet() bool {
	return i.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，取值超出 int8 范围时
// 返回溢出错误，其余情况解析失败同样返回对应错误。
func (i *Int8) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int8 = 0
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxInt8 {
		return fmt.Errorf("json: %d overflows max int8 value", x)
	}

	i.Int8 = int8(x)
	i.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 int8 解析，
// 解析失败时返回错误且不更新内部取值。
func (i *Int8) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 8)
	i.Valid = err == nil
	if i.Valid {
		i.Int8 = int8(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return json.Marshal(nil)
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (i Int8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (i *Int8) SetValid(n int8) {
	i.Int8 = n
	i.Valid = true
	i.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (i Int8) Ptr() *int8 {
	if !i.Valid {
		return nil
	}
	return &i.Int8
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (i Int8) IsZero() bool {
	return !i.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (i *Int8) Scan(value any) error {
	if value == nil {
		i.Int8, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int8, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (i Int8) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int8), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (i Int8) ValueOrDefault() int8 {
	if !i.Valid {
		return 0
	}
	return i.Int8
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Int8) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatInt(int64(a.Int8), 10)
}
