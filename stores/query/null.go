package query

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"reflect"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

var nullJSON = []byte("null")

// Nullable ...
type Nullable interface {
	IsEmpty() bool
	IsValid() bool
}

// NullString 空字符串
type NullString struct {
	sql.NullString
}

// NewString ...
func NewString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: true}}
}

// NewNullString ...
func NewNullString() NullString {
	return NullString{sql.NullString{String: "", Valid: false}}
}

// MarshalJSON ...
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullString) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullString) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullString) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullString) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

// Convert ...
func (n NullString) Convert() string {
	if n.Valid {
		return n.String
	}
	return ""
}

// NullDateTime NullDate 空字符串
type NullDateTime struct {
	sql.NullString
}

// NewDateTime ...
func NewDateTime(s interface{}) NullDateTime {
	return NullDateTime{sql.NullString{String: util.Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
}

// NewNullDateTime ...
func NewNullDateTime() NullDateTime {
	return NullDateTime{sql.NullString{String: "", Valid: false}}
}

// Scan ...
func (nt *NullDateTime) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullDateTime()
		return nil
	} else {
		var s sql.NullString
		if err := s.Scan(value); err != nil {
			return err
		} else {
			*nt, err = NewDateTime(s.String), nil
			return err
		}
	}
}

// MarshalJSON ...
func (ns NullDateTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullDateTime) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullDateTime) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullDateTime) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullDateTime) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

// NullInt64 ...
type NullInt64 struct{ sql.NullInt64 }

// NewInt64 ...
func NewInt64(i int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: i, Valid: true}}
}

// NewNullInt64 ...
func NewNullInt64() NullInt64 {
	return NullInt64{sql.NullInt64{Int64: 0, Valid: false}}
}

// MarshalJSON ...
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON ...
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ni.Int64 = 0
		ni.Valid = false
	} else {
		err = json.Unmarshal(b, &ni.Int64)
		ni.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (n NullInt64) IsEmpty() bool {
	return !n.Valid
}

// IsValid ...
func (n NullInt64) IsValid() bool {
	return n.Valid
}

// Result ...
func (ns NullInt64) Result() interface{} {
	if ns.Valid {
		return ns.Int64
	}
	return "NULL"
}

// Convert ...
func (n NullInt64) Convert() string {
	if n.Valid {
		return cast.ToString(n.Int64)
	}
	return ""
}

// NullInt32 ...
type NullInt32 struct{ sql.NullInt32 }

// NewInt32 ...
func NewInt32(i int32) NullInt32 {
	return NullInt32{sql.NullInt32{Int32: i, Valid: true}}
}

// NewNullInt32 ...
func NewNullInt32() NullInt32 {
	return NullInt32{sql.NullInt32{Int32: 0, Valid: false}}
}

// MarshalJSON ...
func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ni.Int32)
}

// UnmarshalJSON ...
func (ni *NullInt32) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ni.Int32 = 0
		ni.Valid = false
	} else {
		err = json.Unmarshal(b, &ni.Int32)
		ni.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (n NullInt32) IsEmpty() bool {
	return !n.Valid
}

// IsValid ...
func (n NullInt32) IsValid() bool {
	return n.Valid
}

// Result ...
func (ns NullInt32) Result() interface{} {
	if ns.Valid {
		return ns.Int32
	}
	return "NULL"
}

// Convert ...
func (n NullInt32) Convert() string {
	if n.Valid {
		return cast.ToString(n.Int32)
	}
	return ""
}

// NullFloat64 ...
type NullFloat64 struct{ sql.NullFloat64 }

// NewFloat64 ...
func NewFloat64(f float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: f, Valid: true}}
}

// NewNullFloat64 ...
func NewNullFloat64() NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: 0.0, Valid: false}}
}

// MarshalJSON ...
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON ...
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nf.Float64 = 0.0
		nf.Valid = false
	} else {
		err = json.Unmarshal(b, &nf.Float64)
		nf.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (n NullFloat64) IsEmpty() bool {
	return !n.Valid
}

// IsValid ...
func (n NullFloat64) IsValid() bool {
	return n.Valid
}

// Result ...
func (ns NullFloat64) Result() interface{} {
	if ns.Valid {
		return ns.Float64
	}
	return "NULL"
}

// Convert ...
func (n NullFloat64) Convert() string {
	if n.Valid {
		return cast.ToString(n.Float64)
	}
	return ""
}

// NullBool ...
type NullBool struct{ sql.NullBool }

// NewBool ...
func NewBool(b bool) NullBool {
	return NullBool{sql.NullBool{Bool: b, Valid: true}}
}

// NewNullBool ...
func NewNullBool() NullBool {
	return NullBool{sql.NullBool{Bool: false, Valid: false}}
}

// MarshalJSON ...
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON ...
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nb.Bool = false
		nb.Valid = false
	} else {
		err = json.Unmarshal(b, &nb.Bool)
		nb.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (n NullBool) IsEmpty() bool {
	return !n.Valid
}

// IsValid ...
func (n NullBool) IsValid() bool {
	return n.Valid
}

// Result ...
func (ns NullBool) Result() interface{} {
	if ns.Valid {
		if ns.Bool {
			return 1
		} else {
			return 0
		}
	}
	return "NULL"
}

// Convert ...
func (n NullBool) Convert() string {
	if n.Valid {
		if n.Bool {
			return cast.ToString("true")
		} else {
			return cast.ToString("false")
		}
	}
	return "false"
}

// NullDate 空字符串
type NullDate struct {
	sql.NullString
}

// NewDate ...
func NewDate(s interface{}) NullDate {
	return NullDate{sql.NullString{String: util.Date("Y-m-d", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
}

// NewNullDate ...
func NewNullDate() NullDate {
	return NullDate{sql.NullString{String: "", Valid: false}}
}

// Scan ...
func (nt *NullDate) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullDate()
		return nil
	} else {
		var s sql.NullString
		if err := s.Scan(value); err != nil {
			return err
		} else {
			*nt, err = NewDate(s.String), nil
			return err
		}
	}
}

// MarshalJSON ...
func (ns NullDate) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullDate) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullDate) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullDate) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullDate) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

// NullTime 空字符串
type NullTime struct {
	sql.NullString
}

// NewTime ...
func NewTime(s interface{}) NullTime {
	return NullTime{sql.NullString{String: util.Date("H:i:s", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
}

// NewNullTime ...
func NewNullTime() NullTime {
	return NullTime{sql.NullString{String: "", Valid: false}}
}

// Scan ...
func (nt *NullTime) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullTime()
		return nil
	} else {
		var s sql.NullString
		if err := s.Scan(value); err != nil {
			return err
		} else {
			*nt, err = NewTime(s.String), nil
			return err
		}
	}
}

// MarshalJSON ...
func (ns NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullTime) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullTime) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullTime) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullTime) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

// NullYear 空字符串
type NullYear struct {
	sql.NullString
}

// NewYear ...
func NewYear(s interface{}) NullYear {
	return NullYear{sql.NullString{String: util.Date("Y", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
}

// NewNullYear ...
func NewNullYear() NullYear {
	return NullYear{sql.NullString{String: "", Valid: false}}
}

// Scan ...
func (nt *NullYear) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullYear()
		return nil
	} else {
		var s sql.NullString
		if err := s.Scan(value); err != nil {
			return err
		} else {
			*nt, err = NewYear(s.String), nil
			return err
		}
	}
}

// MarshalJSON ...
func (ns NullYear) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullYear) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullYear) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullYear) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullYear) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

// NullTimeStamp NullYear 空字符串
type NullTimeStamp struct {
	sql.NullString
}

// NewTimeStamp ...
func NewTimeStamp(s interface{}) NullTimeStamp {
	return NullTimeStamp{sql.NullString{String: cast.ToString(cast.ToTimeInDefaultLocation(s, util.TimeZone()).Unix()), Valid: true}}
}

// NewNullTimeStamp ...
func NewNullTimeStamp() NullTimeStamp {
	return NullTimeStamp{sql.NullString{String: "", Valid: false}}
}

// Scan ...
func (nt *NullTimeStamp) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullTimeStamp()
		return nil
	} else {
		var s sql.NullString
		if err := s.Scan(value); err != nil {
			return err
		} else {
			*nt, err = NewTimeStamp(s.String), nil
			return err
		}
	}
}

// MarshalJSON ...
func (ns NullTimeStamp) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON ...
func (ns *NullTimeStamp) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		ns.String = ""
		ns.Valid = false
	} else {
		err = json.Unmarshal(b, &ns.String)
		ns.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (ns NullTimeStamp) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

// IsValid ...
func (ns NullTimeStamp) IsValid() bool {
	return ns.Valid
}

// Result ...
func (ns NullTimeStamp) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}
