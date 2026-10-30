package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
)

// TimeStamp 是保留年、月、日、时、分、秒与纳秒的可空 time.Time，JSON 与文本均按
// YYYY-MM-DD hh:mm:ss.000000 格式序列化与解析，同时实现 sql.Scanner 与 driver.Valuer。
//
// 零值（Set 为 false）表示未赋值，可直接使用；结构体为值类型，不提供并发保护，
// 同一实例的并发读写需由调用方自行加锁。
type TimeStamp struct {
	// TimeStamp 为内部保存的时间戳取值。
	TimeStamp time.Time
	// Valid 表示是否持有非空有效值。
	Valid bool
	// Set 表示是否已被显式赋值，包含赋值为 null 的情形。
	Set bool
}

// NewTimeStamp 使用给定时间创建 TimeStamp，保留其纳秒部分；valid 为 false 时得到空值。
func NewTimeStamp(t time.Time, valid bool) TimeStamp {
	return TimeStamp{
		TimeStamp: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()),
		Valid:     valid,
		Set:       true,
	}
}

// TimeStampFrom 创建一个始终有效的 TimeStamp。
func TimeStampFrom(t time.Time) TimeStamp {
	return NewTimeStamp(t, true)
}

// TimeStampFromPtr 根据指针创建 TimeStamp；指针为 nil 时得到空值。
func TimeStampFromPtr(t *time.Time) TimeStamp {
	if t == nil {
		return NewTimeStamp(time.Time{}, false)
	}
	return NewTimeStamp(*t, true)
}

// IsValid 报告是否持有显式赋值且非空的取值。
func (t TimeStamp) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet 报告是否已被显式赋值，包含赋值为 null 的情形。
func (t TimeStamp) IsSet() bool {
	return t.Set
}

// MarshalJSON 实现 json.Marshaler；空值序列化为 JSON null，有效值按带微秒的时间戳格式输出；
// 年份超出 [0, 9999] 范围时返回错误。
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	// 序列化实现参考标准库 time.Time 的 JSON 逻辑。
	if y := t.TimeStamp.Year(); y < 0 || y >= 10000 {
		// RFC 3339 规定年份必须为四位数字，详见 golang.org/issue/4556#c15。
		return nil, errors.New("DateTime.MarshalJSON: year outside of range [0,9999]")
	}
	// 预留 YYYY-MM-DD hh:mm:ss.zzzzzz 二十六字节与两侧引号的容量。
	b := make([]byte, 0, 26+2)
	b = append(b, '"')
	b = t.TimeStamp.AppendFormat(b, TimeStampSQL)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON 实现 json.Unmarshaler；输入为 JSON null 或空对象时置为空值，其余按带微秒的
// 时间戳格式在配置时区解析，解析失败返回对应错误。
func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) || len(data) == 2 {
		t.Valid = false
		t.TimeStamp = time.Time{}
		return nil
	}

	// 反序列化实现参考标准库 time.Time 的 JSON 逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.TimeStamp, err = time.ParseInLocation(`"`+TimeStampSQL+`"`, string(data), util.TimeZone())
	if err != nil {
		return err
	}

	t.Valid = true
	return nil
}

// MarshalText 实现 encoding.TextMarshaler；空值返回 NullBytes，有效值按带微秒的时间戳格式输出；
// 年份超出 [0, 9999] 范围时返回错误。
func (t TimeStamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// 序列化实现参考标准库 time.Time 的文本逻辑。
	if y := t.TimeStamp.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("TimeStamp.MarshalText: year outside of range [0,9999]")
	}

	// 预留 YYYY-MM-DD hh:mm:ss.zzzzzz 二十六字节的容量。
	b := make([]byte, 0, 26+2)
	return t.TimeStamp.AppendFormat(b, TimeStampSQL), nil
}

// UnmarshalText 实现 encoding.TextUnmarshaler；空文本置为空值，其余按带微秒的时间戳格式
// 在配置时区解析，解析失败返回对应错误。
func (t *TimeStamp) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}

	// 反序列化实现参考标准库 time.Time 的文本逻辑；小数秒由 Parse 隐式处理。
	var err error
	t.TimeStamp, err = time.ParseInLocation(TimeStampSQL, string(text), util.TimeZone())
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid 使用给定时间更新 TimeStamp，保留其纳秒部分，同时标记为有效、已赋值。
func (t *TimeStamp) SetValid(v time.Time) {
	t.TimeStamp = time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr 返回指向内部取值的指针；空值时返回 nil。
func (t TimeStamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.TimeStamp
}

// IsZero 报告是否为空值，供结构体标签 omitempty 判定使用。
func (t TimeStamp) IsZero() bool {
	return !t.Valid
}

// Scan 实现 sql.Scanner；value 为 nil 时清空状态，仅接受 time.Time 并保留纳秒，其余类型返回错误。
func (t *TimeStamp) Scan(value any) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.TimeStamp = time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), x.Nanosecond(), x.Location())
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
func (t TimeStamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.TimeStamp, nil
}

// ValueOrDefault 返回内部取值；空值时返回 time.Time 零值。
func (t TimeStamp) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.TimeStamp
}

// Result 返回内部取值的字符串形式；空值返回 "NULL"。
func (t TimeStamp) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.TimeStamp.Format(TimeStampSQL)
}
