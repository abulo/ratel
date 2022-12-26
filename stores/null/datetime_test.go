package null

import (
	"encoding/json"
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

func TestUnmarshalDateTimeJSON(t *testing.T) {
	var dt DateTime
	err := json.Unmarshal(datetimeJSON, &dt)
	maybePanic(err)
	assertDateTime(t, dt, "UnmarshalJSON() json")

	var null DateTime
	err = json.Unmarshal(nullDateTimeJSON, &null)
	maybePanic(err)
	assertNullDateTime(t, null, "null time json")
	if !null.Set {
		t.Error("should be Set")
	}

	var invalid DateTime
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.ParseError, not %T", err)
	}
	assertNullDateTime(t, invalid, "invalid from object json")

	var bad DateTime
	err = json.Unmarshal(badDateTimeObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullDateTime(t, bad, "bad from object json")

	var wrongType DateTime
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullDateTime(t, wrongType, "wrong type object json")
}

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
