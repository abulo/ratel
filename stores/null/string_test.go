package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assertStr(t, str, "StringFrom() string")

	zero := StringFrom("")
	if !zero.Valid {
		t.Error("StringFrom(0)", "is invalid, but should be valid")
	}
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	sptr := &s
	str := StringFromPtr(sptr)
	assertStr(t, str, "StringFromPtr() string")

	null := StringFromPtr(nil)
	assertNullStr(t, null, "StringFromPtr(nil)")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assertStr(t, str, "string json")

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullStr(t, null, "null json")
	if !null.Set {
		t.Error("should be Set")
	}

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := StringFrom("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	null := StringFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, `null`, "null json marshal")
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")
}

// Tests omitempty... broken until Go 1.4
// type stringInStruct struct {
// 	Test String `json:"test,omitempty"`
// }
// func TestMarshalStringInStruct(t *testing.T) {
// 	obj := stringInStruct{Test: StringFrom("")}
// 	data, err := json.Marshal(obj)
// 	maybePanic(err)
// 	assertJSONEquals(t, data, `{}`, "null string in struct")
// }

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := NewString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := StringFrom("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := NewString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestStringSetValid(t *testing.T) {
	change := NewString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStr(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func TestNewString(t *testing.T) {
	type args struct {
		s     string
		valid bool
	}
	tests := []struct {
		name string
		args args
		want String
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewString(tt.args.s, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("String.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_IsSet(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsSet(); got != tt.want {
				t.Errorf("String.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       *String
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("String.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		s       *String
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("String.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestString_SetValid(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		s    *String
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetValid(tt.args.v)
		})
	}
}

func TestString_Ptr(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_IsZero(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsZero(); got != tt.want {
				t.Errorf("String.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		s       *String
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("String.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestString_Value(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   String
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); got != tt.want {
				t.Errorf("String.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   String
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("String.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
