package query

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"reflect"

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
	return NullString{sql.NullString{String: "", Valid: true}}
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
func (ns NullString) Convert() string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullDateTime NullDate 空字符串
type NullDateTime struct {
	sql.NullString
}

// NewDateTime ...
func NewDateTime(s interface{}) NullDateTime {
	// return NullDateTime{sql.NullString{String: util.Date("Y-m-d H:i:s", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
	return NullDateTime{sql.NullString{String: cast.ToString(s), Valid: true}}
}

// NewNullDateTime ...
func NewNullDateTime() NullDateTime {
	return NullDateTime{sql.NullString{String: "", Valid: true}}
}

// Scan ...
func (nt *NullDateTime) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullDateTime()
		return nil
	}
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*nt = NewDateTime(s.String)
	return nil
}

// MarshalJSON ...
func (nt NullDateTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nt.String)
}

// UnmarshalJSON ...
func (nt *NullDateTime) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nt.String = ""
		nt.Valid = false
	} else {
		err = json.Unmarshal(b, &nt.String)
		nt.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (nt NullDateTime) IsEmpty() bool {
	return !nt.Valid || nt.String == ""
}

// IsValid ...
func (nt NullDateTime) IsValid() bool {
	return nt.Valid
}

// Result ...
func (nt NullDateTime) Result() interface{} {
	if nt.Valid {
		return nt.String
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
	return NullInt64{sql.NullInt64{Int64: 0, Valid: true}}
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
func (ni NullInt64) IsEmpty() bool {
	return !ni.Valid
}

// IsValid ...
func (ni NullInt64) IsValid() bool {
	return ni.Valid
}

// Result ...
func (ni NullInt64) Result() interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return "NULL"
}

// Convert ...
func (ni NullInt64) Convert() string {
	if ni.Valid {
		return cast.ToString(ni.Int64)
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
	return NullInt32{sql.NullInt32{Int32: 0, Valid: true}}
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
func (ni NullInt32) IsEmpty() bool {
	return !ni.Valid
}

// IsValid ...
func (ni NullInt32) IsValid() bool {
	return ni.Valid
}

// Result ...
func (ni NullInt32) Result() interface{} {
	if ni.Valid {
		return ni.Int32
	}
	return "NULL"
}

// Convert ...
func (ni NullInt32) Convert() string {
	if ni.Valid {
		return cast.ToString(ni.Int32)
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
	return NullFloat64{sql.NullFloat64{Float64: 0.0, Valid: true}}
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
func (nf NullFloat64) IsEmpty() bool {
	return !nf.Valid
}

// IsValid ...
func (nf NullFloat64) IsValid() bool {
	return nf.Valid
}

// Result ...
func (nf NullFloat64) Result() interface{} {
	if nf.Valid {
		return nf.Float64
	}
	return "NULL"
}

// Convert ...
func (nf NullFloat64) Convert() string {
	if nf.Valid {
		return cast.ToString(nf.Float64)
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
	return NullBool{sql.NullBool{Bool: false, Valid: true}}
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
func (nb NullBool) IsEmpty() bool {
	return !nb.Valid
}

// IsValid ...
func (nb NullBool) IsValid() bool {
	return nb.Valid
}

// Result ...
func (nb NullBool) Result() interface{} {
	if nb.Valid {
		if nb.Bool {
			return 1
		}
		return 0
	}
	return "NULL"
}

// Convert ...
func (nb NullBool) Convert() string {
	if nb.Valid {
		if nb.Bool {
			return cast.ToString("true")
		}
		return cast.ToString("false")

	}
	return "false"
}

// NullDate 空字符串
type NullDate struct {
	sql.NullString
}

// NewDate ...
func NewDate(s interface{}) NullDate {
	// return NullDate{sql.NullString{String: util.Date("Y-m-d", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
	return NullDate{sql.NullString{String: cast.ToString(s), Valid: true}}
}

// NewNullDate ...
func NewNullDate() NullDate {
	return NullDate{sql.NullString{String: "", Valid: true}}
}

// Scan ...
func (nt *NullDate) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullDate()
		return nil
	}
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*nt = NewDate(s.String)
	return nil

}

// MarshalJSON ...
func (nt NullDate) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nt.String)
}

// UnmarshalJSON ...
func (nt *NullDate) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nt.String = ""
		nt.Valid = false
	} else {
		err = json.Unmarshal(b, &nt.String)
		nt.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (nt NullDate) IsEmpty() bool {
	return !nt.Valid || nt.String == ""
}

// IsValid ...
func (nt NullDate) IsValid() bool {
	return nt.Valid
}

// Result ...
func (nt NullDate) Result() interface{} {
	if nt.Valid {
		return nt.String
	}
	return "NULL"
}

// NullTime 空字符串
type NullTime struct {
	sql.NullString
}

// NewTime ...
func NewTime(s interface{}) NullTime {
	// return NullTime{sql.NullString{String: util.Date("H:i:s", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
	return NullTime{sql.NullString{String: cast.ToString(s), Valid: true}}
}

// NewNullTime ...
func NewNullTime() NullTime {
	return NullTime{sql.NullString{String: "", Valid: true}}
}

// Scan ...
func (nt *NullTime) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullTime()
		return nil
	}
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*nt = NewTime(s.String)
	return nil
}

// MarshalJSON ...
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nt.String)
}

// UnmarshalJSON ...
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nt.String = ""
		nt.Valid = false
	} else {
		err = json.Unmarshal(b, &nt.String)
		nt.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (nt NullTime) IsEmpty() bool {
	return !nt.Valid || nt.String == ""
}

// IsValid ...
func (nt NullTime) IsValid() bool {
	return nt.Valid
}

// Result ...
func (nt NullTime) Result() interface{} {
	if nt.Valid {
		return nt.String
	}
	return "NULL"
}

// NullYear 空字符串
type NullYear struct {
	sql.NullString
}

// NewYear ...
func NewYear(s interface{}) NullYear {
	// return NullYear{sql.NullString{String: util.Date("Y", cast.ToTimeInDefaultLocation(s, util.TimeZone())), Valid: true}}
	return NullYear{sql.NullString{String: cast.ToString(s), Valid: true}}
}

// NewNullYear ...
func NewNullYear() NullYear {
	return NullYear{sql.NullString{String: "", Valid: true}}
}

// Scan ...
func (nt *NullYear) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullYear()
		return nil
	}
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*nt = NewYear(s.String)
	return nil
}

// MarshalJSON ...
func (nt NullYear) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nt.String)
}

// UnmarshalJSON ...
func (nt *NullYear) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nt.String = ""
		nt.Valid = false
	} else {
		err = json.Unmarshal(b, &nt.String)
		nt.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (nt NullYear) IsEmpty() bool {
	return !nt.Valid || nt.String == ""
}

// IsValid ...
func (nt NullYear) IsValid() bool {
	return nt.Valid
}

// Result ...
func (nt NullYear) Result() interface{} {
	if nt.Valid {
		return nt.String
	}
	return "NULL"
}

// NullTimeStamp NullYear 空字符串
type NullTimeStamp struct {
	sql.NullString
}

// NewTimeStamp ...
func NewTimeStamp(s interface{}) NullTimeStamp {
	// return NullTimeStamp{sql.NullString{String: cast.ToString(cast.ToTimeInDefaultLocation(s, util.TimeZone()).Unix()), Valid: true}}
	return NullTimeStamp{sql.NullString{String: cast.ToString(s), Valid: true}}
}

// NewNullTimeStamp ...
func NewNullTimeStamp() NullTimeStamp {
	return NullTimeStamp{sql.NullString{String: "", Valid: true}}
}

// Scan ...
func (nt *NullTimeStamp) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NewNullTimeStamp()
		return nil
	}
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*nt = NewTimeStamp(s.String)
	return nil
}

// MarshalJSON ...
func (nt NullTimeStamp) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return nullJSON, nil
	}
	return json.Marshal(nt.String)
}

// UnmarshalJSON ...
func (nt *NullTimeStamp) UnmarshalJSON(b []byte) error {
	var err error = nil
	if bytes.Equal(nullJSON, b) {
		nt.String = ""
		nt.Valid = false
	} else {
		err = json.Unmarshal(b, &nt.String)
		nt.Valid = (err == nil)
	}
	return err
}

// IsEmpty ...
func (nt NullTimeStamp) IsEmpty() bool {
	return !nt.Valid || nt.String == ""
}

// IsValid ...
func (nt NullTimeStamp) IsValid() bool {
	return nt.Valid
}

// Result ...
func (nt NullTimeStamp) Result() interface{} {
	if nt.Valid {
		return nt.String
	}
	return "NULL"
}
