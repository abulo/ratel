package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Float64 是可空的 float64，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Float64 struct {
	// Float64 为内部保存的 float64 取值。
	Float64 float64
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewFloat64 使用给定取值与有效性标记创建 Float64；valid 为 false 时得到空值。
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		Float64: f,
		Valid:   valid,
		Set:     true,
	}
}

// Float64From 创建一个始终有效的 Float64。
func Float64From(f float64) Float64 {
	return NewFloat64(f, true)
}

// Float64FromPtr 根据指针创建 Float64；指针为 nil 时得到空值。
func Float64FromPtr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}
	return NewFloat64(*f, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (f Float64) IsValid() bool {
	return f.Set && f.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (f Float64) IsSet() bool {
	return f.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按浮点数解析，
// 解析失败时返回对应错误。
func (f *Float64) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(data, NullBytes) {
		f.Float64 = 0
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Float64); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按 float64 解析，
// 解析失败时返回错误。
func (f *Float64) UnmarshalText(text []byte) error {
	f.Set = true
	if len(text) == 0 {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return json.Marshal(nil)
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (f Float64) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (f *Float64) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
	f.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (f Float64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (f Float64) IsZero() bool {
	return !f.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (f *Float64) Scan(value any) error {
	if value == nil {
		f.Float64, f.Valid, f.Set = 0, false, false
		return nil
	}
	f.Valid, f.Set = true, true
	return convert.ConvertAssign(&f.Float64, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (f Float64) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (t Float64) ValueOrDefault() float64 {
	if !t.Valid {
		return 0.0
	}
	return t.Float64
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Float64) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatFloat(a.Float64, 'f', -1, 64)
}
