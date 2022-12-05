package query

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

var NULL = []byte("NULL")

// embeds  for custom json un/marshalling
type Any = interface{}

// NullTime 新定义
type NullTime sql.NullTime

func NewNullTime(s Any) NullTime {
	if reflect.TypeOf(s) == nil {
		return NullTime{}
	}
	return NullTime{
		Time:  cast.ToTime(s),
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullTime) UnmarshalJSON(data []byte) error {
	location := util.TimeZone()
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" {
		return nil
	}
	tt, err := time.ParseInLocation(`"`+"2006-01-02T15:04:05"+`"`, string(data), location)
	if err != nil {
		return err
	}
	n.Time = tt
	if n.Time.IsZero() {
		n.Valid = false
	} else {
		n.Valid = true
	}
	return err
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Time.IsZero() || !n.Valid {
		return NULL, nil
	}
	return []byte(n.Time.Format(time.RFC3339)), nil
}

// Result for NullString
func (n NullTime) Result() interface{} {
	if n.Valid {
		return n.Time.Format(time.RFC3339)
	}
	return cast.ToString(NULL)
}

// NullString 新定义
type NullString sql.NullString

// NewNullString 函数将一个字符串转换为sql.NullString
func NewNullString(s Any) NullString {
	res := cast.ToString(s)
	if len(res) == 0 {
		return NullString{}
	}
	return NullString{
		String: res,
		Valid:  true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullString) UnmarshalJSON(data []byte) error {

	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.String); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.String)
}

// Result for NullString
func (n NullString) Result() interface{} {
	if n.Valid {
		return n.String
	}
	return cast.ToString(NULL)
}

type NullBool sql.NullBool

func NewNullBool(s Any) NullBool {
	if reflect.TypeOf(s) == nil {
		return NullBool{}
	}
	return NullBool{
		Bool:  cast.ToBool(s),
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullBool) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Bool); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Bool)
}

func (n NullBool) Result() interface{} {
	if n.Valid {
		return n.Bool
	}
	return cast.ToString(NULL)
}

type NullByte sql.NullByte

func NewNullByte(s Any) NullByte {
	if reflect.TypeOf(s) == nil {
		return NullByte{}
	}
	return NullByte{
		Byte:  cast.ToUint8(s),
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullByte) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Byte); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullByte) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Byte)
}

// Result for NullString
func (n NullByte) Result() interface{} {
	if n.Valid {
		return cast.ToUint8(n.Byte)
	}
	return cast.ToString(NULL)
}

type NullFloat32 struct {
	Float32 float32
	Valid   bool // Valid is true if Float32 is not NULL
}

func NewNullFloat32(s Any) NullFloat32 {
	if reflect.TypeOf(s) == nil {
		return NullFloat32{}
	}
	return NullFloat32{
		Float32: cast.ToFloat32(s),
		Valid:   true,
	}
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value Any) error {
	if value == nil {
		n.Float32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	n.Float32 = cast.ToFloat32(value)
	return nil
}

func (n NullFloat32) Value() (Any, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float32, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullFloat32) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Float32); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullFloat32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Float32)
}

// Result for NullString
func (n NullFloat32) Result() interface{} {
	if n.Valid {
		return n.Float32
	}
	return cast.ToString(NULL)
}

type NullFloat64 sql.NullFloat64

func NewNullFloat64(s Any) NullFloat64 {
	if reflect.TypeOf(s) == nil {
		return NullFloat64{}
	}
	return NullFloat64{
		Float64: cast.ToFloat64(s),
		Valid:   true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Float64); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Float64)
}

// Result for NullString
func (n NullFloat64) Result() interface{} {
	if n.Valid {
		return n.Float64
	}
	return cast.ToString(NULL)
}

type NullInt16 sql.NullInt16

// NewNullInt16 函数将一个字符串转换为sql.NullInt16
func NewNullInt16(s Any) NullInt16 {
	if reflect.TypeOf(s) == nil {
		return NullInt16{}
	}
	res := cast.ToInt16(s)
	return NullInt16{
		Int16: res,
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullInt16) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Int16); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullInt16) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Int16)
}

// Result for NullInt16
func (n NullInt16) Result() interface{} {
	if n.Valid {
		return n.Int16
	}
	return cast.ToString(NULL)
}

type NullInt32 sql.NullInt32

// NewNullInt32 函数将一个字符串转换为sql.NullInt32
func NewNullInt32(s Any) NullInt32 {
	if reflect.TypeOf(s) == nil {
		return NullInt32{}
	}
	res := cast.ToInt32(s)
	return NullInt32{
		Int32: res,
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullInt32) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Int32); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Int32)
}

// Result for NullString
func (n NullInt32) Result() interface{} {
	if n.Valid {
		return n.Int32
	}
	return cast.ToString(NULL)
}

type NullInt64 sql.NullInt64

// NewNullInt64 函数将一个字符串转换为sql.NullInt64
func NewNullInt64(s Any) NullInt64 {
	if reflect.TypeOf(s) == nil {
		return NullInt64{}
	}
	res := cast.ToInt64(s)
	return NullInt64{
		Int64: res,
		Valid: true,
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullInt64) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Int64); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Int64)
}

// Result for NullInt64
func (n NullInt64) Result() interface{} {
	if n.Valid {
		return n.Int64
	}
	return cast.ToString(NULL)
}

// NullBytes can be an []byte or a null value.
type NullBytes struct {
	Bytes []byte
	Valid bool // Valid is true if Bytes is not NULL
}

// NewNullBytes 函数将一个字符串转换为sql.NullBytes
func NewNullBytes(s Any) NullBytes {
	if reflect.TypeOf(s) == nil {
		return NullBytes{}
	}
	return NullBytes{
		Bytes: []byte(cast.ToString(s)),
		Valid: true,
	}
}

// Scan implements the Scanner interface.
func (n *NullBytes) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	n.Bytes = []byte(cast.ToString(value))
	return nil
}

// Value implements the driver Valuer interface.
func (n NullBytes) Value() (Any, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (n *NullBytes) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "NULL" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &n.Bytes); err != nil {
		return nil
	}
	n.Valid = true
	return nil
}

func (n NullBytes) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return NULL, nil
	}
	return json.Marshal(n.Bytes)
}

// Result for NullBytes
func (n NullBytes) Result() interface{} {
	if n.Valid {
		return n.Bytes
	}
	return cast.ToString(NULL)
}
