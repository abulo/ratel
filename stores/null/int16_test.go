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
	int16JSON = []byte(`32766`)
)

func TestInt16From(t *testing.T) {
	i := Int16From(32766)
	assertInt16(t, i, "Int16From()")

	zero := Int16From(0)
	if !zero.Valid {
		t.Error("Int16From(0)", "is invalid, but should be valid")
	}
}

func TestInt16FromPtr(t *testing.T) {
	n := int16(32766)
	iptr := &n
	i := Int16FromPtr(iptr)
	assertInt16(t, i, "Int16FromPtr()")

	null := Int16FromPtr(nil)
	assertNullInt16(t, null, "Int16FromPtr(nil)")
}

func TestUnmarshalInt16(t *testing.T) {
	var i Int16
	err := json.Unmarshal(int16JSON, &i)
	maybePanic(err)
	assertInt16(t, i, "int16 json")

	var null Int16
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt16(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Int16
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt16(t, badType, "wrong type json")

	var invalid Int16
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt16(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber16(t *testing.T) {
	var i Int16
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int16")
	}
}

func TestUnmarshalInt16Overflow(t *testing.T) {
	int16Overflow := uint16(math.MaxInt16)

	// Max int16 should decode successfully
	var i Int16
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int16Overflow), 10)), &i)
	maybePanic(err)
	// Attempt to overflow
	int16Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int16Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int16")
	}
}

func TestTextUnmarshalInt16(t *testing.T) {
	var i Int16
	err := i.UnmarshalText([]byte("32766"))
	maybePanic(err)
	assertInt16(t, i, "UnmarshalText() int16")

	var blank Int16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt16(t, blank, "UnmarshalText() empty int16")
}

func TestMarshalInt16(t *testing.T) {
	i := Int16From(32766)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "32766", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt16(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "NULL", "null json marshal")
}

func TestMarshalInt16Text(t *testing.T) {
	i := Int16From(32766)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "32766", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt16(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt16Pointer(t *testing.T) {
	i := Int16From(32766)
	ptr := i.Ptr()
	if *ptr != 32766 {
		t.Errorf("bad %s int16: %#v ≠ %d\n", "pointer", ptr, 32766)
	}

	null := NewInt16(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int16: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt16IsZero(t *testing.T) {
	i := Int16From(32766)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt16(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt16(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt16SetValid(t *testing.T) {
	change := NewInt16(0, false)
	assertNullInt16(t, change, "SetValid()")
	change.SetValid(32766)
	assertInt16(t, change, "SetValid()")
}

func TestInt16Scan(t *testing.T) {
	var i Int16
	err := i.Scan(32766)
	maybePanic(err)
	assertInt16(t, i, "scanned int16")

	var null Int16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt16(t, null, "scanned null")
}

func assertInt16(t *testing.T, i Int16, from string) {
	if i.Int16 != 32766 {
		t.Errorf("bad %s int16: %d ≠ %d\n", from, i.Int16, 32766)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt16(t *testing.T, i Int16, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewInt16(t *testing.T) {
	type args struct {
		i     int16
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt16(tt.args.i, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_IsValid(t *testing.T) {
	tests := []struct {
		name string
		i    Int16
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsValid(); got != tt.want {
				t.Errorf("Int16.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_IsSet(t *testing.T) {
	tests := []struct {
		name string
		i    Int16
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsSet(); got != tt.want {
				t.Errorf("Int16.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		i       *Int16
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Int16.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt16_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		i       *Int16
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Int16.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt16_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		i       Int16
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int16.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		i       Int16
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int16.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_SetValid(t *testing.T) {
	type args struct {
		n int16
	}
	tests := []struct {
		name string
		i    *Int16
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.i.SetValid(tt.args.n)
		})
	}
}

func TestInt16_Ptr(t *testing.T) {
	tests := []struct {
		name string
		i    Int16
		want *int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_IsZero(t *testing.T) {
	tests := []struct {
		name string
		i    Int16
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsZero(); got != tt.want {
				t.Errorf("Int16.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		i       *Int16
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Int16.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt16_Value(t *testing.T) {
	tests := []struct {
		name    string
		i       Int16
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int16.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		i    Int16
		want int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.ValueOrDefault(); got != tt.want {
				t.Errorf("Int16.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16_Result(t *testing.T) {
	tests := []struct {
		name string
		a    Int16
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Result(); got != tt.want {
				t.Errorf("Int16.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
