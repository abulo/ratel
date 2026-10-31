package fontx

import (
	"image/color"
	"os"
	"sync"

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
	// faceMu 保护 faceCache 的并发访问（PERF-7-43）。
	faceMu sync.Mutex
	// faceCache 以字号为键缓存已构建的排版面；DPI 与 Hinting 在此固定为 72 / HintingNone。
	faceCache map[float64]font.Face
}

// LoadFont 读取并解析指定路径的字体文件；文件读取失败或字体解析失败时返回错误。
func LoadFont(fontFile string) (*TrueTypeFont, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	return &TrueTypeFont{Font: f}, err
}

// NewFreeTypeContext 创建一个新的 freetype 绘制上下文。
func NewFreeTypeContext() *freetype.Context {
	return freetype.NewContext()
}

// GetMetrics 按 FontConfig 中的字号计算字体度量信息（字高、上行、下行等），
// 内部固定使用 DPI 72 与无提示模式。
func (fo *TrueTypeFont) GetMetrics(fc *FontConfig) font.Metrics {
	// PERF-7-43 [P2] 按字号缓存 font.Face，同一字体同一字号反复取度量时复用已构建排版面，避免每次 NewFace
	fo.faceMu.Lock()
	defer fo.faceMu.Unlock()
	if fo.faceCache == nil {
		fo.faceCache = make(map[float64]font.Face)
	}
	face, ok := fo.faceCache[fc.FontSize]
	if !ok {
		opt := &truetype.Options{
			Size:    fc.FontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}
		face = truetype.NewFace(fo.Font, opt)
		fo.faceCache[fc.FontSize] = face
	}
	return face.Metrics()
}
