package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Int 是可空的 int，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Int struct {
	// Int 为内部保存的整数取值。
	Int int
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewInt 使用给定取值与有效性标记创建 Int；valid 为 false 时得到空值。
func NewInt(i int, valid bool) Int {
	return Int{
		Int:   i,
		Valid: valid,
		Set:   true,
	}
}

// IntFrom 创建一个始终有效的 Int。
func IntFrom(i int) Int {
	return NewInt(i, true)
}

// IntFromPtr 根据指针创建 Int；指针为 nil 时得到空值。
func IntFromPtr(i *int) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (i Int) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (i Int) IsSet() bool {
	return i.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按整数解析，
// 解析失败时返回对应错误。
func (i *Int) UnmarshalJSON(data []byte) error {
	i.Set = true

	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int = 0
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	i.Int = int(x)
	i.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制整数解析，
// 解析失败时返回错误且不更新内部取值。
func (i *Int) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 0)
	i.Valid = err == nil
	if i.Valid {
		i.Int = int(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (i *Int) SetValid(n int) {
	i.Int = n
	i.Valid = true
	i.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (i Int) Ptr() *int {
	if !i.Valid {
		return nil
	}
	return &i.Int
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (i Int) IsZero() bool {
	return !i.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (i *Int) Scan(value any) error {
	if value == nil {
		i.Int, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (i Int) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (i Int) ValueOrDefault() int {
	if !i.Valid {
		return 0
	}
	return i.Int
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (a Int) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatInt(int64(a.Int), 10)
}
