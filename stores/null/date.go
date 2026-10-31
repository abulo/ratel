package null

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
)

// Date 是仅保留年、月、日的可空 time.Time，时分秒统一置零；JSON 与文本均按 YYYY-MM-DD 格式
// 序列化与解析，同时实现 sql.Scanner 与 driver.Valuer。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type Date struct {
	// Date 为内部保存的日期取值，时分秒为零。
	Date time.Time
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewDate 使用给定时间创建 Date，仅保留其年、月、日并将时分秒置零；valid 为 false 时得到空值。
func NewDate(t time.Time, valid bool) Date {
	return Date{
		Date:  time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		Valid: valid,
		Set:   true,
	}
}

// DateFrom 创建一个始终有效的 Date。
func DateFrom(t time.Time) Date {
	return NewDate(t, true)
}

// DateFromPtr 根据指针创建 Date；指针为 nil 时得到空值。
func DateFromPtr(t *time.Time) Date {
	if t == nil {
		return NewDate(time.Time{}, false)
	}
	return NewDate(*t, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (t Date) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (t Date) IsSet() bool {
	return t.Set
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null，有效值按 YYYY-MM-DD 输出；
// 年份超出 [0, 9999] 范围时返回错误。
func (t Date) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil // PERF-2-17 [P2] 复用包级共享 null 字节，省去 json 反射序列化与分配
	}

	// 序列化实现参考标准库 time.Time 的 JSON 逻辑。
	if y := t.Date.Year(); y < 0 || y >= 10000 {
		// RFC 3339 规定年份必须为四位数字，详见 golang.org/issue/4556#c15。
		return nil, errors.New("Date.MarshalJSON: year outside of range [0,9999]")
	}
	// 预留 YYYY-MM-DD 十二字节与两侧引号的容量。
	b := make([]byte, 0, 12+2)
	b = append(b, '"')
	b = t.Date.AppendFormat(b, RFC3339DateOnly)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 或空对象时置为空值，其余按 YYYY-MM-DD
// 在配置时区解析，解析失败返回对应错误。
func (t *Date) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) || len(data) == 2 {
		t.Valid = false
		t.Date = time.Time{}
		return nil
	}

	// 反序列化实现参考标准库 time.Time 的 JSON 逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.Date, err = time.ParseInLocation(`"`+RFC3339DateOnly+`"`, string(data), util.TimeZone())
	if err != nil {
		return err
	}

	t.Valid = true
	return nil
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 NullBytes，有效值按 YYYY-MM-DD 输出；
// 年份超出 [0, 9999] 范围时返回错误。
func (t Date) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// 序列化实现参考标准库 time.Time 的文本逻辑。
	if y := t.Date.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Date.MarshalText: year outside of range [0,9999]")
	}

	// 预留 YYYY-MM-DD 十二字节的容量。
	b := make([]byte, 0, 12+2)
	return t.Date.AppendFormat(b, RFC3339DateOnly), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，其余按 YYYY-MM-DD 在配置时区解析，
// 解析失败返回对应错误。
func (t *Date) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}

	// 反序列化实现参考标准库 time.Time 的文本逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.Date, err = time.ParseInLocation(RFC3339DateOnly, string(text), util.TimeZone())
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid 使用给定时间更新 Date，仅保留年、月、日并将时分秒置零，同时标记为有效、已赋值。
func (t *Date) SetValid(v time.Time) {
	t.Date = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (t Date) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Date
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (t Date) IsZero() bool {
	return !t.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，仅接受 time.Time 并裁剪为日期，其余类型返回错误。
func (t *Date) Scan(value any) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Date = time.Date(x.Year(), x.Month(), x.Day(), 0, 0, 0, 0, x.Location())
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
func (t Date) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Date, nil
}

// ValueOrDefault 返回内部取值；空值时返回 time.Time 零值。
func (t Date) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Date
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t Date) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.Date.Format(RFC3339DateOnly)
}
