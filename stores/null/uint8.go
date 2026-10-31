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

// Uint8 是可空的 uint8，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Uint8 struct {
	// Uint8 为内部保存的 uint8 取值。
	Uint8 uint8
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewUint8 使用给定取值与有效性标记创建 Uint8；valid 为 false 时得到空值。
func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		Uint8: i,
		Valid: valid,
		Set:   true,
	}
}

// Uint8From 创建一个始终有效的 Uint8。
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, true)
}

// Uint8FromPtr 根据指针创建 Uint8；指针为 nil 时得到空值。
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	return NewUint8(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (u Uint8) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (u Uint8) IsSet() bool {
	return u.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，取值超出 uint8 范围时
// 返回溢出错误，其余情况解析失败同样返回对应错误。
func (u *Uint8) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint8 = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxUint8 {
		return fmt.Errorf("json: %d overflows max uint8 value", x)
	}

	u.Uint8 = uint8(x)
	u.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 uint8 解析，
// 解析失败时返回错误且不更新内部取值。
func (u *Uint8) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 8)
	u.Valid = err == nil
	if u.Valid {
		u.Uint8 = uint8(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (u Uint8) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatUint(uint64(u.Uint8), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (u Uint8) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint8), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (u *Uint8) SetValid(n uint8) {
	u.Uint8 = n
	u.Valid = true
	u.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (u Uint8) Ptr() *uint8 {
	if !u.Valid {
		return nil
	}
	return &u.Uint8
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (u Uint8) IsZero() bool {
	return !u.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (u *Uint8) Scan(value any) error {
	if value == nil {
		u.Uint8, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true
	return convert.ConvertAssign(&u.Uint8, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (u Uint8) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint8), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (u Uint8) ValueOrDefault() uint8 {
	if !u.Valid {
		return 0
	}
	return u.Uint8
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Uint8) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return strconv.FormatUint(uint64(t.Uint8), 10) // PERF-2-19 [P2] 已知无符号整数直接格式化，去掉 cast 反射分发
}
