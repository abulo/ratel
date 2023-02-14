package yamlv3

// see https://pkg.go.dev/gopkg.in/yaml.v3
import (
	"github.com/abulo/ratel/config"
	"gopkg.in/yaml.v3"
)

// Decoder the yaml content decoder
var Decoder config.Decoder = yaml.Unmarshal

// Encoder the yaml content encoder
var Encoder config.Encoder = yaml.Marshal

// Driver for yaml
var Driver = config.NewDriver(config.Yaml, Decoder, Encoder)
