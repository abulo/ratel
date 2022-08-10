/*
Package other is an example of a custom driver
*/
package other

import (
	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/ini"
)

const driverName = "other"

// Encoder ...
var (
	// Encoder is the encoder for this driver
	Encoder = ini.Encoder
	// Decoder is the decoder for this driver
	Decoder = ini.Decoder
	// Driver is the exported symbol
	Driver = config.NewDriver(driverName, Decoder, Encoder)
)
