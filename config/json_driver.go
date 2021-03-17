package config

// default json driver(encoder/decoder)
import (
	"encoding/json"
)

// JSONAllowComments support write comments on json file.
var JSONAllowComments = true

// JSONDecoder for json decode
var JSONDecoder Decoder = func(data []byte, v interface{}) (err error) {
	if JSONAllowComments {
		str := StripComments(string(data))
		return json.Unmarshal([]byte(str), v)
	}

	return json.Unmarshal(data, v)
}

// JSONEncoder for json encode
var JSONEncoder Encoder = json.Marshal

// JSONDriver instance fot json
var JSONDriver = &jsonDriver{name: JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name          string
	ClearComments bool
}

// Name of the driver
func (d *jsonDriver) Name() string {
	return d.name
}

// GetDecoder for json
func (d *jsonDriver) GetDecoder() Decoder {
	return JSONDecoder
}

// GetEncoder for json
func (d *jsonDriver) GetEncoder() Encoder {
	return JSONEncoder
}
