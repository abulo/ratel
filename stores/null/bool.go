package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Bool is a nullable bool.
type Bool struct {
	Bool  bool
	Valid bool
	Set   bool
}

// NewBool creates a new Bool
func NewBool(b, valid bool) Bool {
	return Bool{
		Bool:  b,
		Valid: valid,
		Set:   true,
	}
}

// BoolFrom creates a new Bool that will always be valid.
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

// BoolFromPtr creates a new Bool that will be null if f is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (b Bool) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (b Bool) IsSet() bool {
	return b.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.Set = true

	if bytes.Equal(data, NullBytes) {
		b.Bool = false
		b.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return err
	}

	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bool) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Valid = false
		return nil
	}

	str := string(text)
	switch str {
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return NullBytes, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
	b.Set = true
}

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		b.Bool, b.Valid, b.Set = false, false, false
		return nil
	}
	b.Valid, b.Set = true, true
	return convert.ConvertAssign(&b.Bool, value)
}

// Value implements the driver Valuer interface.
func (b Bool) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}
