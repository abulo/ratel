package null

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

// Time 是可空的 time.Time，实现 sql.Scanner、driver.Valuer 以及 JSON 与文本编解码接口，
// 序列化格式沿用 time.Time 自身的 RFC3339 格式。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Time struct {
	// Time 为内部保存的时间取值。
	Time time.Time
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewTime 使用给定取值与有效性标记创建 Time；valid 为 false 时得到空值。
func NewTime(t time.Time, valid bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
		Set:   true,
	}
}

// TimeFrom 创建一个始终有效的 Time。
func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

// TimeFromPtr 根据指针创建 Time；指针为 nil 时得到空值。
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (t Time) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (t Time) IsSet() bool {
	return t.Set
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null，有效值委托 time.Time.MarshalJSON。
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 或空对象时置为空值，其余内容委托
// time.Time.UnmarshalJSON 解析，解析失败返回对应错误。
func (t *Time) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) || len(data) == 2 {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}

	if err := t.Time.UnmarshalJSON(data); err != nil {
		return err
	}

	t.Valid = true
	return nil
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 NullBytes，有效值委托 time.Time.MarshalText。
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，其余内容委托
// time.Time.UnmarshalText 解析，解析失败返回对应错误。
func (t *Time) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid 将内部取值置为 v，并同时标记为有效、已赋值。
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
	t.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (t Time) IsZero() bool {
	return !t.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，仅接受 time.Time，其余类型返回扫描错误。
func (t *Time) Scan(value any) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid, t.Set = false, false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	if err == nil {
		t.Valid, t.Set = true, true
	}
	return err
}

// Value 实现 driver.Valuer；空值返回 nil。
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// ValueOrDefault 返回内部取值；空值时返回 time.Time 零值。
func (t Time) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Time) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.Time.String()
}
