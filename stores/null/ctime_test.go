package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var (
	ctimeString    = "21:54:21"
	ctimeJSON      = []byte(`"` + ctimeString + `"`)
	nullCTimeJSON  = []byte(`null`)
	ctimeValue, _  = time.Parse(RFC3339TimeOnly, ctimeString)
	badCTimeObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalCTimeText(t *testing.T) {
	ct := CTimeFrom(ctimeValue)
	txt, err := ct.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, ctimeString, "marshal text")

	var unmarshal CTime
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertCTime(t, unmarshal, "unmarshal text")

	var invalid CTime
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullCTime(t, invalid, "bad string")
}

func TestMarshalCTime(t *testing.T) {
	ct := CTimeFrom(ctimeValue)
	data, err := json.Marshal(ct)
	maybePanic(err)
	assertJSONEquals(t, data, string(ctimeJSON), "non-empty json marshal")

	ct.Valid = false
	data, err = json.Marshal(ct)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestCTimeFrom(t *testing.T) {
	ct := CTimeFrom(ctimeValue)
	assertCTime(t, ct, "DateFrom() time.Time")
}

func TestCTimeFromPtr(t *testing.T) {
	ct := CTimeFromPtr(&ctimeValue)
	assertCTime(t, ct, "DateFromPtr() time")

	null := CTimeFromPtr(nil)
	assertNullCTime(t, null, "DateFromPtr(nil)")
}

func TestCTimeSetValid(t *testing.T) {
	var ct time.Time
	change := NewCTime(ct, false)
	assertNullCTime(t, change, "SetValid()")
	change.SetValid(ctimeValue)
	assertCTime(t, change, "SetValid()")
}

func TestCTimePointer(t *testing.T) {
	ct := CTimeFrom(ctimeValue)
	ptr := ct.Ptr()
	if *ptr != ctimeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, ctimeValue)
	}

	var nt time.Time
	null := NewCTime(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestCTimeIsZero(t *testing.T) {
	ct := CTimeFrom(time.Now())
	if ct.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := CTimeFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestCTimeScanValue(t *testing.T) {
	var ct CTime
	err := ct.Scan(ctimeValue)
	maybePanic(err)
	assertCTime(t, ct, "scanned time")
	if v, err := ct.Value(); v != ctimeValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null CTime
	err = null.Scan(nil)
	maybePanic(err)
	assertNullCTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong CTime
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullCTime(t, wrong, "scanned wrong")
}

func assertCTime(t *testing.T, ct CTime, from string) {
	if ct.CTime != ctimeValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ct.CTime, ctimeValue)
	}
	if !ct.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullCTime(t *testing.T, ct CTime, from string) {
	if ct.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewCTime(t *testing.T) {
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want CTime
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCTime(tt.args.t, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_IsValid(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsValid(); got != tt.want {
				t.Errorf("CTime.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_IsSet(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsSet(); got != tt.want {
				t.Errorf("CTime.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      CTime
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("CTime.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CTime.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		tr      *CTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCTime_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		tr      CTime
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("CTime.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CTime.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		tr      *CTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("CTime.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCTime_SetValid(t *testing.T) {
	type args struct {
		v time.Time
	}
	tests := []struct {
		name string
		tr   *CTime
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetValid(tt.args.v)
		})
	}
}

func TestCTime_Ptr(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want *time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CTime.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_IsZero(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsZero(); got != tt.want {
				t.Errorf("CTime.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		tr      *CTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CTime.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCTime_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      CTime
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("CTime.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CTime.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CTime.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCTime_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   CTime
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("CTime.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
