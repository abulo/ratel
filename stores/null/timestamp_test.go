package null

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var (
	timestampString    = "2012-12-21 21:30:12.654321"
	timestampJSON      = []byte(`"` + timestampString + `"`)
	nullTimeStampJSON  = []byte(`null`)
	timestampValue, _  = time.Parse(DateTimeSQL, timestampString)
	badTimeStampObject = []byte(`{"hello": "world"}`)
)

func TestUnmarshalTimeStampJSON(t *testing.T) {
	var ts TimeStamp
	err := json.Unmarshal(timestampJSON, &ts)
	maybePanic(err)
	assertTimeStamp(t, ts, "UnmarshalJSON() json")

	var null TimeStamp
	err = json.Unmarshal(nullTimeStampJSON, &null)
	maybePanic(err)
	assertNullTimeStamp(t, null, "null time json")
	if !null.Set {
		t.Error("should be Set")
	}

	var invalid TimeStamp
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.ParseError, not %T", err)
	}
	assertNullTimeStamp(t, invalid, "invalid from object json")

	var bad TimeStamp
	err = json.Unmarshal(badTimeStampObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullTimeStamp(t, bad, "bad from object json")

	var wrongType TimeStamp
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullTimeStamp(t, wrongType, "wrong type object json")
}

func TestUnmarshalTimeStampText(t *testing.T) {
	ts := TimeStampFrom(timestampValue)
	txt, err := ts.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timestampString, "marshal text")

	var unmarshal TimeStamp
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTimeStamp(t, unmarshal, "unmarshal text")

	var invalid TimeStamp
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTimeStamp(t, invalid, "bad string")
}

func TestMarshalTimeStamp(t *testing.T) {
	ts := TimeStampFrom(timestampValue)
	data, err := json.Marshal(ts)
	maybePanic(err)
	assertJSONEquals(t, data, string(timestampJSON), "non-empty json marshal")

	ts.Valid = false
	data, err = json.Marshal(ts)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestTimeStampFrom(t *testing.T) {
	ts := TimeStampFrom(timestampValue)
	assertTimeStamp(t, ts, "DateFrom() time.Time")
}

func TestTimeStampFromPtr(t *testing.T) {
	ts := TimeStampFromPtr(&timestampValue)
	assertTimeStamp(t, ts, "DateFromPtr() time")

	null := TimeStampFromPtr(nil)
	assertNullTimeStamp(t, null, "DateFromPtr(nil)")
}

func TestTimeStampSetValid(t *testing.T) {
	var ts time.Time
	change := NewTimeStamp(ts, false)
	assertNullTimeStamp(t, change, "SetValid()")
	change.SetValid(timestampValue)
	assertTimeStamp(t, change, "SetValid()")
}

func TestTimeStampPointer(t *testing.T) {
	dt := TimeStampFrom(timestampValue)
	ptr := dt.Ptr()
	if *ptr != timestampValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timestampValue)
	}

	var nt time.Time
	null := NewTimeStamp(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeStampIsZero(t *testing.T) {
	dt := TimeStampFrom(time.Now())
	if dt.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := TimeStampFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestTimeStampScanValue(t *testing.T) {
	var ts TimeStamp
	err := ts.Scan(timestampValue)
	maybePanic(err)
	assertTimeStamp(t, ts, "scanned time")
	if v, err := ts.Value(); v != timestampValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null TimeStamp
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTimeStamp(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong TimeStamp
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTimeStamp(t, wrong, "scanned wrong")
}

func assertTimeStamp(t *testing.T, ts TimeStamp, from string) {
	if ts.TimeStamp != timestampValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ts.TimeStamp, timestampValue)
	}
	if !ts.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTimeStamp(t *testing.T, dt TimeStamp, from string) {
	if dt.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewTimeStamp(t *testing.T) {
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want TimeStamp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTimeStamp(tt.args.t, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_IsValid(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsValid(); got != tt.want {
				t.Errorf("TimeStamp.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_IsSet(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsSet(); got != tt.want {
				t.Errorf("TimeStamp.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      TimeStamp
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeStamp.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		tr      *TimeStamp
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeStamp_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		tr      TimeStamp
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeStamp.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		tr      *TimeStamp
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeStamp_SetValid(t *testing.T) {
	type args struct {
		v time.Time
	}
	tests := []struct {
		name string
		tr   *TimeStamp
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

func TestTimeStamp_Ptr(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want *time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeStamp.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_IsZero(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsZero(); got != tt.want {
				t.Errorf("TimeStamp.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		tr      *TimeStamp
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeStamp_Value(t *testing.T) {
	tests := []struct {
		name    string
		tr      TimeStamp
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStamp.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeStamp.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_ValueOrDefault(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ValueOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeStamp.ValueOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStamp_Result(t *testing.T) {
	tests := []struct {
		name string
		tr   TimeStamp
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Result(); got != tt.want {
				t.Errorf("TimeStamp.Result() = %v, want %v", got, tt.want)
			}
		})
	}
}
