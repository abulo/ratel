package img

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"os"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/img/fontx"
	"github.com/abulo/ratel/v3/util"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/h2non/bimg"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

const (
	// DefaultDPI 为文字绘制使用的默认 DPI，固定为 72。
	DefaultDPI = 72
)

// Image 包装一张待处理的源图，图像变换方法就地修改内部源图并返回自身，支持链式调用。
//
// 零值内部源图为 nil，须通过 OpenLocalFile、OpenReader 或 NewImage 构造；
// 方法间不保证并发安全。
type Image struct {
	src image.Image
}

// OpenLocalFile 从本地文件解码图像并返回 Image；仅支持 jpg、png、webp、jpeg 格式，
// 不支持的格式或读取、解码失败时返回错误。
func OpenLocalFile(filename string, opts ...imaging.DecodeOption) (*Image, error) {

	var (
		extension  string
		fileBuffer []byte
		err        error
	)

	if fileBuffer, err = os.ReadFile(filename); err != nil {
		return nil, err
	}

	extension = bimg.NewImage(fileBuffer).Type()
	extAry := []string{"jpg", "png", "webp", "jpeg"}
	if !util.InArray(extension, extAry) {
		return nil, errors.New("No Support file")
	}
	// imaging 无法直接解码 webp，先转换为临时 jpeg 文件再打开。
	newFilename := filename + ".jpg"
	if extension == "webp" {
		if !util.FileExists(newFilename) {
			buffer, err := bimg.Read(filename)
			if err != nil {
				return nil, err
			}
			newCovImage, err := bimg.NewImage(buffer).Convert(bimg.JPEG)
			if err != nil {
				return nil, err
			}
			err = bimg.Write(newFilename, newCovImage)
			if err != nil {
				return nil, err
			}
		}
		filename = newFilename
	}
	img, err := imaging.Open(filename, opts...)
	if err != nil {
		return nil, err
	}
	return &Image{img}, nil
}

// GetSource 返回内部持有的底层 image.Image。
func (img *Image) GetSource() image.Image {
	return img.src
}

// SetSource 将内部源图替换为 src，并返回自身以支持链式调用。
func (img *Image) SetSource(src image.Image) *Image {
	img.src = src
	return img
}

// Width 返回图像像素宽度。
func (img *Image) Width() int {
	return img.src.Bounds().Size().X
}

// Height 返回图像像素高度。
func (img *Image) Height() int {
	return img.src.Bounds().Size().Y
}

// OpenReader 从读取器解码图像并返回 Image；解码失败时返回错误。
func OpenReader(src io.Reader, opts ...imaging.DecodeOption) (*Image, error) {
	img, err := imaging.Decode(src, opts...)
	if err != nil {
		return nil, err
	}
	return &Image{img}, nil
}

// Clone 返回包装当前源图的新 Image；两个实例共享同一份底层图像数据。
func (img *Image) Clone() *Image {
	return &Image{img.src}
}

// Resize 使用指定重采样滤镜将图像缩放到给定宽高；宽或高为 0 时保持原图宽高比。
// 变换就地作用于内部源图并返回自身。
func (img *Image) Resize(width int, height int, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Resize(img.src, width, height, filter)
	return img
}

// Crop 按指定锚点从图像中裁出 width×height 的矩形区域，变换就地作用于内部源图并返回自身。
func (img *Image) Crop(width int, height int, anchor imaging.Anchor) *Image {
	img.src = imaging.CropAnchor(img.src, width, height, anchor)
	return img
}

// Blur 使用高斯模糊处理图像，sigma 必须为正数，数值越大模糊程度越高；返回自身。
func (img *Image) Blur(sigma float64) *Image {
	img.src = imaging.Blur(img.src, sigma)
	return img
}

// Gray 将图像转换为灰度图，返回自身。
func (img *Image) Gray() *Image {
	img.src = imaging.Grayscale(img.src)
	return img
}

// AdjustContrast 按 percentage 调整图像对比度，取值范围 (-100, 100)：
// 0 为原图，-100 为纯灰色；返回自身。
func (img *Image) AdjustContrast(percentage float64) *Image {
	img.src = imaging.AdjustContrast(img.src, percentage)
	return img
}

// Sharpen 按给定参数生成图像的锐化版本，参数语义与底层 imaging.Sharpen 一致，返回自身。
func (img *Image) Sharpen(percentage float64) *Image {
	img.src = imaging.Sharpen(img.src, percentage)
	return img
}

// Invert 对图像做反色（负片）处理，返回自身。
func (img *Image) Invert() *Image {
	img.src = imaging.Invert(img.src)
	return img
}

// Convolve3x3 使用给定的 3×3 卷积核对图像做卷积处理，返回自身。
func (img *Image) Convolve3x3(kernel [9]float64) *Image {
	img.src = imaging.Convolve3x3(img.src, kernel, nil)
	return img
}

// Convolve5x5 使用给定的 5×5 卷积核对图像做卷积处理，返回自身。
func (img *Image) Convolve5x5(kernel [25]float64) *Image {
	img.src = imaging.Convolve5x5(img.src, kernel, nil)
	return img
}

// AdjustBrightness 按 percentage 调整图像亮度，取值范围 (-100, 100)：
// 0 为原图，-100 为纯黑，100 为纯白；返回自身。
func (img *Image) AdjustBrightness(percentage float64) *Image {
	img.src = imaging.AdjustBrightness(img.src, percentage)
	return img
}

// AdjustGamma 对图像做伽马校正，gamma 必须为正数：1.0 为原图，小于 1.0 变暗，
// 大于 1.0 变亮；返回自身。
func (img *Image) AdjustGamma(gamma float64) *Image {
	img.src = imaging.AdjustGamma(img.src, gamma)
	return img
}

// AdjustSaturation 按 percentage 调整图像饱和度，取值范围 (-100, 100)：
// 0 为原图，100 饱和度加倍，-100 为灰度；返回自身。
func (img *Image) AdjustSaturation(percentage float64) *Image {
	img.src = imaging.AdjustSaturation(img.src, percentage)
	return img
}

// AdjustSigmoid 使用 S 形函数做非线性对比度调整，可保留高光与阴影细节。
// midpoint 为对比度中点，取值 0 到 1，通常取 0.5；factor 为对比度强弱，
// 通常在 (-10, 10) 之间，正值增强对比度，负值减弱；返回自身。
func (img *Image) AdjustSigmoid(midpoint float64, factor float64) *Image {
	img.src = imaging.AdjustSigmoid(img.src, midpoint, factor)
	return img
}

// Rotate 将图像按 angle 角度（单位：度）逆时针旋转，bgColor 填充旋转后露出的空白区域，返回自身。
func (img *Image) Rotate(angle float64, bgColor color.Color) *Image {
	img.src = imaging.Rotate(img.src, angle, bgColor)
	return img
}

// Transverse 依次对图像执行 imaging.Transverse 与 imaging.Rotate90 变换，返回自身。
func (img *Image) Transverse() *Image {
	img.src = imaging.Rotate90(imaging.Transverse(img.src))
	return img
}

// Fit 使用指定重采样滤镜等比缩放图像，使整体不超过给定的最大宽高，返回自身。
func (img *Image) Fit(width int, height int, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Fit(img.src, width, height, filter)
	return img
}

// Fill 将源图像按锚点缩放并裁剪填充到指定宽高，过程中保持比例不变形，返回自身。
func (img *Image) Fill(width int, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Fill(img.src, width, height, anchor, filter)
	return img
}

// Paste 将 top 粘贴到当前图像的 pos 位置，返回自身。
func (img *Image) Paste(top *Image, pos image.Point) *Image {
	img.src = imaging.Paste(img.src, top.src, pos)
	return img
}

// Overlay 在 pos 位置以 opacity（取值 0.0 到 1.0）将 top 叠加到当前图像，返回自身。
func (img *Image) Overlay(top *Image, pos image.Point, opacity float64) *Image {
	img.src = imaging.Overlay(img.src, top.src, pos, opacity)
	return img
}

// AddWaterMark 按锚点与边距计算水印位置，并以 opacity（0.0 到 1.0）叠加水印图，返回自身。
func (img *Image) AddWaterMark(watermark *Image, anchor imaging.Anchor, marginX int, marginY int, opacity float64) *Image {
	pot := CalculatePt(img.src.Bounds().Size(), watermark.GetSource().Bounds().Size(), anchor, marginX, marginY)
	img.src = imaging.Overlay(img.src, watermark.src, pot, opacity)
	return img
}

// Compress 在不改变图像尺寸的前提下按指定 JPEG 质量重新编码，quality 取值 1 到 100，
// 数值越高画质越好；编码失败时返回 nil 与错误，成功时返回自身。
func (img *Image) Compress(quality int) (*Image, error) {
	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, img.src, &jpeg.Options{Quality: quality}); err != nil {
		logger.Logger.Error("Compress error:", err)
		return nil, err
	}
	encoded, _ := imaging.Decode(&buffer)
	img.src = encoded
	return img, nil
}

// DrawText 在当前图像上绘制文字 content。
//
// fc 提供字体、字号与颜色，m 为对应字号的字体度量；anchor 决定文字对齐方式，
// marginX、marginY 为距最近边框的边距，对齐方式为 Center 等居中锚点时边距
// 参照左侧或上侧边框。绘制失败时返回图像自身与错误。
func (img *Image) DrawText(content string, fc *fontx.FontConfig, m font.Metrics, anchor imaging.Anchor, marginX int, marginY int) (*Image, error) {
	ctx := freetype.NewContext()
	ctx.SetDPI(DefaultDPI)
	ctx.SetFont(fc.Font)
	ctx.SetFontSize(fc.FontSize)
	ctx.SetClip(img.src.Bounds())
	ctx.SetDst((img.src.(any)).(draw.Image))
	ctx.SetSrc(image.NewUniform(fc.Color))
	ctx.SetHinting(font.HintingNone)

	offset := 0
	if anchor == imaging.Top || anchor == imaging.TopLeft {
		offset = m.Ascent.Ceil() - m.Descent.Ceil()
	} else if anchor == imaging.Left || anchor == imaging.Right || anchor == imaging.Center {
		offset = (m.Ascent.Ceil() - m.Descent.Ceil()) / 2
	} else if anchor == imaging.BottomLeft || anchor == imaging.Bottom || anchor == imaging.BottomRight {
		offset = -m.Descent.Ceil() + m.Descent.Ceil()
	}

	pot := CalculatePt2(img.src.Bounds().Max, image.Point{0, 0}, anchor, marginX, marginY)
	_, err := ctx.DrawString(content, freetype.Pt(pot.X, pot.Y+offset))
	if err != nil {
		return img, err
	}
	return img, nil
}

// Paste 将 img 粘贴到 background 的 pos 位置，返回合成后的新 Image。
func Paste(background Image, img Image, pos image.Point) *Image {
	return &Image{imaging.Paste(background.src, img.src, pos)}
}

// Overlay 在 pos 位置以 opacity（取值 0.0 到 1.0）将 img 叠加到 background，
// 返回合成后的新 Image。
func Overlay(background Image, img Image, pos image.Point, opacity float64) *Image {
	return &Image{imaging.Overlay(background.src, img.src, pos, opacity)}
}

// CalculatePt 按锚点 anchor 与边距 marginX、marginY，依据目标图与水印图尺寸计算水印落点坐标。
func CalculatePt(targetSize image.Point,
	watermark image.Point,
	anchor imaging.Anchor,
	marginX int, marginY int) image.Point {
	if anchor == imaging.Top {
		return image.Point{
			X: (targetSize.X - watermark.X) / 2,
			Y: marginY,
		}
	}
	if anchor == imaging.TopLeft {
		return image.Point{
			X: marginX,
			Y: marginY,
		}
	}
	if anchor == imaging.TopRight {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: marginY,
		}
	}
	if anchor == imaging.Bottom {
		return image.Point{
			X: (targetSize.X - watermark.X) / 2,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.BottomLeft {
		return image.Point{
			X: marginX,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.BottomRight {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.Left {
		return image.Point{
			X: marginX,
			Y: (targetSize.Y - watermark.Y) / 2,
		}
	}
	if anchor == imaging.Right {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: (targetSize.Y - watermark.Y) / 2,
		}
	}
	return image.Point{
		X: (targetSize.X - watermark.X) / 2,
		Y: (targetSize.Y - watermark.Y) / 2,
	}
}

// CalculatePt2 按锚点 anchor 与边距 marginX、marginY 计算文字绘制落点坐标，
// 居中锚点下边距会叠加到居中位置上。
func CalculatePt2(targetSize image.Point,
	watermark image.Point,
	anchor imaging.Anchor,
	marginX int, marginY int) image.Point {
	if anchor == imaging.Top {
		return image.Point{
			X: (targetSize.X-watermark.X)/2 + marginX,
			Y: marginY,
		}
	}
	if anchor == imaging.TopLeft {
		return image.Point{
			X: marginX,
			Y: marginY,
		}
	}
	if anchor == imaging.TopRight {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: marginY,
		}
	}
	if anchor == imaging.Bottom {
		return image.Point{
			X: (targetSize.X-watermark.X)/2 + marginX,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.BottomLeft {
		return image.Point{
			X: marginX,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.BottomRight {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: (targetSize.Y - watermark.Y) - marginY,
		}
	}
	if anchor == imaging.Left {
		return image.Point{
			X: marginX,
			Y: (targetSize.Y-watermark.Y)/2 + marginY,
		}
	}
	if anchor == imaging.Right {
		return image.Point{
			X: (targetSize.X - watermark.X) - marginX,
			Y: (targetSize.Y-watermark.Y)/2 + marginY,
		}
	}
	return image.Point{
		X: (targetSize.X-watermark.X)/2 + marginX,
		Y: (targetSize.Y-watermark.Y)/2 + marginY,
	}
}

// SaveToFile 按文件扩展名推断编码格式，将图像写入 filename 指定的文件；
// 格式识别、创建文件或编码失败时返回错误；编码成功后若文件关闭失败，返回关闭错误。
func SaveToFile(img *Image, filename string, opts ...imaging.EncodeOption) error {
	f, err := imaging.FormatFromFilename(filename)
	if err != nil {
		return err
	}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = imaging.Encode(out, img.src, f, opts...)
	errc := out.Close()
	if err == nil {
		err = errc
	}
	return err
}

// Save 按指定格式将图像编码写入 out；编码失败时返回错误。
func Save(img *Image, out io.Writer, format imaging.Format, opts ...imaging.EncodeOption) error {
	return imaging.Encode(out, img.src, format, opts...)
}

// NewRGBA 按矩形区域新建一张 image.RGBA 图像。
func NewRGBA(r image.Rectangle) *image.RGBA {
	return image.NewRGBA(r)
}

// NewImage 将给定 image.Image 包装为 Image。
func NewImage(src image.Image) *Image {
	return &Image{src}
}
