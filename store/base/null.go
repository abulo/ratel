package base

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"reflect"

	"github.com/abulo/ratel/util"
)

var nullJSON = []byte("null")

type Nullable interface {
	IsEmpty() bool
	IsValid() bool
}

//NullString 空字符串
type NullString struct {
	sql.NullString
}

func NewString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: true}}
}

func NewNullString() NullString {
	return NullString{sql.NullString{String: "", Valid: false}}
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

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

func (ns NullString) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

func (ns NullString) IsValid() bool {
	return ns.Valid
}

func (ns NullString) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

func (n NullString) Convert() string {
	if n.Valid {
		return n.String
	}
	return ""
}

//NullDate 空字符串
type NullDateTime struct {
	sql.NullString
}

func NewDateTime(s interface{}) NullDateTime {
	return NullDateTime{sql.NullString{String: util.Date("Y-m-d H:i:s", util.ToTime(s)), Valid: true}}
}

func NewNullDateTime() NullDateTime {
	return NullDateTime{sql.NullString{String: "", Valid: false}}
}

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

func (ns NullDateTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

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

func (ns NullDateTime) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

func (ns NullDateTime) IsValid() bool {
	return ns.Valid
}

func (ns NullDateTime) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

type NullInt64 struct{ sql.NullInt64 }

func NewInt64(i int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: i, Valid: true}}
}

func NewNullInt64() NullInt64 {
	return NullInt64{sql.NullInt64{Int64: 0, Valid: false}}
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ni.Int64)
}

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

func (n NullInt64) IsEmpty() bool {
	return !n.Valid
}

func (n NullInt64) IsValid() bool {
	return n.Valid
}

func (ns NullInt64) Result() interface{} {
	if ns.Valid {
		return ns.Int64
	}
	return "NULL"
}

func (n NullInt64) Convert() string {
	if n.Valid {
		return util.ToString(n.Int64)
	}
	return ""
}

type NullInt32 struct{ sql.NullInt32 }

func NewInt32(i int32) NullInt32 {
	return NullInt32{sql.NullInt32{Int32: i, Valid: true}}
}

func NewNullInt32() NullInt32 {
	return NullInt32{sql.NullInt32{Int32: 0, Valid: false}}
}

func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ni.Int32)
}

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

func (n NullInt32) IsEmpty() bool {
	return !n.Valid
}

func (n NullInt32) IsValid() bool {
	return n.Valid
}

func (ns NullInt32) Result() interface{} {
	if ns.Valid {
		return ns.Int32
	}
	return "NULL"
}

func (n NullInt32) Convert() string {
	if n.Valid {
		return util.ToString(n.Int32)
	}
	return ""
}

type NullFloat64 struct{ sql.NullFloat64 }

func NewFloat64(f float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: f, Valid: true}}
}

func NewNullFloat64() NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: 0.0, Valid: false}}
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nf.Float64)
}

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

func (n NullFloat64) IsEmpty() bool {
	return !n.Valid
}

func (n NullFloat64) IsValid() bool {
	return n.Valid
}

func (ns NullFloat64) Result() interface{} {
	if ns.Valid {
		return ns.Float64
	}
	return "NULL"
}

func (n NullFloat64) Convert() string {
	if n.Valid {
		return util.ToString(n.Float64)
	}
	return ""
}

type NullBool struct{ sql.NullBool }

func NewBool(b bool) NullBool {
	return NullBool{sql.NullBool{Bool: b, Valid: true}}
}

func NewNullBool() NullBool {
	return NullBool{sql.NullBool{Bool: false, Valid: false}}
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nb.Bool)
}

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

func (n NullBool) IsEmpty() bool {
	return !n.Valid
}

func (n NullBool) IsValid() bool {
	return n.Valid
}

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

func (n NullBool) Convert() string {
	if n.Valid {
		if n.Bool {
			return util.ToString("true")
		} else {
			return util.ToString("false")
		}
	}
	return "false"
}

//NullDate 空字符串
type NullDate struct {
	sql.NullString
}

func NewDate(s interface{}) NullDate {
	return NullDate{sql.NullString{String: util.Date("Y-m-d", util.ToTime(s)), Valid: true}}
}

func NewNullDate() NullDate {
	return NullDate{sql.NullString{String: "", Valid: false}}
}

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

func (ns NullDate) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return nullJSON, nil
	}
	return json.Marshal(ns.String)
}

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

func (ns NullDate) IsEmpty() bool {
	return !ns.Valid || ns.String == ""
}

func (ns NullDate) IsValid() bool {
	return ns.Valid
}

func (ns NullDate) Result() interface{} {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}
