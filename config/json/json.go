// Package json use the https://github.com/json-iterator/go for parse json
package json

import (
	"github.com/abulo/ratel/config"
	"github.com/gookit/goutil/jsonutil"
	jsoniter "github.com/json-iterator/go"
)

var parser = jsoniter.ConfigCompatibleWithStandardLibrary

// Decoder ...
var (
	// Decoder for json
	Decoder config.Decoder = func(data []byte, v interface{}) (err error) {
		if config.JSONAllowComments {
			str := jsonutil.StripComments(string(data))
			return parser.Unmarshal([]byte(str), v)
		}
		return parser.Unmarshal(data, v)
	}

	// Encoder for json
	Encoder config.Encoder = parser.Marshal
	// Driver for json
	Driver = config.NewDriver(config.JSON, Decoder, Encoder)
)
