package null

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"testing"
)

var (
	uint8JSON = []byte(`254`)
)

func TestUint8From(t *testing.T) {
	i := Uint8From(254)
	assertUint8(t, i, "Uint8From()")

	zero := Uint8From(0)
	if !zero.Valid {
		t.Error("Uint8From(0)", "is invalid, but should be valid")
	}
}

func TestUint8FromPtr(t *testing.T) {
	n := uint8(254)
	iptr := &n
	i := Uint8FromPtr(iptr)
	assertUint8(t, i, "Uint8FromPtr()")

	null := Uint8FromPtr(nil)
	assertNullUint8(t, null, "Uint8FromPtr(nil)")
}

func TestUnmarshalUint8(t *testing.T) {
	var i Uint8
	err := json.Unmarshal(uint8JSON, &i)
	maybePanic(err)
	assertUint8(t, i, "uint8 json")

	var null Uint8
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint8(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Uint8
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint8(t, badType, "wrong type json")

	var invalid Uint8
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint8(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber8(t *testing.T) {
	var i Uint8
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to uint8")
	}
}

func TestUnmarshalUint8Overflow(t *testing.T) {
	uint8Overflow := int64(math.MaxUint8)

	// Max uint8 should decode successfully
	var i Uint8
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(uint8Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	uint8Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(uint8Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows uint8")
	}
}

func TestTextUnmarshalUint8(t *testing.T) {
	var i Uint8
	err := i.UnmarshalText([]byte("254"))
	maybePanic(err)
	assertUint8(t, i, "UnmarshalText() uint8")

	var blank Uint8
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint8(t, blank, "UnmarshalText() empty uint8")
}

func TestMarshalUint8(t *testing.T) {
	i := Uint8From(254)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "254", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint8(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "NULL", "null json marshal")
}

func TestMarshalUint8Text(t *testing.T) {
	i := Uint8From(254)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "254", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint8(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint8Pointer(t *testing.T) {
	i := Uint8From(254)
	ptr := i.Ptr()
	if *ptr != 254 {
		t.Errorf("bad %s uint8: %#v ≠ %d\n", "pointer", ptr, 254)
	}

	null := NewUint8(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint8: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint8IsZero(t *testing.T) {
	i := Uint8From(254)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint8(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint8(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint8SetValid(t *testing.T) {
	change := NewUint8(0, false)
	assertNullUint8(t, change, "SetValid()")
	change.SetValid(254)
	assertUint8(t, change, "SetValid()")
}

func TestUint8Scan(t *testing.T) {
	var i Uint8
	err := i.Scan(254)
	maybePanic(err)
	assertUint8(t, i, "scanned uint8")

	var null Uint8
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint8(t, null, "scanned null")
}

func assertUint8(t *testing.T, i Uint8, from string) {
	if i.Uint8 != 254 {
		t.Errorf("bad %s uint8: %d ≠ %d\n", from, i.Uint8, 254)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint8(t *testing.T, i Uint8, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewUint8(t *testing.T) {
	type args struct {
		i     uint8
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUint8(tt.args.i, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_IsValid(t *testing.T) {
	tests := []struct {
		name string
		u    Uint8
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsValid(); got != tt.want {
				t.Errorf("Uint8.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_IsSet(t *testing.T) {
	tests := []struct {
		name string
		u    Uint8
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsSet(); got != tt.want {
				t.Errorf("Uint8.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		u       *Uint8
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Uint8.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint8_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		u       *Uint8
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Uint8.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint8_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint8
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint8.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint8
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint8.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_SetValid(t *testing.T) {
	type args struct {
		n uint8
	}
	tests := []struct {
		name string
		u    *Uint8
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.SetValid(tt.args.n)
		})
	}
}

func TestUint8_Ptr(t *testing.T) {
	tests := []struct {
		name string
		u    Uint8
		want *uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint8.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_IsZero(t *testing.T) {
	tests := []struct {
		name string
		u    Uint8
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsZero(); got != tt.want {
				t.Errorf("Uint8.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		u       *Uint8
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Uint8.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint8_Value(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint8
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint8.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		u    Uint8
		want uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ValueOrDefault(); got != tt.want {
				t.Errorf("Uint8.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Uint8
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Uint8.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
