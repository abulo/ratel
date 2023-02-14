/*
Package hclv2 is driver use HCL format content as config source

about HCL, please see https://github.com/hashicorp/hcl
docs for HCL v2 https://pkg.go.dev/github.com/hashicorp/hcl/v2
*/
package hclv2

import (
	"github.com/pkg/errors"

	"github.com/abulo/ratel/v2/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Decoder the hcl content decoder
var Decoder config.Decoder = func(blob []byte, v interface{}) (err error) {
	// return hclsimple.Decode("hcl2/config.hcl", blob, nil, v)
	file, diags := hclsyntax.ParseConfig(
		blob,
		"hcl2/config.hcl",
		hcl.Pos{Line: 0, Column: 0},
	)
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
var Driver = config.NewDriver(config.Hcl, Decoder, Encoder)
