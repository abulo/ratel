package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Float32 is a nullable float32.
type Float32 struct {
	Float32 float32
	Valid   bool
	Set     bool
}

// NewFloat32 creates a new Float32
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		Float32: f,
		Valid:   valid,
		Set:     true,
	}
}

// Float32From creates a new Float32 that will always be valid.
func Float32From(f float32) Float32 {
	return NewFloat32(f, true)
}

// Float32FromPtr creates a new Float32 that be null if f is nil.
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (f Float32) IsValid() bool {
	return f.Set && f.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (f Float32) IsSet() bool {
	return f.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Float32) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(data, NullBytes) {
		f.Valid = false
		f.Float32 = 0
		return nil
	}

	var x float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	f.Float32 = float32(x)
	f.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Float32) UnmarshalText(text []byte) error {
	f.Set = true
	if len(text) == 0 {
		f.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseFloat(string(text), 32)
	f.Valid = err == nil
	if f.Valid {
		f.Float32 = float32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (f Float32) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// SetValid changes this Float32's value and also sets it to be non-null.
func (f *Float32) SetValid(n float32) {
	f.Float32 = n
	f.Valid = true
	f.Set = true
}

// Ptr returns a pointer to this Float32's value, or a nil pointer if this Float32 is null.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero returns true for invalid Float32s, for future omitempty support (Go 1.4?)
func (f Float32) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *Float32) Scan(value any) error {
	if value == nil {
		f.Float32, f.Valid, f.Set = 0, false, false
		return nil
	}
	f.Valid, f.Set = true, true
	return convert.ConvertAssign(&f.Float32, value)
}

// Value implements the driver Valuer interface.
func (f Float32) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return float64(f.Float32), nil
}

// ValueOrDefault returns the inner value if valid, otherwise zero.
func (t Float32) ValueOrDefault() float32 {
	if !t.Valid {
		return 0.0
	}
	return t.Float32
}

// String returns the string representation of the float or null.
func (a Float32) Result() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatFloat(float64(a.Float32), 'f', -1, 32)
}
