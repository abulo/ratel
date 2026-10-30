package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// Uint16 是可空的 uint16，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Uint16 struct {
	// Uint16 为内部保存的 uint16 取值。
	Uint16 uint16
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewUint16 使用给定取值与有效性标记创建 Uint16；valid 为 false 时得到空值。
func NewUint16(i uint16, valid bool) Uint16 {
	return Uint16{
		Uint16: i,
		Valid:  valid,
		Set:    true,
	}
}

// Uint16From 创建一个始终有效的 Uint16。
func Uint16From(i uint16) Uint16 {
	return NewUint16(i, true)
}

// Uint16FromPtr 根据指针创建 Uint16；指针为 nil 时得到空值。
func Uint16FromPtr(i *uint16) Uint16 {
	if i == nil {
		return NewUint16(0, false)
	}
	return NewUint16(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (u Uint16) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (u Uint16) IsSet() bool {
	return u.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，取值超出 uint16 范围时
// 返回溢出错误，其余情况解析失败同样返回对应错误。
func (u *Uint16) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint16 = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxUint16 {
		return fmt.Errorf("json: %d overflows max uint8 value", x)
	}

	u.Uint16 = uint16(x)
	u.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 uint16 解析，
// 解析失败时返回错误且不更新内部取值。
func (u *Uint16) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 16)
	u.Valid = err == nil
	if u.Valid {
		u.Uint16 = uint16(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (u Uint16) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return json.Marshal(nil)
	}
	return []byte(strconv.FormatUint(uint64(u.Uint16), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (u Uint16) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint16), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (u *Uint16) SetValid(n uint16) {
	u.Uint16 = n
	u.Valid = true
	u.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (u Uint16) Ptr() *uint16 {
	if !u.Valid {
		return nil
	}
	return &u.Uint16
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (u Uint16) IsZero() bool {
	return !u.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (u *Uint16) Scan(value any) error {
	if value == nil {
		u.Uint16, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true
	return convert.ConvertAssign(&u.Uint16, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (u Uint16) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint16), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (u Uint16) ValueOrDefault() uint16 {
	if !u.Valid {
		return 0
	}
	return u.Uint16
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Uint16) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.Uint16)
}
