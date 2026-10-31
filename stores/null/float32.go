package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Float32 是可空的 float32，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Float32 struct {
	// Float32 为内部保存的 float32 取值。
	Float32 float32
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewFloat32 使用给定取值与有效性标记创建 Float32；valid 为 false 时得到空值。
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		Float32: f,
		Valid:   valid,
		Set:     true,
	}
}

// Float32From 创建一个始终有效的 Float32。
func Float32From(f float32) Float32 {
	return NewFloat32(f, true)
}

// Float32FromPtr 根据指针创建 Float32；指针为 nil 时得到空值。
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (f Float32) IsValid() bool {
	return f.Set && f.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (f Float32) IsSet() bool {
	return f.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按浮点数解析，
// 解析失败时返回对应错误。
func (f *Float32) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(data, NullBytes) {
		f.Valid = false
		f.Float32 = 0
		return nil
	}

	var x float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	f.Float32 = float32(x)
	f.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按 float32 解析，
// 解析失败时返回错误且不更新内部取值。
func (f *Float32) UnmarshalText(text []byte) error {
	f.Set = true
	if len(text) == 0 {
		f.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseFloat(string(text), 32)
	f.Valid = err == nil
	if f.Valid {
		f.Float32 = float32(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (f Float32) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (f *Float32) SetValid(n float32) {
	f.Float32 = n
	f.Valid = true
	f.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (f Float32) IsZero() bool {
	return !f.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (f *Float32) Scan(value any) error {
	if value == nil {
		f.Float32, f.Valid, f.Set = 0, false, false
		return nil
	}
	f.Valid, f.Set = true, true
	return convert.ConvertAssign(&f.Float32, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (f Float32) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return float64(f.Float32), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (t Float32) ValueOrDefault() float32 {
	if !t.Valid {
		return 0.0
	}
	return t.Float32
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Float32) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatFloat(float64(a.Float32), 'f', -1, 32)
}
