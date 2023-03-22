package null

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

var (
	dateString    = "2012-12-21"
	dateJSON      = []byte(`"` + dateString + `"`)
	nullDateJSON  = []byte(`null`)
	dateValue, _  = time.Parse(RFC3339DateOnly, dateString)
	badDateObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalDateText(t *testing.T) {
	dt := DateFrom(dateValue)
	txt, err := dt.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, dateString, "marshal text")

	var unmarshal Date
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertDate(t, unmarshal, "unmarshal text")

	var invalid Date
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDate(t, invalid, "bad string")
}

func TestDateFrom(t *testing.T) {
	dt := DateFrom(dateValue)
	assertDate(t, dt, "DateFrom() time.Time")
}

func TestDateFromPtr(t *testing.T) {
	dt := DateFromPtr(&dateValue)
	assertDate(t, dt, "DateFromPtr() time")

	null := DateFromPtr(nil)
	assertNullDate(t, null, "DateFromPtr(nil)")
}

func TestDateSetValid(t *testing.T) {
	var dt time.Time
	change := NewDate(dt, false)
	assertNullDate(t, change, "SetValid()")
	change.SetValid(dateValue)
	assertDate(t, change, "SetValid()")
}

func TestDatePointer(t *testing.T) {
	dt := DateFrom(dateValue)
	ptr := dt.Ptr()
	if *ptr != dateValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, dateValue)
	}

	var nt time.Time
	null := NewDate(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestDateIsZero(t *testing.T) {
	dt := DateFrom(time.Now())
	if dt.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := DateFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestDateScanValue(t *testing.T) {
	var dt Date
	err := dt.Scan(dateValue)
	maybePanic(err)
	assertDate(t, dt, "scanned time")
	if v, err := dt.Value(); v != dateValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Date
	err = null.Scan(nil)
	maybePanic(err)
	assertNullDate(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Date
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullDate(t, wrong, "scanned wrong")
}

func assertDate(t *testing.T, dt Date, from string) {
	if dt.Date != dateValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, dt.Date, dateValue)
	}
	if !dt.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullDate(t *testing.T, dt Date, from string) {
	if dt.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewDate(t *testing.T) {
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDate(tt.args.t, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsValid(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsValid(); got != tt.want {
				t.Errorf("Date.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsSet(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsSet(); got != tt.want {
				t.Errorf("Date.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      Date
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		tr      *Date
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		tr      Date
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		tr      *Date
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_SetValid(t *testing.T) {
	type args struct {
		v time.Time
	}
	tests := []struct {
		name string
		tr   *Date
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

func TestDate_Ptr(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want *time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsZero(); got != tt.want {
				t.Errorf("Date.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		tr      *Date
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Date.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      Date
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   Date
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("Date.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
