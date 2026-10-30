package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// Uint 是可空的 uint，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Uint struct {
	// Uint 为内部保存的 uint 取值。
	Uint uint
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewUint 使用给定取值与有效性标记创建 Uint；valid 为 false 时得到空值。
func NewUint(i uint, valid bool) Uint {
	return Uint{
		Uint:  i,
		Valid: valid,
		Set:   true,
	}
}

// UintFrom 创建一个始终有效的 Uint。
func UintFrom(i uint) Uint {
	return NewUint(i, true)
}

// UintFromPtr 根据指针创建 Uint；指针为 nil 时得到空值。
func UintFromPtr(i *uint) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	return NewUint(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (u Uint) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (u Uint) IsSet() bool {
	return u.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按无符号整数解析，
// 解析失败时返回对应错误。
func (u *Uint) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	u.Uint = uint(x)
	u.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制无符号整数解析，
// 解析失败时返回错误且不更新内部取值。
func (u *Uint) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 0)
	u.Valid = err == nil
	if u.Valid {
		u.Uint = uint(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (u Uint) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return json.Marshal(nil)
	}
	return []byte(strconv.FormatUint(uint64(u.Uint), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (u Uint) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (u *Uint) SetValid(n uint) {
	u.Uint = n
	u.Valid = true
	u.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (u Uint) Ptr() *uint {
	if !u.Valid {
		return nil
	}
	return &u.Uint
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (u Uint) IsZero() bool {
	return !u.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (u *Uint) Scan(value any) error {
	if value == nil {
		u.Uint, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true
	return convert.ConvertAssign(&u.Uint, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (u Uint) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (u Uint) ValueOrDefault() uint {
	if !u.Valid {
		return 0
	}
	return u.Uint
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Uint) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.Uint)
}
