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

// Int16 是可空的 int16，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Int16 struct {
	// Int16 为内部保存的 int16 取值。
	Int16 int16
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewInt16 使用给定取值与有效性标记创建 Int16；valid 为 false 时得到空值。
func NewInt16(i int16, valid bool) Int16 {
	return Int16{
		Int16: i,
		Valid: valid,
		Set:   true,
	}
}

// Int16From 创建一个始终有效的 Int16。
func Int16From(i int16) Int16 {
	return NewInt16(i, true)
}

// Int16FromPtr 根据指针创建 Int16；指针为 nil 时得到空值。
func Int16FromPtr(i *int16) Int16 {
	if i == nil {
		return NewInt16(0, false)
	}
	return NewInt16(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (i Int16) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (i Int16) IsSet() bool {
	return i.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，取值超出 int16 范围时
// 返回溢出错误，其余情况解析失败同样返回对应错误。
func (i *Int16) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Int16, i.Valid = 0, false
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxInt16 {
		return fmt.Errorf("json: %d overflows max int16 value", x)
	}

	i.Int16, i.Valid = int16(x), true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 int16 解析，
// 解析失败时返回错误且不更新内部取值。
func (i *Int16) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 16)
	i.Valid = err == nil
	if i.Valid {
		i.Int16 = int16(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (i Int16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (i Int16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int16), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (i *Int16) SetValid(n int16) {
	i.Int16 = n
	i.Valid = true
	i.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (i Int16) Ptr() *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (i Int16) IsZero() bool {
	return !i.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (i *Int16) Scan(value any) error {
	if value == nil {
		i.Int16, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int16, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (i Int16) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int16), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (i Int16) ValueOrDefault() int16 {
	if !i.Valid {
		return 0
	}
	return i.Int16
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Int16) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatInt(int64(a.Int16), 10)
}
