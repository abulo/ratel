package null

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// TimeStamp is a nullable time.Time that accept only fullyear,month,date,hour,minute,second, and microsecond .
// When it is serialize to JSON and vice versa only accept of SQL DateTime Format YYYY-MM-DD hh:mm:ss
// ([4 digit year]-[2 digit month]-[2 digit year] [2 digit hour]:[2 digit minute][2 digit second].[6 digit microsecond]).
// It supports SQL and JSON serialization.
type TimeStamp struct {
	TimeStamp time.Time
	Valid     bool
	Set       bool
}

// NewTimeStamp creates a new Time. that accept only fullyear,month,date,hour,minute,second and microsecond.
func NewTimeStamp(t time.Time, valid bool) TimeStamp {
	return TimeStamp{
		TimeStamp: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()),
		Valid:     valid,
		Set:       true,
	}
}

// TimeStampFrom creates a new TimeStamp that will always be valid.
func TimeStampFrom(t time.Time) TimeStamp {
	return NewTimeStamp(t, true)
}

// TimeStampFromPtr creates a new TimeStamp that will be null if t is nil.
func TimeStampFromPtr(t *time.Time) TimeStamp {
	if t == nil {
		return NewTimeStamp(time.Time{}, false)
	}
	return NewTimeStamp(*t, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (t TimeStamp) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (t TimeStamp) IsSet() bool {
	return t.Set
}

// MarshalJSON implements json.Marshaler.
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	if y := t.TimeStamp.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("DateTime.MarshalJSON: year outside of range [0,9999]")
	}
	b := make([]byte, 0, 26+2) // 26 byte YYYY-MM-DD hh:mm:ss.zzzzzz + ""
	b = append(b, '"')
	b = t.TimeStamp.AppendFormat(b, TimeStampSQL)
	b = append(b, '"')
	return b, nil
	// ---
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) {
		t.Valid = false
		t.TimeStamp = time.Time{}
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.TimeStamp, err = time.Parse(`"`+TimeStampSQL+`"`, string(data))
	if err != nil {
		return err
	}
	// ---

	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t TimeStamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	if y := t.TimeStamp.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("TimeStamp.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, 26+2) // 26 byte YYYY-MM-DD hh:mm:ss.zzzzzz
	return t.TimeStamp.AppendFormat(b, TimeStampSQL), nil
	// ---
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *TimeStamp) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.TimeStamp, err = time.Parse(TimeStampSQL, string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
	// ---
}

// SetValid changes this Time's value that accept only fullyear,month,date,hour,minute,second and microsecond.
// and sets it to be non-null.
func (t *TimeStamp) SetValid(v time.Time) {
	t.TimeStamp = time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr returns a pointer to this Time's value that accept only fullyear,month,date,hour,minute,second and microsecond.
// , or a nil pointer if this Time is null.
func (t TimeStamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.TimeStamp
}

// IsZero returns true for an invalid Time's value, for potential future omitempty support.
func (t TimeStamp) IsZero() bool {
	return !t.Valid
}

// Scan implements the Scanner interface. that accept only fullyear,month,date,hour,minute,second and microsecond.
func (t *TimeStamp) Scan(value interface{}) error {
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

// Value implements the driver Valuer interface.
func (t TimeStamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.TimeStamp, nil
}
