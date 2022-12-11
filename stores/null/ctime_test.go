package null

import (
	"encoding/json"
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

func TestUnmarshalCTimeJSON(t *testing.T) {
	var ct CTime
	err := json.Unmarshal(ctimeJSON, &ct)
	maybePanic(err)
	assertCTime(t, ct, "UnmarshalJSON() json")

	var null CTime
	err = json.Unmarshal(nullCTimeJSON, &null)
	maybePanic(err)
	assertNullCTime(t, null, "null time json")
	if !null.Set {
		t.Error("should be Set")
	}

	var invalid CTime
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*time.ParseError); !ok {
		t.Errorf("expected json.ParseError, not %T", err)
	}
	assertNullCTime(t, invalid, "invalid from object json")

	var bad CTime
	err = json.Unmarshal(badCTimeObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullCTime(t, bad, "bad from object json")

	var wrongType CTime
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullCTime(t, wrongType, "wrong type object json")
}

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
