package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// Uint8 is an nullable uint8.
type Uint8 struct {
	Uint8 uint8
	Valid bool
	Set   bool
}

// NewUint8 creates a new Uint8
func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		Uint8: i,
		Valid: valid,
		Set:   true,
	}
}

// Uint8From creates a new Uint8 that will always be valid.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, true)
}

// Uint8FromPtr creates a new Uint8 that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	return NewUint8(*i, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (u Uint8) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (u Uint8) IsSet() bool {
	return u.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint8) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint8 = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxUint8 {
		return fmt.Errorf("json: %d overflows max uint8 value", x)
	}

	u.Uint8 = uint8(x)
	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint8) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 8)
	u.Valid = err == nil
	if u.Valid {
		u.Uint8 = uint8(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Uint8) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint8), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint8) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint8), 10)), nil
}

// SetValid changes this Uint8's value and also sets it to be non-null.
func (u *Uint8) SetValid(n uint8) {
	u.Uint8 = n
	u.Valid = true
	u.Set = true
}

// Ptr returns a pointer to this Uint8's value, or a nil pointer if this Uint8 is null.
func (u Uint8) Ptr() *uint8 {
	if !u.Valid {
		return nil
	}
	return &u.Uint8
}

// IsZero returns true for invalid Uint8's, for future omitempty support (Go 1.4?)
func (u Uint8) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Uint8) Scan(value any) error {
	if value == nil {
		u.Uint8, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true
	return convert.ConvertAssign(&u.Uint8, value)
}

// Value implements the driver Valuer interface.
func (u Uint8) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint8), nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (u Uint8) ValueOrDefault() uint8 {
	if !u.Valid {
		return 0
	}
	return u.Uint8
}

// String returns the string representation of the float or null.
func (t Uint8) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.Uint8)
}
