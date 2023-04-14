package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// NullBytes is a global byte slice of JSON null
var NullBytes = []byte("NULL")

// Bytes is a nullable []byte.
type Bytes struct {
	Bytes []byte
	Valid bool
	Set   bool
}

// NewBytes creates a new Bytes
func NewBytes(b []byte, valid bool) Bytes {
	return Bytes{
		Bytes: b,
		Valid: valid,
		Set:   true,
	}
}

// BytesFrom creates a new Bytes that will be invalid if nil.
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, b != nil)
}

// BytesFromPtr creates a new Bytes that will be invalid if nil.
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil {
		return NewBytes(nil, false)
	}
	n := NewBytes(*b, true)
	return n
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (b Bytes) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (b Bytes) IsSet() bool {
	return b.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	b.Set = true

	if bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Bytes = nil
		return nil
	}

	var bv []byte
	if err := json.Unmarshal(data, &bv); err != nil {
		return err
	}

	b.Bytes = bv
	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return json.Marshal(nil)
	}
	if len(b.Bytes) == 0 {
		return NullBytes, nil
	}
	return json.Marshal(b.Bytes)
}

// MarshalText implements encoding.TextMarshaler.
func (b Bytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid changes this Bytes's value and also sets it to be non-null.
func (b *Bytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
	b.Set = true
}

// Ptr returns a pointer to this Bytes's value, or a nil pointer if this Bytes is null.
func (b Bytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero returns true for null or zero Bytes's, for future omitempty support (Go 1.4?)
func (b Bytes) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Bytes) Scan(value any) error {
	if value == nil {
		b.Bytes, b.Valid, b.Set = nil, false, false
		return nil
	}
	b.Valid, b.Set = true, true
	return convert.ConvertAssign(&b.Bytes, value)
}

// Value implements the driver Valuer interface.
func (b Bytes) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (t Bytes) ValueOrDefault() []byte {
	if !t.Valid {
		return []byte{}
	}
	return t.Bytes
}

// String returns the string representation of the float or null.
func (t Bytes) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.Bytes)
}
