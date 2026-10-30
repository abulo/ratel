package fontx

import (
	"image/color"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// FontConfig 保存绘制文字所需的字体、字号与颜色配置。
type FontConfig struct {
	// Font 为已加载的 TrueType 字体。
	Font *truetype.Font
	// FontSize 为字号大小。
	FontSize float64
	// Color 为文字颜色。
	Color color.Color
}

// TrueTypeFont 包装已解析的 TrueType 字体，提供字体度量信息获取能力。
type TrueTypeFont struct {
	// Font 为已解析的 TrueType 字体实例。
	Font *truetype.Font
}

// LoadFont 读取并解析指定路径的字体文件；文件读取失败或字体解析失败时返回错误。
func LoadFont(fontFile string) (*TrueTypeFont, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	return &TrueTypeFont{f}, err
}

// NewFreeTypeContext 创建一个新的 freetype 绘制上下文。
func NewFreeTypeContext() *freetype.Context {
	return freetype.NewContext()
}

// GetMetrics 按 FontConfig 中的字号计算字体度量信息（字高、上行、下行等），
// 内部固定使用 DPI 72 与无提示模式。
func (fo *TrueTypeFont) GetMetrics(fc *FontConfig) font.Metrics {
	opt := &truetype.Options{
		Size:    fc.FontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	face := truetype.NewFace(fo.Font, opt)
	return face.Metrics()
}
