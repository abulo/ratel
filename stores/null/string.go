package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// String 是可空的 string，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type String struct {
	// String 为内部保存的字符串取值。
	String string
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// StringFrom 创建一个始终有效的 String。
func StringFrom(s string) String {
	return NewString(s, true)
}

// StringFromPtr 根据指针创建 String；指针为 nil 时得到空值。
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

// NewString 使用给定取值与有效性标记创建 String；valid 为 false 时得到空值。
func NewString(s string, valid bool) String {
	return String{
		String: s,
		Valid:  valid,
		Set:    true,
	}
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (s String) IsValid() bool {
	return s.Set && s.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (s String) IsSet() bool {
	return s.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按字符串解析，
// 解析失败时返回对应错误。
func (s *String) UnmarshalJSON(data []byte) error {
	s.Set = true
	if bytes.Equal(data, NullBytes) {
		s.String = ""
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(s.String)
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，非空文本直接写入内部取值。
func (s *String) UnmarshalText(text []byte) error {
	s.Set = true
	if len(text) == 0 {
		s.Valid = false
		return nil
	}

	s.String = string(text)
	s.Valid = true
	return nil
}

// SetValid 将内部取值置为 v，并同时标记为有效、已赋值。
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
	s.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (s String) IsZero() bool {
	return !s.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (s *String) Scan(value any) error {
	if value == nil {
		s.String, s.Valid, s.Set = "", false, false
		return nil
	}
	s.Valid, s.Set = true, true
	return convert.ConvertAssign(&s.String, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (t String) ValueOrDefault() string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t String) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.String)
}
