package sql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/stores/null"
	"github.com/pkg/errors"
)

// Format
func Format(query string, args ...any) (string, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return query, nil
	}

	var b strings.Builder
	var argIndex int
	bytes := len(query)

	for i := 0; i < bytes; i++ {
		ch := query[i]
		switch ch {
		case '?':
			if argIndex >= numArgs {
				return "", fmt.Errorf("%d ? in sql, but less arguments provided", argIndex)
			}

			writeValue(&b, args[argIndex])
			argIndex++
		case ':', '$':
			var j int
			for j = i + 1; j < bytes; j++ {
				char := query[j]
				if char < '0' || '9' < char {
					break
				}
			}

			if j > i+1 {
				index, err := strconv.Atoi(query[i+1 : j])
				if err != nil {
					return "", err
				}

				// index starts from 1 for pg or oracle
				if index > argIndex {
					argIndex = index
				}

				index--
				if index < 0 || numArgs <= index {
					return "", fmt.Errorf("wrong index %d in sql", index)
				}

				writeValue(&b, args[index])
				i = j - 1
			}
		case '\'', '"', '`':
			b.WriteByte(ch)

			for j := i + 1; j < bytes; j++ {
				cur := query[j]
				b.WriteByte(cur)

				if cur == '\\' {
					j++
					if j >= bytes {
						return "", errors.New("no char after escape char")
					}

					b.WriteByte(query[j])
				} else if cur == ch {
					i = j
					break
				}
			}
		default:
			b.WriteByte(ch)
		}
	}

	if argIndex < numArgs {
		return "", fmt.Errorf("%d arguments provided, not matching sql", argIndex)
	}

	return b.String(), nil
}

func writeValue(buf *strings.Builder, arg any) {
	switch v := arg.(type) {
	case bool:
		if v {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	case string:
		buf.WriteByte('\'')
		buf.WriteString(escape(v))
		buf.WriteByte('\'')
	case time.Time:
		buf.WriteByte('\'')
		buf.WriteString(v.String())
		buf.WriteByte('\'')
	case *time.Time:
		buf.WriteByte('\'')
		buf.WriteString(v.String())
		buf.WriteByte('\'')
	case null.Bool:
		buf.WriteString(v.Result())
	case null.Byte:
		buf.WriteString(v.Result())
	case null.Bytes:
		buf.WriteString(v.Result())
	case null.CTime:
		buf.WriteString(v.Result())
	case null.Date:
		buf.WriteString(v.Result())
	case null.DateTime:
		buf.WriteString(v.Result())
	case null.Float32:
		buf.WriteString(v.Result())
	case null.Float64:
		buf.WriteString(v.Result())
	case null.Int:
		buf.WriteString(v.Result())
	case null.Int8:
		buf.WriteString(v.Result())
	case null.Int16:
		buf.WriteString(v.Result())
	case null.Int32:
		buf.WriteString(v.Result())
	case null.Int64:
		buf.WriteString(v.Result())
	case null.JSON:
		buf.WriteString(v.Result())
	case null.String:
		buf.WriteString(v.Result())
	case null.Time:
		buf.WriteString(v.Result())
	case null.TimeStamp:
		buf.WriteString(v.Result())
	case null.Uint:
		buf.WriteString(v.Result())
	case null.Uint8:
		buf.WriteString(v.Result())
	case null.Uint16:
		buf.WriteString(v.Result())
	case null.Uint32:
		buf.WriteString(v.Result())
	case null.Uint64:
		buf.WriteString(v.Result())
	default:
		buf.WriteString(replace(v))
	}
}

func escape(input string) string {
	var b strings.Builder

	for _, ch := range input {
		switch ch {
		case '\x00':
			b.WriteString(`\x00`)
		case '\r':
			b.WriteString(`\r`)
		case '\n':
			b.WriteString(`\n`)
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '"':
			b.WriteString(`\"`)
		case '\x1a':
			b.WriteString(`\x1a`)
		default:
			b.WriteRune(ch)
		}
	}

	return b.String()
}

func replace(v any) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return replaceOfValue(val)
}

func replaceOfValue(val reflect.Value) string {
	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}
