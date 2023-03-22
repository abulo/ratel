package null

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

var (
	boolJSON = []byte(`true`)
)

func TestBoolFrom(t *testing.T) {
	b := BoolFrom(true)
	assertBool(t, b, "BoolFrom()")

	zero := BoolFrom(false)
	if !zero.Valid {
		t.Error("BoolFrom(false)", "is invalid, but should be valid")
	}
}

func TestBoolFromPtr(t *testing.T) {
	n := true
	bptr := &n
	b := BoolFromPtr(bptr)
	assertBool(t, b, "BoolFromPtr()")

	null := BoolFromPtr(nil)
	assertNullBool(t, null, "BoolFromPtr(nil)")
}

func TestTextUnmarshalBool(t *testing.T) {
	var b Bool
	err := b.UnmarshalText([]byte("true"))
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero Bool
	err = zero.UnmarshalText([]byte("false"))
	maybePanic(err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank Bool
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var invalid Bool
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, invalid, "invalid json")
}

func TestMarshalBoolText(t *testing.T) {
	b := BoolFrom(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := NewBool(false, true)
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := NewBool(false, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBoolPointer(t *testing.T) {
	b := BoolFrom(true)
	ptr := b.Ptr()
	if *ptr != true {
		t.Errorf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	null := NewBool(false, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolIsZero(t *testing.T) {
	b := BoolFrom(true)
	if b.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewBool(false, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewBool(false, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBoolSetValid(t *testing.T) {
	change := NewBool(false, false)
	assertNullBool(t, change, "SetValid()")
	change.SetValid(true)
	assertBool(t, change, "SetValid()")
}

func TestBoolScan(t *testing.T) {
	var b Bool
	err := b.Scan(true)
	maybePanic(err)
	assertBool(t, b, "scanned bool")

	var null Bool
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBool(t, null, "scanned null")
}

func assertBool(t *testing.T, b Bool, from string) {
	if b.Bool != true {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b Bool, from string) {
	if b.Bool != false {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBool(t *testing.T, b Bool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewBool(t *testing.T) {
	type args struct {
		b     bool
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBool(tt.args.b, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_IsValid(t *testing.T) {
	tests := []struct {
		name string
		b    Bool
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsValid(); got != tt.want {
				t.Errorf("Bool.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_IsSet(t *testing.T) {
	tests := []struct {
		name string
		b    Bool
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsSet(); got != tt.want {
				t.Errorf("Bool.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *Bool
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBool_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		b       *Bool
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       Bool
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		b       Bool
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_SetValid(t *testing.T) {
	type args struct {
		v bool
	}
	tests := []struct {
		name string
		b    *Bool
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetValid(tt.args.v)
		})
	}
}

func TestBool_Ptr(t *testing.T) {
	tests := []struct {
		name string
		b    Bool
		want *bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_IsZero(t *testing.T) {
	tests := []struct {
		name string
		b    Bool
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsZero(); got != tt.want {
				t.Errorf("Bool.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		b       *Bool
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Bool.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBool_Value(t *testing.T) {
	tests := []struct {
		name    string
		b       Bool
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Bool
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); got != tt.want {
				t.Errorf("Bool.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Bool
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Bool.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
