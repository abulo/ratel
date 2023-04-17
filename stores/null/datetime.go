package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// DateTime is a nullable time.Time that accept only fullyear,month,date,hour,minute,second and ignore elses.
// When it is serialize to JSON and vice versa only accept of SQL DateTime Format YYYY-MM-DD hh:mm:ss
// ([4 digit year]-[2 digit month]-[2 digit year] [2 digit hour]:[2 digit minute][2 digit second]).
// It supports SQL and JSON serialization.
type DateTime struct {
	DateTime time.Time
	Valid    bool
	Set      bool
}

// NewDateTime creates a new Time. that accept only fullyear,month,date,hour,minute,second and ignore elses.
func NewDateTime(t time.Time, valid bool) DateTime {
	return DateTime{
		DateTime: time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location()),
		Valid:    valid,
		Set:      true,
	}
}

// DateTimeFrom creates a new DateTime that will always be valid.
func DateTimeFrom(t time.Time) DateTime {
	return NewDateTime(t, true)
}

// DateTimeFromPtr creates a new DateTime that will be null if t is nil.
func DateTimeFromPtr(t *time.Time) DateTime {
	if t == nil {
		return NewDateTime(time.Time{}, false)
	}
	return NewDateTime(*t, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (t DateTime) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (t DateTime) IsSet() bool {
	return t.Set
}

// MarshalJSON implements json.Marshaler.
func (t DateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	// customize from golang time/time.go
	if y := t.DateTime.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("DateTime.MarshalJSON: year outside of range [0,9999]")
	}
	b := make([]byte, 0, 19+2) // 19 byte YYYY-MM-DD hh:mm:ss+ ""
	b = append(b, '"')
	b = t.DateTime.AppendFormat(b, DateTimeSQL)
	b = append(b, '"')
	return b, nil
	// ---
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *DateTime) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) {
		t.Valid = false
		t.DateTime = time.Time{}
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.DateTime, err = time.Parse(`"`+DateTimeSQL+`"`, string(data))
	if err != nil {
		return err
	}
	// ---

	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t DateTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	if y := t.DateTime.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("DateTime.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, 19+2) // 19 byte YYYY-MM-DD hh:mm:ss
	return t.DateTime.AppendFormat(b, DateTimeSQL), nil
	// ---
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *DateTime) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.DateTime, err = time.Parse(DateTimeSQL, string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
	// ---
}

// SetValid changes this Time's value that accept only fullyear,month,date and ignore elses.
// and sets it to be non-null.
func (t *DateTime) SetValid(v time.Time) {
	t.DateTime = time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), 0, v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr returns a pointer to this Time's value that accept only fullyear,month,date and ignore elses.
// , or a nil pointer if this Time is null.
func (t DateTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.DateTime
}

// IsZero returns true for an invalid Time's value, for potential future omitempty support.
func (t DateTime) IsZero() bool {
	return !t.Valid
}

// Scan implements the Scanner interface. that accept only fullyear,month,date and ignore elses.
func (t *DateTime) Scan(value any) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.DateTime = time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), 0, x.Location())
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
func (t DateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.DateTime, nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (t DateTime) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.DateTime
}

// String returns the string representation of the float or null.
func (t DateTime) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.DateTime.Format(DateTimeSQL)
}
