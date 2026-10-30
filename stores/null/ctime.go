package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/util"
)

// CTime 是仅保留时、分、秒的可空 time.Time，年、月、日固定为 0、1、1；JSON 与文本均按
// hh:mm:ss 格式序列化与解析，同时实现 sql.Scanner 与 driver.Valuer。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type CTime struct {
	// CTime 为内部保存的时间取值，仅时、分、秒有效。
	CTime time.Time
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewCTime 使用给定时间创建 CTime，仅保留时、分、秒；valid 为 false 时得到空值。
func NewCTime(t time.Time, valid bool) CTime {
	return CTime{
		// 与标准库 time.Parse 对空时间的默认取值保持一致：年 0、月 1、日 1。
		CTime: time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location()),
		Valid: valid,
		Set:   true,
	}
}

// CTimeFrom 创建一个始终有效的 CTime。
func CTimeFrom(t time.Time) CTime {
	return NewCTime(t, true)
}

// CTimeFromPtr 根据指针创建 CTime；指针为 nil 时得到空值。
func CTimeFromPtr(t *time.Time) CTime {
	if t == nil {
		return NewCTime(time.Time{}, false)
	}
	return NewCTime(*t, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (t CTime) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (t CTime) IsSet() bool {
	return t.Set
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null，有效值按 hh:mm:ss 输出。
func (t CTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	// 序列化实现参考标准库 time.Time 的 JSON 逻辑。
	// 预留 hh:mm:ss 八字节与两侧引号的容量。
	b := make([]byte, 0, 8+2)
	b = append(b, '"')
	b = t.CTime.AppendFormat(b, RFC3339TimeOnly)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 或空对象时置为空值，其余按 hh:mm:ss
// 在配置时区解析，解析失败返回对应错误。
func (t *CTime) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) || len(data) == 2 {
		t.Valid = false
		t.CTime = time.Time{}
		return nil
	}

	// 反序列化实现参考标准库 time.Time 的 JSON 逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.CTime, err = time.ParseInLocation(`"`+RFC3339TimeOnly+`"`, string(data), util.TimeZone())
	if err != nil {
		return err
	}

	t.Valid = true
	return nil
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 NullBytes，有效值按 hh:mm:ss 输出。
func (t CTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// 序列化实现参考标准库 time.Time 的文本逻辑。
	// 预留 hh:mm:ss 八字节的容量。
	b := make([]byte, 0, 8+2)
	return t.CTime.AppendFormat(b, RFC3339TimeOnly), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，其余按 hh:mm:ss 在配置时区解析，
// 解析失败返回对应错误。
func (t *CTime) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}
	// 反序列化实现参考标准库 time.Time 的文本逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.CTime, err = time.ParseInLocation(RFC3339TimeOnly, string(text), util.TimeZone())
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid 使用给定时间更新 CTime，仅保留时、分、秒，同时标记为有效、已赋值。
func (t *CTime) SetValid(v time.Time) {
	// 与标准库 time.Parse 对空时间的默认取值保持一致：年 0、月 1、日 1。
	t.CTime = time.Date(0, 1, 1, v.Hour(), v.Minute(), v.Second(), 0, v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (t CTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.CTime
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (t CTime) IsZero() bool {
	return !t.Valid
}

// Scan 实现 sql.Scanner；接受字节切片与 time.Time，value 为 nil 时清空状态，其余类型返回错误。
func (t *CTime) Scan(value any) error {
	var err error
	switch x := value.(type) {
	case []byte:
		t.CTime, _ = time.ParseInLocation(RFC3339TimeOnly, string(x), util.TimeZone())
	case time.Time:
		// 与标准库 time.Parse 对空时间的默认取值保持一致：年 0、月 1、日 1。
		t.CTime = time.Date(0, 1, 1, x.Hour(), x.Minute(), x.Second(), 0, x.Location())
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
func (t CTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.CTime, nil
}

// ValueOrDefault 返回内部取值；空值时返回 time.Time 零值。
func (t CTime) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.CTime
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t CTime) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.CTime.Format(RFC3339TimeOnly)
}
