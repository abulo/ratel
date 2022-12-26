## null [![GoDoc](https://godoc.org/github.com/mbahjadol/null?status.svg)](https://godoc.org/github.com/mbahjadol/null) [![Coverage](http://gocover.io/_badge/github.com/mbahjadol/null)](http://gocover.io/github.com/mbahjadol/null)

`null` is a library with reasonable options for dealing with nullable SQL and
JSON values.

Types in `null` will only be considered null on null input, and will JSON
encode to `null`.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this
library in place of `sql.NullXXX`. All types also implement:
`encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`,
`json.Unmarshaler` and `sql.Scanner`.

---

### Installation

Null installation, ensure you're using go modules.

```
go get github.com/mbahjadol/null
```

### Usage

The following are all types supported in this package. All types will marshal
to JSON null if Invalid or SQL source data is null.

| Type | Description | Notes |
|------|-------------|-------|
| `null.Date` | Nullable `time.Time` | That accept only full-year, month, date and ignore elses. Marshals to JSON null if SQL source data is null. Uses customize `time.Time`'s json & text marshaler that support format with YYYY-MM-DD. |
| `null.CTime` | Nullable `time.Time` | That accept only hour, minute, second and ignore elses. Marshals to JSON null if SQL source data is null. Uses customize `time.Time`'s json & text marshaler that support format with hh:mm:ss . It is use name CTime (Customize Time) to differentiate null.Time and partially support for RFC3339 only time |
| `null.DateTime` | Nullable `time.Time` | That accept full-year,month,date,hour,minute,second and ignore elses. Marshals to JSON null if SQL source data is null. Uses customize `time.Time`'s json & text marshaler that support format with YYYY-MM-DD hh:mm:ss |
| `null.TimeStamp` | Nullable `time.Time` | That accept full-year,month,date,hour,minute,second and microsecond. Marshals to JSON null if SQL source data is null. Uses customize `time.Time`'s json & text marshaler that support format with YYYY-MM-DD hh:mm:ss.zzzzzz |
| `null.Time` | Nullable `time.Time` | Marshals to JSON null if SQL source data is null. Uses `time.Time`'s marshaler. It marshall json with RFC3339.|
| `null.JSON` | Nullable `[]byte` | Will marshal to JSON null if Invalid. `[]byte{}` input will not produce an Invalid JSON, but `[]byte(nil)` will. This should be used for storing raw JSON in the database. Also has `null.JSON.Marshal` and `null.JSON.Unmarshal` helpers to marshal and unmarshal foreign objects. |
| `null.Bytes` | Nullable `[]byte` | `[]byte{}` input will not produce an Invalid Bytes, but `[]byte(nil)` will. This should be used for storing binary data (bytes in PSQL for example) in the database. |
| `null.String` | Nullable `string` | |
| `null.Byte` | Nullable `byte` | |
| `null.Bool` | Nullable `bool` | |
| `null.Float32` | Nullable `float32` | |
| `null.Float64` | Nullable `float64` | |
| `null.Int` | Nullable `int` | |
| `null.Int8` | Nullable `int8` | |
| `null.Int16` | Nullable `int16` | |
| `null.Int32` | Nullable `int32` | |
| `null.Int64` | Nullable `int64` | |
| `null.Uint` | Nullable `uint` | |
| `null.Uint8` | Nullable `uint8` | |
| `null.Uint16` | Nullable `uint16` | |
| `null.Uint32` | Nullable `uint32` | |
| `null.Uint64` | Nullable `uint64` | | |

### Bugs

`json`'s `",omitempty"` struct tag does not work correctly right now. It will
never omit a null or empty String. This might be [fixed
eventually](https://github.com/golang/go/issues/4357).


### License

BSD-3 (See License file)
