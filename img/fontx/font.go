package fontx

import (
	"image/color"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type FontConfig struct {
	Font     *truetype.Font
	FontSize float64
	Color    color.Color
}

type TrueTypeFont struct {
	Font *truetype.Font
}

// LoadFont loads a font file and returns *truetype.Font.
func LoadFont(fontFile string) (*TrueTypeFont, error) {
	// Read the font data.
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	return &TrueTypeFont{f}, err
}

func NewFreeTypeContext() *freetype.Context {
	return freetype.NewContext()
}

// GetMetrics gets font metrics by options.
func (fo *TrueTypeFont) GetMetrics(fc *FontConfig) font.Metrics {
	opt := &truetype.Options{
		Size:    fc.FontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	face := truetype.NewFace(fo.Font, opt)
	return face.Metrics()
}
