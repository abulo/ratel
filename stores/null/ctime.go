package null

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

// CTime (Customize Time) is a nullable time.Time that accept only hour,minutes,seconds and ignore elses.
// When it is serialize to JSON and vice versa only accept of front half RFC 3339
// ([2 digit hour]:[2 digit minutes]:[2 digit seconds]). It supports SQL and JSON serialization.
type CTime struct {
	CTime time.Time
	Valid bool
	Set   bool
}

// NewCTime creates a new Time. that accept only hour,minutes,seconds and ignore elses.
func NewCTime(t time.Time, valid bool) CTime {
	return CTime{
		// because time.parser setting default year 0, month 1, and day 1 for emptied time.Time so we follow it
		CTime: time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location()),
		Valid: valid,
		Set:   true,
	}
}

// CTimeFrom creates a new CTime that will always be valid.
func CTimeFrom(t time.Time) CTime {
	return NewCTime(t, true)
}

// CTimeFromPtr creates a new CTime that will be null if t is nil.
func CTimeFromPtr(t *time.Time) CTime {
	if t == nil {
		return NewCTime(time.Time{}, false)
	}
	return NewCTime(*t, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (t CTime) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (t CTime) IsSet() bool {
	return t.Set
}

// MarshalJSON implements json.Marshaler.
func (t CTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	b := make([]byte, 0, 8+2) // 8 byte hh:mm:ss + ""
	b = append(b, '"')
	b = t.CTime.AppendFormat(b, RFC3339TimeOnly)
	b = append(b, '"')
	return b, nil
	// ---
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *CTime) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) {
		t.Valid = false
		t.CTime = time.Time{}
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.CTime, err = time.Parse(`"`+RFC3339TimeOnly+`"`, string(data))
	if err != nil {
		return err
	}
	// ---

	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t CTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	b := make([]byte, 0, 8+2) // 8 byte hh:mm:ss
	return t.CTime.AppendFormat(b, RFC3339TimeOnly), nil
	// ---
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *CTime) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}
	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.CTime, err = time.Parse(RFC3339TimeOnly, string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
	// ---
}

// SetValid changes this Time's value that accept only hour,minutes,seconds and ignore elses.
// and sets it to be non-null.
func (t *CTime) SetValid(v time.Time) {
	// because time.parser setting default year 0, month 1, and day 1 for emptied time.Time so we follow it
	t.CTime = time.Date(0, 1, 1, v.Hour(), v.Minute(), v.Second(), 0, v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr returns a pointer to this Time's value that accept only hour,minutes,seconds and ignore elses.
// , or a nil pointer if this Time is null.
func (t CTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.CTime
}

// IsZero returns true for an invalid Time's value, for potential future omitempty support.
func (t CTime) IsZero() bool {
	return !t.Valid
}

// Scan implements the Scanner interface. that accept only hour,minutes,seconds and ignore elses.
func (t *CTime) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case []byte:
		t.CTime = time.Date(0, 1, 1, cast.ToInt(cast.ToString(x[0:2])), cast.ToInt(cast.ToString(x[3:5])), cast.ToInt(cast.ToString(x[6:8])), 0, util.TimeZone())
	case time.Time:
		// because time.parser setting default year 0, month 1, and day 1 for emptied time.Time so we follow it
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

// Value implements the driver Valuer interface.
func (t CTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.CTime, nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (t CTime) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.CTime
}

// String returns the string representation of the float or null.
func (t CTime) Result() string {
	if !t.Valid {
		return "null"
	}
	return t.CTime.Format(RFC3339TimeOnly)
}
