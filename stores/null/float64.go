package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Float64 is a nullable float64.
type Float64 struct {
	Float64 float64
	Valid   bool
	Set     bool
}

// NewFloat64 creates a new Float64
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		Float64: f,
		Valid:   valid,
		Set:     true,
	}
}

// Float64From creates a new Float64 that will always be valid.
func Float64From(f float64) Float64 {
	return NewFloat64(f, true)
}

// Float64FromPtr creates a new Float64 that be null if f is nil.
func Float64FromPtr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}
	return NewFloat64(*f, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (f Float64) IsValid() bool {
	return f.Set && f.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (f Float64) IsSet() bool {
	return f.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Float64) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(data, NullBytes) {
		f.Float64 = 0
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Float64); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Float64) UnmarshalText(text []byte) error {
	f.Set = true
	if len(text) == 0 {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (f Float64) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// SetValid changes this Float64's value and also sets it to be non-null.
func (f *Float64) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
	f.Set = true
}

// Ptr returns a pointer to this Float64's value, or a nil pointer if this Float64 is null.
func (f Float64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero returns true for invalid Float64s, for future omitempty support (Go 1.4?)
func (f Float64) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *Float64) Scan(value any) error {
	if value == nil {
		f.Float64, f.Valid, f.Set = 0, false, false
		return nil
	}
	f.Valid, f.Set = true, true
	return convert.ConvertAssign(&f.Float64, value)
}

// Value implements the driver Valuer interface.
func (f Float64) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

// ValueOrDefault returns the inner value if valid, otherwise zero.
func (t Float64) ValueOrDefault() float64 {
	if !t.Valid {
		return 0.0
	}
	return t.Float64
}

// String returns the string representation of the float or null.
func (a Float64) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatFloat(a.Float64, 'f', -1, 64)
}
