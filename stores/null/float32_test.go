package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

var (
	float32JSON = []byte(`1.2345`)
)

func TestFloat32From(t *testing.T) {
	f := Float32From(1.2345)
	assertFloat32(t, f, "Float32From()")

	zero := Float32From(0)
	if !zero.Valid {
		t.Error("Float32From(0)", "is invalid, but should be valid")
	}
}

func TestFloat32FromPtr(t *testing.T) {
	n := float32(1.2345)
	iptr := &n
	f := Float32FromPtr(iptr)
	assertFloat32(t, f, "Float32FromPtr()")

	null := Float32FromPtr(nil)
	assertNullFloat32(t, null, "Float32FromPtr(nil)")
}

func TestUnmarshalFloat32(t *testing.T) {
	var f Float32
	err := json.Unmarshal(float32JSON, &f)
	maybePanic(err)
	assertFloat32(t, f, "float32 json")

	var null Float32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat32(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Float32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat32(t, badType, "wrong type json")

	var invalid Float32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat32(t *testing.T) {
	var f Float32
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat32(t, f, "UnmarshalText() float32")

	var blank Float32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat32(t, blank, "UnmarshalText() empty float32")
}

func TestMarshalFloat32Text(t *testing.T) {
	f := Float32From(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat32(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloat32Pointer(t *testing.T) {
	f := Float32From(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float32: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat32(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat32IsZero(t *testing.T) {
	f := Float32From(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat32(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat32(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloat32SetValid(t *testing.T) {
	change := NewFloat32(0, false)
	assertNullFloat32(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat32(t, change, "SetValid()")
}

func TestFloat32Scan(t *testing.T) {
	var f Float32
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat32(t, f, "scanned float32")

	var null Float32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat32(t, null, "scanned null")
}

func assertFloat32(t *testing.T, f Float32, from string) {
	if f.Float32 != 1.2345 {
		t.Errorf("bad %s float32: %f ≠ %f\n", from, f.Float32, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat32(t *testing.T, f Float32, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewFloat32(t *testing.T) {
	type args struct {
		f     float32
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFloat32(tt.args.f, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_IsValid(t *testing.T) {
	tests := []struct {
		name string
		f    Float32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsValid(); got != tt.want {
				t.Errorf("Float32.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_IsSet(t *testing.T) {
	tests := []struct {
		name string
		f    Float32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsSet(); got != tt.want {
				t.Errorf("Float32.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		f       *Float32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Float32.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat32_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		f       *Float32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Float32.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat32_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		f       Float32
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float32.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float32.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		f       Float32
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float32.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float32.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_SetValid(t *testing.T) {
	type args struct {
		n float32
	}
	tests := []struct {
		name string
		f    *Float32
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.SetValid(tt.args.n)
		})
	}
}

func TestFloat32_Ptr(t *testing.T) {
	tests := []struct {
		name string
		f    Float32
		want *float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float32.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_IsZero(t *testing.T) {
	tests := []struct {
		name string
		f    Float32
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsZero(); got != tt.want {
				t.Errorf("Float32.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		f       *Float32
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Float32.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat32_Value(t *testing.T) {
	tests := []struct {
		name    string
		f       Float32
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float32.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float32.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Float32
		want float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); got != tt.want {
				t.Errorf("Float32.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32_Result(t *testing.T) {
	tests := []struct {
		name string
		a    Float32
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Result(); got != tt.want {
				t.Errorf("Float32.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
