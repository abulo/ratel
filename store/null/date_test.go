package null

import (
	"encoding/json"
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

func TestUnmarshalDateJSON(t *testing.T) {
	var dt Date
	err := json.Unmarshal(dateJSON, &dt)
	maybePanic(err)
	assertDate(t, dt, "UnmarshalJSON() json")

	var null Date
	err = json.Unmarshal(nullDateJSON, &null)
	maybePanic(err)
	assertNullDate(t, null, "null time json")
	if !null.Set {
		t.Error("should be Set")
	}

	var invalid Date
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.ParseError, not %T", err)
	}
	assertNullDate(t, invalid, "invalid from object json")

	var bad Date
	err = json.Unmarshal(badDateObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullDate(t, bad, "bad from object json")

	var wrongType Date
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullDate(t, wrongType, "wrong type object json")
}

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

func TestMarshalDate(t *testing.T) {
	dt := DateFrom(dateValue)
	data, err := json.Marshal(dt)
	maybePanic(err)
	assertJSONEquals(t, data, string(dateJSON), "non-empty json marshal")

	dt.Valid = false
	data, err = json.Marshal(dt)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
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
