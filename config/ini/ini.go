package ini

import (
	"github.com/abulo/ratel/config"
	"github.com/abulo/ratel/config/parser/ini/parser"
)

// Decoder the ini content decoder
var Decoder config.Decoder = parser.Decode

// Encoder encode data to ini content
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	return parser.Encode(ptr)
}

// Driver for ini
var Driver = &iniDriver{config.Ini}

// iniDriver for ini format content
type iniDriver struct {
	name string
}

// Name get name
func (d *iniDriver) Name() string {
	return d.name
}

// GetDecoder for ini
func (d *iniDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for ini
func (d *iniDriver) GetEncoder() config.Encoder {
	return Encoder
}
