package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
	"github.com/spf13/cast"
)

// Uint64 is an nullable uint64.
type Uint64 struct {
	Uint64 uint64
	Valid  bool
	Set    bool
}

// NewUint64 creates a new Uint64
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  valid,
		Set:    true,
	}
}

// Uint64From creates a new Uint64 that will always be valid.
func Uint64From(i uint64) Uint64 {
	return NewUint64(i, true)
}

// Uint64FromPtr creates a new Uint64 that be null if i is nil.
func Uint64FromPtr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}
	return NewUint64(*i, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (u Uint64) IsValid() bool {
	return u.Set && u.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (u Uint64) IsSet() bool {
	return u.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint64) UnmarshalJSON(data []byte) error {
	u.Set = true
	if bytes.Equal(data, NullBytes) {
		u.Uint64 = 0
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.Uint64); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint64) UnmarshalText(text []byte) error {
	u.Set = true
	if len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 64)
	u.Valid = err == nil
	if u.Valid {
		u.Uint64 = uint64(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Uint64) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatUint(u.Uint64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint64) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(u.Uint64, 10)), nil
}

// SetValid changes this Uint64's value and also sets it to be non-null.
func (u *Uint64) SetValid(n uint64) {
	u.Uint64 = n
	u.Valid = true
	u.Set = true
}

// Ptr returns a pointer to this Uint64's value, or a nil pointer if this Uint64 is null.
func (u Uint64) Ptr() *uint64 {
	if !u.Valid {
		return nil
	}
	return &u.Uint64
}

// IsZero returns true for invalid Uint64's, for future omitempty support (Go 1.4?)
func (u Uint64) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Uint64) Scan(value any) error {
	if value == nil {
		u.Uint64, u.Valid, u.Set = 0, false, false
		return nil
	}
	u.Valid, u.Set = true, true

	// If value is negative int64, convert it to uint64
	if i, ok := value.(int64); ok && i < 0 {
		return convert.ConvertAssign(&u.Uint64, uint64(i))
	}

	return convert.ConvertAssign(&u.Uint64, value)
}

// Value implements the driver Valuer interface.
func (u Uint64) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	// If u.Uint64 overflows the range of int64, convert it to string
	if u.Uint64 >= 1<<63 {
		return strconv.FormatUint(u.Uint64, 10), nil
	}

	return int64(u.Uint64), nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (u Uint64) ValueOrDefault() uint64 {
	if !u.Valid {
		return 0
	}
	return u.Uint64
}

// String returns the string representation of the float or null.
func (t Uint64) Result() string {
	if !t.Valid {
		return "NULL"
	}
	return cast.ToString(t.Uint64)
}
