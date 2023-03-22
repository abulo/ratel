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
	int64JSON = []byte(`9223372036854775806`)
)

func TestInt64From(t *testing.T) {
	i := Int64From(9223372036854775806)
	assertInt64(t, i, "Int64From()")

	zero := Int64From(0)
	if !zero.Valid {
		t.Error("Int64From(0)", "is invalid, but should be valid")
	}
}

func TestInt64FromPtr(t *testing.T) {
	n := int64(9223372036854775806)
	iptr := &n
	i := Int64FromPtr(iptr)
	assertInt64(t, i, "Int64FromPtr()")

	null := Int64FromPtr(nil)
	assertNullInt64(t, null, "Int64FromPtr(nil)")
}

func TestUnmarshalInt64(t *testing.T) {
	var i Int64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertInt64(t, i, "int64 json")

	var null Int64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt64(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Int64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt64(t, badType, "wrong type json")

	var invalid Int64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt64(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber64(t *testing.T) {
	var i Int64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int64
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int64Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int64Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt64(t *testing.T) {
	var i Int64
	err := i.UnmarshalText([]byte("9223372036854775806"))
	maybePanic(err)
	assertInt64(t, i, "UnmarshalText() int64")

	var blank Int64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt64(t, blank, "UnmarshalText() empty int64")
}

func TestMarshalInt64(t *testing.T) {
	i := Int64From(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "NULL", "null json marshal")
}

func TestMarshalInt64Text(t *testing.T) {
	i := Int64From(9223372036854775806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt64Pointer(t *testing.T) {
	i := Int64From(9223372036854775806)
	ptr := i.Ptr()
	if *ptr != 9223372036854775806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9223372036854775806)
	}

	null := NewInt64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt64IsZero(t *testing.T) {
	i := Int64From(9223372036854775806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt64SetValid(t *testing.T) {
	change := NewInt64(0, false)
	assertNullInt64(t, change, "SetValid()")
	change.SetValid(9223372036854775806)
	assertInt64(t, change, "SetValid()")
}

func TestInt64Scan(t *testing.T) {
	var i Int64
	err := i.Scan(9223372036854775806)
	maybePanic(err)
	assertInt64(t, i, "scanned int64")

	var null Int64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt64(t, null, "scanned null")
}

func assertInt64(t *testing.T, i Int64, from string) {
	if i.Int64 != 9223372036854775806 {
		t.Errorf("bad %s int64: %d ≠ %d\n", from, i.Int64, 9223372036854775806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt64(t *testing.T, i Int64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewInt64(t *testing.T) {
	type args struct {
		i     int64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt64(tt.args.i, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_IsValid(t *testing.T) {
	tests := []struct {
		name string
		i    Int64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsValid(); got != tt.want {
				t.Errorf("Int64.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_IsSet(t *testing.T) {
	tests := []struct {
		name string
		i    Int64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsSet(); got != tt.want {
				t.Errorf("Int64.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		i       *Int64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Int64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		i       *Int64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Int64.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		i       Int64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		i       Int64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_SetValid(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		i    *Int64
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

func TestInt64_Ptr(t *testing.T) {
	tests := []struct {
		name string
		i    Int64
		want *int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_IsZero(t *testing.T) {
	tests := []struct {
		name string
		i    Int64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsZero(); got != tt.want {
				t.Errorf("Int64.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		i       *Int64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Int64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64_Value(t *testing.T) {
	tests := []struct {
		name    string
		i       Int64
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		i    Int64
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.ValueOrDefault(); got != tt.want {
				t.Errorf("Int64.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Result(t *testing.T) {
	tests := []struct {
		name string
		a    Int64
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Result(); got != tt.want {
				t.Errorf("Int64.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
