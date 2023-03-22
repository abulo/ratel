package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

var (
	hello        = []byte("hello")
	bytesJSON    = []byte(`"hello"`)
	b64BytesJSON = []byte(`"aGVsbG8="`)
)

func TestBytesFrom(t *testing.T) {
	i := BytesFrom(hello)
	assertBytes(t, i, "BytesFrom()")

	zero := BytesFrom(nil)
	if zero.Valid {
		t.Error("BytesFrom(nil)", "is valid, but should be invalid")
	}

	zero = BytesFrom([]byte{})
	if !zero.Valid {
		t.Error("BytesFrom([]byte{})", "is invalid, but should be valid")
	}
}

func TestBytesFromPtr(t *testing.T) {
	n := hello
	iptr := &n
	i := BytesFromPtr(iptr)
	assertBytes(t, i, "BytesFromPtr()")

	null := BytesFromPtr(nil)
	assertNullBytes(t, null, "BytesFromPtr(nil)")
}

func TestUnmarshalBytes(t *testing.T) {
	var i Bytes
	err := json.Unmarshal(b64BytesJSON, &i)
	maybePanic(err)
	assertBytes(t, i, "[]byte json")

	var ni Bytes
	err = ni.UnmarshalJSON([]byte{})
	if err == nil {
		t.Errorf("Expected error")
	}

	var null Bytes
	err = null.UnmarshalJSON(NullBytes)
	if err != nil {
		t.Error(err)
	}
	if null.Valid {
		t.Errorf("expected Valid to be false, got true")
	}
	if null.Bytes != nil {
		t.Errorf("Expected Bytes to be nil, but was not: %#v %#v", null.Bytes, []byte(`null`))
	}
	if !null.Set {
		t.Errorf("Expected Set to be true; got false")
	}
}

func TestTextUnmarshalBytes(t *testing.T) {
	var i Bytes
	err := i.UnmarshalText(hello)
	maybePanic(err)
	assertBytes(t, i, "UnmarshalText() []byte")

	var blank Bytes
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBytes(t, blank, "UnmarshalText() empty []byte")
}

func TestMarshalBytesText(t *testing.T) {
	i := BytesFrom(bytesJSON)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewBytes(nil, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBytesPointer(t *testing.T) {
	i := BytesFrom([]byte(`"hello"`))
	ptr := i.Ptr()
	if !bytes.Equal(*ptr, bytesJSON) {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "pointer", ptr, `"hello"`)
	}

	null := NewBytes(nil, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBytesIsZero(t *testing.T) {
	i := BytesFrom(bytesJSON)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewBytes(nil, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewBytes(nil, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBytesSetValid(t *testing.T) {
	change := NewBytes(nil, false)
	assertNullBytes(t, change, "SetValid()")
	change.SetValid(hello)
	assertBytes(t, change, "SetValid()")
}

func TestBytesScan(t *testing.T) {
	var i Bytes
	err := i.Scan(`hello`)
	maybePanic(err)
	assertBytes(t, i, "Scan() []byte")

	var null Bytes
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBytes(t, null, "scanned null")
}

func assertBytes(t *testing.T, i Bytes, from string) {
	if !bytes.Equal(i.Bytes, hello) {
		t.Errorf("bad %s []byte: %v ≠ %v\n", from, string(i.Bytes), "hello")
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBytes(t *testing.T, i Bytes, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewBytes(t *testing.T) {
	type args struct {
		b     []byte
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Bytes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBytes(tt.args.b, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_IsValid(t *testing.T) {
	tests := []struct {
		name string
		b    Bytes
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsValid(); got != tt.want {
				t.Errorf("Bytes.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_IsSet(t *testing.T) {
	tests := []struct {
		name string
		b    Bytes
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsSet(); got != tt.want {
				t.Errorf("Bytes.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *Bytes
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Bytes.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytes_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		b       *Bytes
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Bytes.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytes_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       Bytes
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		b       Bytes
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_SetValid(t *testing.T) {
	type args struct {
		n []byte
	}
	tests := []struct {
		name string
		b    *Bytes
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

func TestBytes_Ptr(t *testing.T) {
	tests := []struct {
		name string
		b    Bytes
		want *[]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_IsZero(t *testing.T) {
	tests := []struct {
		name string
		b    Bytes
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsZero(); got != tt.want {
				t.Errorf("Bytes.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		b       *Bytes
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Bytes.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytes_Value(t *testing.T) {
	tests := []struct {
		name    string
		b       Bytes
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Bytes
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Bytes
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Bytes.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
