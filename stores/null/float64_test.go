package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

var (
	float64JSON = []byte(`1.2345`)
)

func TestFloat64From(t *testing.T) {
	f := Float64From(1.2345)
	assertFloat64(t, f, "Float64From()")

	zero := Float64From(0)
	if !zero.Valid {
		t.Error("Float64From(0)", "is invalid, but should be valid")
	}
}

func TestFloat64FromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := Float64FromPtr(iptr)
	assertFloat64(t, f, "Float64FromPtr()")

	null := Float64FromPtr(nil)
	assertNullFloat64(t, null, "Float64FromPtr(nil)")
}

func TestUnmarshalFloat64(t *testing.T) {
	var f Float64
	err := json.Unmarshal(float64JSON, &f)
	maybePanic(err)
	assertFloat64(t, f, "float64 json")

	var null Float64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat64(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType Float64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat64(t, badType, "wrong type json")

	var invalid Float64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat64(t *testing.T) {
	var f Float64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float64")

	var blank Float64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float64")
}

func TestMarshalFloat64Text(t *testing.T) {
	f := Float64From(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloat64Pointer(t *testing.T) {
	f := Float64From(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat64IsZero(t *testing.T) {
	f := Float64From(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloat64SetValid(t *testing.T) {
	change := NewFloat64(0, false)
	assertNullFloat64(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat64(t, change, "SetValid()")
}

func TestFloat64Scan(t *testing.T) {
	var f Float64
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat64(t, f, "scanned float64")

	var null Float64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat64(t, null, "scanned null")
}

func assertFloat64(t *testing.T, f Float64, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float64: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat64(t *testing.T, f Float64, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewFloat64(t *testing.T) {
	type args struct {
		f     float64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFloat64(tt.args.f, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_IsValid(t *testing.T) {
	tests := []struct {
		name string
		f    Float64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsValid(); got != tt.want {
				t.Errorf("Float64.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_IsSet(t *testing.T) {
	tests := []struct {
		name string
		f    Float64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsSet(); got != tt.want {
				t.Errorf("Float64.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		f       *Float64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Float64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat64_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		f       *Float64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Float64.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		f       Float64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		f       Float64
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_SetValid(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name string
		f    *Float64
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

func TestFloat64_Ptr(t *testing.T) {
	tests := []struct {
		name string
		f    Float64
		want *float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_IsZero(t *testing.T) {
	tests := []struct {
		name string
		f    Float64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsZero(); got != tt.want {
				t.Errorf("Float64.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		f       *Float64
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Float64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloat64_Value(t *testing.T) {
	tests := []struct {
		name    string
		f       Float64
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Float64
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); got != tt.want {
				t.Errorf("Float64.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64_Result(t *testing.T) {
	tests := []struct {
		name string
		a    Float64
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Result(); got != tt.want {
				t.Errorf("Float64.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
