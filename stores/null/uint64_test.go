package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

var (
	uint64JSON = []byte(`18446744073709551614`)
)

func TestUint64From(t *testing.T) {
	i := Uint64From(18446744073709551614)
	assertUint64(t, i, "Uint64From()")

	zero := Uint64From(0)
	if !zero.Valid {
		t.Error("Uint64From(0)", "is invalid, but should be valid")
	}
}

func TestUint64FromPtr(t *testing.T) {
	n := uint64(18446744073709551614)
	iptr := &n
	i := Uint64FromPtr(iptr)
	assertUint64(t, i, "Uint64FromPtr()")

	null := Uint64FromPtr(nil)
	assertNullUint64(t, null, "Uint64FromPtr(nil)")
}

func TestUnmarshalUint64(t *testing.T) {
	var i Uint64
	err := json.Unmarshal(uint64JSON, &i)
	maybePanic(err)
	assertUint64(t, i, "uint64 json")

	var null Uint64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint64(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Uint64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint64(t, badType, "wrong type json")

	var invalid Uint64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint64(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber64(t *testing.T) {
	var i Uint64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to uint64")
	}
}

func TestTextUnmarshalUint64(t *testing.T) {
	var i Uint64
	err := i.UnmarshalText([]byte("18446744073709551614"))
	maybePanic(err)
	assertUint64(t, i, "UnmarshalText() uint64")

	var blank Uint64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint64(t, blank, "UnmarshalText() empty uint64")
}

func TestMarshalUint64(t *testing.T) {
	i := Uint64From(18446744073709551614)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "18446744073709551614", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "NULL", "null json marshal")
}

func TestMarshalUint64Text(t *testing.T) {
	i := Uint64From(18446744073709551614)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "18446744073709551614", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint64Pointer(t *testing.T) {
	i := Uint64From(18446744073709551614)
	ptr := i.Ptr()
	if *ptr != 18446744073709551614 {
		t.Errorf("bad %s uint64: %#v ≠ %d\n", "pointer", ptr, uint64(18446744073709551614))
	}

	null := NewUint64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint64IsZero(t *testing.T) {
	i := Uint64From(18446744073709551614)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint64SetValid(t *testing.T) {
	change := NewUint64(0, false)
	assertNullUint64(t, change, "SetValid()")
	change.SetValid(18446744073709551614)
	assertUint64(t, change, "SetValid()")
}

func TestUint64Scan(t *testing.T) {
	var i Uint64
	err := i.Scan(uint64(18446744073709551614))
	maybePanic(err)
	assertUint64(t, i, "scanned uint64")

	err = i.Scan(int64(-2))
	maybePanic(err)
	assertUint64(t, i, "scanned int64")

	err = i.Scan(nil)
	maybePanic(err)
	assertNullUint64(t, i, "scanned null")
}

func assertUint64(t *testing.T, i Uint64, from string) {
	if i.Uint64 != 18446744073709551614 {
		t.Errorf("bad %s uint64: %d ≠ %d\n", from, i.Uint64, uint64(18446744073709551614))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint64(t *testing.T, i Uint64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewUint64(t *testing.T) {
	type args struct {
		i     uint64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUint64(tt.args.i, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_IsValid(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsValid(); got != tt.want {
				t.Errorf("Uint64.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_IsSet(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsSet(); got != tt.want {
				t.Errorf("Uint64.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		u       *Uint64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Uint64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint64_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		u       *Uint64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Uint64.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_SetValid(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name string
		u    *Uint64
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

func TestUint64_Ptr(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want *uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_IsZero(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsZero(); got != tt.want {
				t.Errorf("Uint64.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		u       *Uint64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Uint64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUint64_Value(t *testing.T) {
	tests := []struct {
		name    string
		u       Uint64
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ValueOrDefault(); got != tt.want {
				t.Errorf("Uint64.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Uint64
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Uint64.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
