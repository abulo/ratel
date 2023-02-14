package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/abulo/ratel/v2/stores/null/convert"
	"github.com/spf13/cast"
)

// String is a nullable string. It supports SQL and JSON serialization.
type String struct {
	String string
	Valid  bool
	Set    bool
}

// StringFrom creates a new String that will never be blank.
func StringFrom(s string) String {
	return NewString(s, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		String: s,
		Valid:  valid,
		Set:    true,
	}
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (s String) IsValid() bool {
	return s.Set && s.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (s String) IsSet() bool {
	return s.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(data []byte) error {
	s.Set = true
	if bytes.Equal(data, NullBytes) {
		s.String = ""
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return NullBytes, nil
	}
	return json.Marshal(s.String)
}

// MarshalText implements encoding.TextMarshaler.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *String) UnmarshalText(text []byte) error {
	s.Set = true
	if len(text) == 0 {
		s.Valid = false
		return nil
	}

	s.String = string(text)
	s.Valid = true
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
	s.Set = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s String) IsZero() bool {
	return !s.Valid
}

// Scan implements the Scanner interface.
func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.String, s.Valid, s.Set = "", false, false
		return nil
	}
	s.Valid, s.Set = true, true
	return convert.ConvertAssign(&s.String, value)
}

// Value implements the driver Valuer interface.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (t String) ValueOrDefault() string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Result returns the string representation of the float or null.
func (t String) Result() string {
	if !t.Valid {
		return "null"
	}
	return cast.ToString(t.String)
}
