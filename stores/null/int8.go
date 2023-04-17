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

// Int8 is an nullable int8.
type Int8 struct {
	Int8  int8
	Valid bool
	Set   bool
}

// NewInt8 creates a new Int8
func NewInt8(i int8, valid bool) Int8 {
	return Int8{
		Int8:  i,
		Valid: valid,
		Set:   true,
	}
}

// Int8From creates a new Int8 that will always be valid.
func Int8From(i int8) Int8 {
	return NewInt8(i, true)
}

// Int8FromPtr creates a new Int8 that be null if i is nil.
func Int8FromPtr(i *int8) Int8 {
	if i == nil {
		return NewInt8(0, false)
	}
	return NewInt8(*i, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (i Int8) IsValid() bool {
	return i.Set && i.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (i Int8) IsSet() bool {
	return i.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int8) UnmarshalJSON(data []byte) error {
	i.Set = true
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int8 = 0
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxInt8 {
		return fmt.Errorf("json: %d overflows max int8 value", x)
	}

	i.Int8 = int8(x)
	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int8) UnmarshalText(text []byte) error {
	i.Set = true
	if len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 8)
	i.Valid = err == nil
	if i.Valid {
		i.Int8 = int8(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// SetValid changes this Int8's value and also sets it to be non-null.
func (i *Int8) SetValid(n int8) {
	i.Int8 = n
	i.Valid = true
	i.Set = true
}

// Ptr returns a pointer to this Int8's value, or a nil pointer if this Int8 is null.
func (i Int8) Ptr() *int8 {
	if !i.Valid {
		return nil
	}
	return &i.Int8
}

// IsZero returns true for invalid Int8's, for future omitempty support (Go 1.4?)
func (i Int8) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (i *Int8) Scan(value any) error {
	if value == nil {
		i.Int8, i.Valid, i.Set = 0, false, false
		return nil
	}
	i.Valid, i.Set = true, true
	return convert.ConvertAssign(&i.Int8, value)
}

// Value implements the driver Valuer interface.
func (i Int8) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int8), nil
}

// ValueOrDefault returns the inner value if valid, otherwise default.
func (i Int8) ValueOrDefault() int8 {
	if !i.Valid {
		return 0
	}
	return i.Int8
}

// String returns the string representation of the int or null.
func (a Int8) Result() string {
	if !a.Valid {
		return "NULL"
	}
	return strconv.FormatInt(int64(a.Int8), 10)
}
