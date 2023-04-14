package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Date is a nullable time.Time that accept only fullyear,month,date and ignore elses.
// When it is serialize to JSON and vice versa only accept of front half RFC 3339
// ([4 digit year]-[2 digit month]-[2 digit year]). It supports SQL and JSON serialization.
type Date struct {
	Date  time.Time
	Valid bool
	Set   bool
}

// NewDate creates a new Time. that accept only fullyear,month,date and ignore elses.
func NewDate(t time.Time, valid bool) Date {
	return Date{
		Date:  time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		Valid: valid,
		Set:   true,
	}
}

// DateFrom creates a new Date that will always be valid.
func DateFrom(t time.Time) Date {
	return NewDate(t, true)
}

// DateFromPtr creates a new Date that will be null if t is nil.
func DateFromPtr(t *time.Time) Date {
	if t == nil {
		return NewDate(time.Time{}, false)
	}
	return NewDate(*t, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (t Date) IsValid() bool {
	return t.Set && t.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (t Date) IsSet() bool {
	return t.Set
}

// MarshalJSON implements json.Marshaler.
func (t Date) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return NullBytes, nil
	}

	// customize from golang time/time.go
	if y := t.Date.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Date.MarshalJSON: year outside of range [0,9999]")
	}
	b := make([]byte, 0, 12+2) // 12 byte YYYY-MM-DD + ""
	b = append(b, '"')
	b = t.Date.AppendFormat(b, RFC3339DateOnly)
	b = append(b, '"')
	return b, nil
	// ---
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Date) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, NullBytes) {
		t.Valid = false
		t.Date = time.Time{}
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.Date, err = time.Parse(`"`+RFC3339DateOnly+`"`, string(data))
	if err != nil {
		return err
	}
	// ---

	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (t Date) MarshalText() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	// customize from golang time/time.go
	if y := t.Date.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Date.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, 12+2) // 12 byte YYYY-MM-DD
	return t.Date.AppendFormat(b, RFC3339DateOnly), nil
	// ---
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Date) UnmarshalText(text []byte) error {
	t.Set = true
	if len(text) == 0 {
		t.Valid = false
		return nil
	}

	// customize from golang time/time.go
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.Date, err = time.Parse(RFC3339DateOnly, string(text))
	if err != nil {
		return err
	}
	t.Valid = true
	return nil
	// ---
}

// SetValid changes this Time's value that accept only fullyear,month,date and ignore elses.
// and sets it to be non-null.
func (t *Date) SetValid(v time.Time) {
	t.Date = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
	t.Valid = true
	t.Set = true
}

// Ptr returns a pointer to this Time's value that accept only fullyear,month,date and ignore elses.
// , or a nil pointer if this Time is null.
func (t Date) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Date
}

// IsZero returns true for an invalid Time's value, for potential future omitempty support.
func (t Date) IsZero() bool {
	return !t.Valid
}

// Scan implements the Scanner interface. that accept only fullyear,month,date and ignore elses.
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

// Value implements the driver Valuer interface.
func (t Date) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Date, nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (t Date) ValueOrDefault() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Date
}

// String returns the string representation of the float or null.
func (t Date) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return t.Date.Format(RFC3339DateOnly)
}
