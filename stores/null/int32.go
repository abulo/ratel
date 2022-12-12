package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/abulo/ratel/v3/stores/null/convert"
)

// Int32 is an nullable int32.
type Int32 struct {
	Int32 int32
	Valid bool
	Set   bool
}

// NewInt32 creates a new Int32
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		Int32: i,
		Valid: valid,
		Set:   true,
	}
}

// Int32From creates a new Int32 that will always be valid.
func Int32From(i int32) Int32 {
	return NewInt32(i, true)
}

// Int32FromPtr creates a new Int32 that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	return NewInt32(*i, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (i Int32) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (i Int32) IsSet() bool {
	return i.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int32) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int32 = 0
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxInt32 {
		return fmt.Errorf("json: %d overflows max int32 value", x)
	}

	i.Int32 = int32(x)
	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int32) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 32)
	i.Valid = err == nil
	if i.Valid {
		i.Int32 = int32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// SetValid changes this Int32's value and also sets it to be non-null.
func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = true
	i.Set = true
}

// Ptr returns a pointer to this Int32's value, or a nil pointer if this Int32 is null.
func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// IsZero returns true for invalid Int32's, for future omitempty support (Go 1.4?)
func (i Int32) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (i *Int32) Scan(value interface{}) error {
	if value == nil {
		i.Int32, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int32, value)
}

// Value implements the driver Valuer interface.
func (i Int32) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int32), nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (i Int32) ValueOrDefault() int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

// String returns the string representation of the int or null.
func (a Int32) Result() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatInt(int64(a.Int32), 10)
}
