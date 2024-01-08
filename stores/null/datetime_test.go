package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var (
	datetimeString    = "2012-12-21 21:30:12"
	datetimeJSON      = []byte(`"` + datetimeString + `"`)
	nullDateTimeJSON  = []byte(`null`)
	datetimeValue, _  = time.Parse(DateTimeSQL, datetimeString)
	badDateTimeObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalDateTimeText(t *testing.T) {
	dt := DateTimeFrom(datetimeValue)
	txt, err := dt.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, datetimeString, "marshal text")

	var unmarshal DateTime
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertDateTime(t, unmarshal, "unmarshal text")

	var invalid DateTime
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDateTime(t, invalid, "bad string")
}

func TestMarshalDateTime(t *testing.T) {
	dt := DateTimeFrom(datetimeValue)
	data, err := json.Marshal(dt)
	maybePanic(err)
	assertJSONEquals(t, data, string(datetimeJSON), "non-empty json marshal")

	dt.Valid = false
	data, err = json.Marshal(dt)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestDateTimeFrom(t *testing.T) {
	dt := DateTimeFrom(datetimeValue)
	assertDateTime(t, dt, "DateFrom() time.Time")
}

func TestDateTimeFromPtr(t *testing.T) {
	dt := DateTimeFromPtr(&datetimeValue)
	assertDateTime(t, dt, "DateFromPtr() time")

	null := DateTimeFromPtr(nil)
	assertNullDateTime(t, null, "DateFromPtr(nil)")
}

func TestDateTimeSetValid(t *testing.T) {
	var dt time.Time
	change := NewDateTime(dt, false)
	assertNullDateTime(t, change, "SetValid()")
	change.SetValid(datetimeValue)
	assertDateTime(t, change, "SetValid()")
}

func TestDateTimePointer(t *testing.T) {
	dt := DateTimeFrom(datetimeValue)
	ptr := dt.Ptr()
	if *ptr != datetimeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, datetimeValue)
	}

	var nt time.Time
	null := NewDateTime(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestDateTimeIsZero(t *testing.T) {
	dt := DateTimeFrom(time.Now())
	if dt.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := DateTimeFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestDateTimeScanValue(t *testing.T) {
	var dt DateTime
	err := dt.Scan(datetimeValue)
	maybePanic(err)
	assertDateTime(t, dt, "scanned time")
	if v, err := dt.Value(); v != datetimeValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null DateTime
	err = null.Scan(nil)
	maybePanic(err)
	assertNullDateTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong DateTime
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDateTime(t, wrong, "scanned wrong")
}

func assertDateTime(t *testing.T, dt DateTime, from string) {
	if dt.DateTime != datetimeValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, dt.DateTime, datetimeValue)
	}
	if !dt.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullDateTime(t *testing.T, dt DateTime, from string) {
	if dt.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewDateTime(t *testing.T) {
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want DateTime
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDateTime(tt.args.t, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_IsValid(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsValid(); got != tt.want {
				t.Errorf("DateTime.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_IsSet(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsSet(); got != tt.want {
				t.Errorf("DateTime.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      DateTime
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateTime.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		tr      *DateTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateTime_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		tr      DateTime
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateTime.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		tr      *DateTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateTime_SetValid(t *testing.T) {
	type args struct {
		v time.Time
	}
	tests := []struct {
		name string
		tr   *DateTime
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

func TestDateTime_Ptr(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want *time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_IsZero(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsZero(); got != tt.want {
				t.Errorf("DateTime.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		tr      *DateTime
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateTime_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      DateTime
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateTime.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   DateTime
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("DateTime.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
