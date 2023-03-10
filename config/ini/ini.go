/*
Package ini is driver use INI format content as config source

about ini parse, please see https://github.com/gookit/ini/parser
*/
package ini

import (
	"github.com/abulo/ratel/v3/config"
	"github.com/gookit/ini/v2/parser"
)

// Decoder the ini content decoder
var Decoder config.Decoder = parser.Decode

// Encoder encode data to ini content
var Encoder config.Encoder = func(ptr any) (out []byte, err error) {
	return parser.Encode(ptr)
}

// Driver for ini
var Driver = config.NewDriver(config.Ini, Decoder, Encoder)
