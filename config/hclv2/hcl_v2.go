package hclv2

import (
	"errors"

	"github.com/abulo/ratel/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Decoder the hcl content decoder
var Decoder config.Decoder = func(blob []byte, v interface{}) (err error) {
	// return hclsimple.Decode("hcl2/config.hcl", blob, nil, v)
	file, diags := hclsyntax.ParseConfig(blob, "hcl2/config.hcl", hcl.Pos{Line: 0, Column: 0})
	// if diags.HasErrors() {
	if len(diags) != 0 {
		return diags
	}

	return gohcl.DecodeBody(file.Body, nil, v)
}

// Encoder the hcl content encoder
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	err = errors.New("HCLv2: is not support encode data to HCL")
	return
}

// Driver instance for hcl
var Driver = &hclDriver{config.Hcl}

// hclDriver for hcl format content
type hclDriver struct {
	name string
}

// Name
func (d *hclDriver) Name() string {
	return d.name
}

// GetDecoder for hcl
func (d *hclDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for hcl
func (d *hclDriver) GetEncoder() config.Encoder {
	return Encoder
}
