package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

func TestByteFrom(t *testing.T) {
	i := ByteFrom('b')
	assertByte(t, i, "ByteFrom()")

	zero := ByteFrom(0)
	if !zero.Valid {
		t.Error("ByteFrom(0)", "is invalid, but should be valid")
	}
}

func TestByteFromPtr(t *testing.T) {
	n := byte('b')
	iptr := &n
	i := ByteFromPtr(iptr)
	assertByte(t, i, "ByteFromPtr()")

	null := ByteFromPtr(nil)
	assertNullByte(t, null, "ByteFromPtr(nil)")
}

func TestUnmarshalNonByteegerNumber(t *testing.T) {
	var i Byte
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestTextUnmarshalByte(t *testing.T) {
	var i Byte
	err := i.UnmarshalText([]byte("b"))
	maybePanic(err)
	assertByte(t, i, "UnmarshalText() int")

	var blank Byte
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullByte(t, blank, "UnmarshalText() empty int")
}

func TestMarshalByte(t *testing.T) {
	i := ByteFrom('b')
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"b"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewByte(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "NULL", "null json marshal")
}

func TestMarshalByteText(t *testing.T) {
	i := ByteFrom('b')
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "b", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewByte(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBytePointer(t *testing.T) {
	i := ByteFrom('b')
	ptr := i.Ptr()
	if *ptr != 'b' {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 'b')
	}

	null := NewByte(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestByteIsZero(t *testing.T) {
	i := ByteFrom('b')
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewByte(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewByte(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestByteSetValid(t *testing.T) {
	change := NewByte(0, false)
	assertNullByte(t, change, "SetValid()")
	change.SetValid('b')
	assertByte(t, change, "SetValid()")
}

func TestByteScan(t *testing.T) {
	var i Byte
	err := i.Scan("b")
	maybePanic(err)
	assertByte(t, i, "scanned int")

	var null Byte
	err = null.Scan(nil)
	maybePanic(err)
	assertNullByte(t, null, "scanned null")
}

func assertByte(t *testing.T, i Byte, from string) {
	if i.Byte != 'b' {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Byte, 'b')
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullByte(t *testing.T, i Byte, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewByte(t *testing.T) {
	type args struct {
		b     byte
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewByte(tt.args.b, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_IsValid(t *testing.T) {
	tests := []struct {
		name string
		b    Byte
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsValid(); got != tt.want {
				t.Errorf("Byte.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_IsSet(t *testing.T) {
	tests := []struct {
		name string
		b    Byte
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsSet(); got != tt.want {
				t.Errorf("Byte.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *Byte
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Byte.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestByte_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		b       *Byte
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Byte.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestByte_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       Byte
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Byte.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Byte.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		b       Byte
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Byte.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Byte.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_SetValid(t *testing.T) {
	type args struct {
		n byte
	}
	tests := []struct {
		name string
		b    *Byte
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetValid(tt.args.n)
		})
	}
}

func TestByte_Ptr(t *testing.T) {
	tests := []struct {
		name string
		b    Byte
		want *byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Byte.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_IsZero(t *testing.T) {
	tests := []struct {
		name string
		b    Byte
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsZero(); got != tt.want {
				t.Errorf("Byte.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		b       *Byte
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Byte.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestByte_Value(t *testing.T) {
	tests := []struct {
		name    string
		b       Byte
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Byte.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Byte.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Byte
		want byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); got != tt.want {
				t.Errorf("Byte.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByte_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Byte
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Byte.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
