package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Uint64 是可空的 uint64，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Uint64 struct {
	// Uint64 为内部保存的 uint64 取值。
	Uint64 uint64
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewUint64 使用给定取值与有效性标记创建 Uint64；valid 为 false 时得到空值。
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  valid,
		Set:    true,
	}
}

// Uint64From 创建一个始终有效的 Uint64。
func Uint64From(i uint64) Uint64 {
	return NewUint64(i, true)
}

// Uint64FromPtr 根据指针创建 Uint64；指针为 nil 时得到空值。
func Uint64FromPtr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}
	return NewUint64(*i, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (u Uint64) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (u Uint64) IsSet() bool {
	return u.Set
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 时置为空值，其余内容按无符号整数解析，
// 解析失败时返回对应错误。
func (u *Uint64) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Uint64 = 0
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.Uint64); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，按十进制 uint64 解析，
// 解析失败时返回错误且不更新内部取值。
func (u *Uint64) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 64)
	u.Valid = err == nil
	if u.Valid {
		u.Uint64 = uint64(res)
	}
	return err
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null。
func (u Uint64) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return []byte(strconv.FormatUint(u.Uint64, 10)), nil
}

// MarshalText 实现 encoding.TextMarshaler；空值序列化为空字节切片。
func (u Uint64) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(u.Uint64, 10)), nil
}

// SetValid 将内部取值置为 n，并同时标记为有效、已赋值。
func (u *Uint64) SetValid(n uint64) {
	u.Uint64 = n
	u.Valid = true
	u.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (u Uint64) Ptr() *uint64 {
	if !u.Valid {
		return nil
	}
	return &u.Uint64
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (u Uint64) IsZero() bool {
	return !u.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，否则按驱动值转换后写入内部取值。
func (u *Uint64) Scan(value any) error {
	if value == nil {
		u.Uint64, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true

	// 驱动层可能传入负的 int64，需先转换为 uint64 再赋值，避免直接转换出错。
	if i, ok := value.(int64); ok && i < 0 {
		return convert.ConvertAssign(&u.Uint64, uint64(i))
	}

	return convert.ConvertAssign(&u.Uint64, value)
}

// Value 实现 driver.Valuer；空值返回 nil。
func (u Uint64) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	// 取值超出 int64 正数范围时无法以 int64 承载，改为十进制字符串返回。
	if u.Uint64 >= 1<<63 {
		return strconv.FormatUint(u.Uint64, 10), nil
	}

	return int64(u.Uint64), nil
}

// ValueOrDefault 返回内部取值；空值时返回类型零值。
func (u Uint64) ValueOrDefault() uint64 {
	if !u.Valid {
		return 0
	}
	return u.Uint64
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Uint64) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return strconv.FormatUint(t.Uint64, 10) // PERF-2-19 [P2] 已知无符号整数直接格式化，去掉 cast 反射分发
}
