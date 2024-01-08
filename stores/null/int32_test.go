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
	int32JSON = []byte(`2147483646`)
)

func TestInt32From(t *testing.T) {
	i := Int32From(2147483646)
	assertInt32(t, i, "Int32From()")

	zero := Int32From(0)
	if !zero.Valid {
		t.Error("Int32From(0)", "is invalid, but should be valid")
	}
}

func TestInt32FromPtr(t *testing.T) {
	n := int32(2147483646)
	iptr := &n
	i := Int32FromPtr(iptr)
	assertInt32(t, i, "Int32FromPtr()")

	null := Int32FromPtr(nil)
	assertNullInt32(t, null, "Int32FromPtr(nil)")
}

func TestUnmarshalInt32(t *testing.T) {
	var i Int32
	err := json.Unmarshal(int32JSON, &i)
	maybePanic(err)
	assertInt32(t, i, "int32 json")

	var null Int32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt32(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Int32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt32(t, badType, "wrong type json")

	var invalid Int32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt32(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber32(t *testing.T) {
	var i Int32
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int32")
	}
}

func TestUnmarshalInt32Overflow(t *testing.T) {
	int32Overflow := uint32(math.MaxInt32)

	// Max int32 should decode successfully
	var i Int32
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int32Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int32")
	}
}

func TestTextUnmarshalInt32(t *testing.T) {
	var i Int32
	err := i.UnmarshalText([]byte("2147483646"))
	maybePanic(err)
	assertInt32(t, i, "UnmarshalText() int32")

	var blank Int32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt32(t, blank, "UnmarshalText() empty int32")
}

func TestMarshalInt32Text(t *testing.T) {
	i := Int32From(2147483646)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "2147483646", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt32(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt32Pointer(t *testing.T) {
	i := Int32From(2147483646)
	ptr := i.Ptr()
	if *ptr != 2147483646 {
		t.Errorf("bad %s int32: %#v ≠ %d\n", "pointer", ptr, 2147483646)
	}

	null := NewInt32(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt32IsZero(t *testing.T) {
	i := Int32From(2147483646)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt32(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt32(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt32SetValid(t *testing.T) {
	change := NewInt32(0, false)
	assertNullInt32(t, change, "SetValid()")
	change.SetValid(2147483646)
	assertInt32(t, change, "SetValid()")
}

func TestInt32Scan(t *testing.T) {
	var i Int32
	err := i.Scan(2147483646)
	maybePanic(err)
	assertInt32(t, i, "scanned int32")

	var null Int32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt32(t, null, "scanned null")
}

func assertInt32(t *testing.T, i Int32, from string) {
	if i.Int32 != 2147483646 {
		t.Errorf("bad %s int32: %d ≠ %d\n", from, i.Int32, 2147483646)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt32(t *testing.T, i Int32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewInt32(t *testing.T) {
	type args struct {
		i     int32
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt32(tt.args.i, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_IsValid(t *testing.T) {
	tests := []struct {
		name string
		i    Int32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsValid(); got != tt.want {
				t.Errorf("Int32.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_IsSet(t *testing.T) {
	tests := []struct {
		name string
		i    Int32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsSet(); got != tt.want {
				t.Errorf("Int32.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		i       *Int32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Int32.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt32_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		i       *Int32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Int32.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt32_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		i       Int32
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		i       Int32
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_SetValid(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		i    *Int32
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

func TestInt32_Ptr(t *testing.T) {
	tests := []struct {
		name string
		i    Int32
		want *int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_IsZero(t *testing.T) {
	tests := []struct {
		name string
		i    Int32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsZero(); got != tt.want {
				t.Errorf("Int32.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		i       *Int32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Int32.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt32_Value(t *testing.T) {
	tests := []struct {
		name    string
		i       Int32
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		i    Int32
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.ValueOrDefault(); got != tt.want {
				t.Errorf("Int32.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Result(t *testing.T) {
	tests := []struct {
		name string
		a    Int32
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Result(); got != tt.want {
				t.Errorf("Int32.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
