package img

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/abulo/ratel/img/fontx"
	"github.com/abulo/ratel/util"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/h2non/bimg"
	"golang.org/x/image/font"
)

const (
	DefaultDPI = 72
)

type Image struct {
	src image.Image
}

// OpenLocalFile gets an Image using a local file.
func OpenLocalFile(filename string, opts ...imaging.DecodeOption) (*Image, error) {
	ext := filepath.Ext(filename)
	extension := strings.ToLower(strings.TrimPrefix(ext, "."))
	extAry := []string{"jpg", "png", "webp", "jpeg"}
	if !util.InArray(extension, extAry) {
		return nil, errors.New("No Support file")
	}
	//如果是 webp 格式, 先将 webp 转换成 jpeg
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

// GetSource returns source image.Image.
func (img *Image) GetSource() image.Image {
	return img.src
}

func (img *Image) SetSource(src image.Image) *Image {
	img.src = src
	return img
}

//Width 获取文件的宽度
func (img *Image) Width() int {
	return img.src.Bounds().Size().X
}

//Width 获取文件的高度
func (img *Image) Height() int {
	// img.src.Decode
	return img.src.Bounds().Size().Y
}

// OpenReader gets an Image using a reader.
func OpenReader(src io.Reader, opts ...imaging.DecodeOption) (*Image, error) {
	img, err := imaging.Decode(src, opts...)
	if err != nil {
		return nil, err
	}
	return &Image{img}, nil
}

// Clone clones and returns a new Image.
func (img *Image) Clone() *Image {
	return &Image{img.src}
}

// Resize resizes the image to the specified width and height using the specified resampling filter
// and returns the transformed image. If one of width or height is 0, the image aspect ratio is preserved.
//
// Example:
//
//  dstImage := imaging.Resize(srcImage, 800, 600, imaging.Lanczos)
// 调整大小使用指定的重采样过滤器将图像调整为指定的宽度和高度，并返回变换后的图像。 如果宽度或高度中的一个为0，则保留图像宽高比。
func (img *Image) Resize(width int, height int, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Resize(img.src, width, height, filter)
	return img
}

// Crop cuts out a rectangular region with the specified size from the image using
// the specified anchor point and returns the cropped image.
// Example:
//
//  imaging.CropAnchor(src, 100, 100, imaging.TopLeft)
// 使用指定的锚点从图像中剪切出具有指定大小的矩形区域，并返回裁剪后的图像。
func (img *Image) Crop(width int, height int, anchor imaging.Anchor) *Image {
	img.src = imaging.CropAnchor(img.src, width, height, anchor)
	return img
}

// Blur produces a blurred version of the image using a Gaussian function.
// Sigma parameter must be positive and indicates how much the image will be blurred.
//
// Example:
//
//  dstImage := imaging.Blur(srcImage, 3.5)
// 模糊使用高斯函数生成图像的模糊版本。 Sigma参数必须为正，表示图像模糊的程度。
func (img *Image) Blur(sigma float64) *Image {
	img.src = imaging.Blur(img.src, sigma)
	return img
}

// Gray produces a grayscale version of the image.
// 使图像变为黑白。
func (img *Image) Gray() *Image {
	img.src = imaging.Grayscale(img.src)
	return img
}

// AdjustContrast changes the contrast of the image using the percentage parameter and returns the adjusted image.
// The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid gray image.
// Examples:
// Decrease image contrast by 10%.
//  dstImage = imaging.AdjustContrast(srcImage, -10)
// Increase image con
//  dstImage = imaging.AdjustContrast(srcImage, 20)
// 使用百分比参数更改图像的对比度并返回调整后的图像。
//
// 百分比必须在范围内（-100,100）。
//
// 百分比= 0给出原始图像。 百分比= -100给出纯灰色图像。
func (img *Image) AdjustContrast(percentage float64) *Image {
	img.src = imaging.AdjustContrast(img.src, percentage)
	return img
}

// Sharpen produces a sharpened version of the image. Sigma parameter must be positive and indicates how much the image will be sharpened.
// Example:
//  dstImage := imaging.Sharpen(srcImage, 3.5)
//锐化生成图像的锐化版本。 Sigma参数必须为正，表示图像将被锐化多少。
func (img *Image) Sharpen(percentage float64) *Image {
	img.src = imaging.Sharpen(img.src, percentage)
	return img
}

// Invert produces an inverted (negated) version of the image.
//
// 反转产生图像的反转（否定）版本。
func (img *Image) Invert() *Image {
	img.src = imaging.Invert(img.src)
	return img
}

// Convolve3x3 convolves the image with the specified 3x3 convolution kernel.
// Default parameters are used if a nil *ConvolveOptions is passed.
//
// Convolve3x3使用指定的3x3卷积内核对图像进行卷积。
// 如果传递了nil * ConvolveOptions，则使用默认参数。
func (img *Image) Convolve3x3(kernel [9]float64) *Image {
	img.src = imaging.Convolve3x3(img.src, kernel, nil)
	return img
}

// Convolve5x5 convolves the image with the specified 5x5 convolution kernel.
// Default parameters are used if a nil *ConvolveOptions is passed.
//
// Convolve5x5使用指定的5x5卷积内核对图像进行卷积。
// 如果传递了nil * ConvolveOptions，则使用默认参数。
func (img *Image) Convolve5x5(kernel [25]float64) *Image {
	img.src = imaging.Convolve5x5(img.src, kernel, nil)
	return img
}

// AdjustBrightness changes the brightness of the image using the percentage parameter and returns the adjusted image. The percentage must be in range (-100, 100). The percentage = 0 gives the original image. The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
// Examples:
// Decrease image brightness by 15%.
//  dstImage = imaging.AdjustBrightness(srcImage, -15)
// Increase image brightness by 10%.
//  dstImage = imaging.AdjustBrightness(srcImage, 10)
// AdjustBrightness使用百分比参数更改图像的亮度，并返回调整后的图像。
// 百分比必须在范围内（-100,100）。 百分比= 0给出原始图像。
// 百分比= -100给出纯黑色图像。 百分比= 100给出纯白图像。
func (img *Image) AdjustBrightness(percentage float64) *Image {
	img.src = imaging.AdjustBrightness(img.src, percentage)
	return img
}

// AdjustGamma performs a gamma correction on the image and returns the adjusted image. Gamma parameter must be positive. Gamma = 1.0 gives the original image. Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.
// Example:
//  dstImage = imaging.AdjustGamma(srcImage, 0.7)
// AdjustGamma对图像执行伽玛校正并返回调整后的图像。
// Gamma参数必须为正数。 Gamma = 1.0给出原始图像。
// 小于1.0的伽玛使图像变暗，大于1.0的伽玛使其变亮。
func (img *Image) AdjustGamma(gamma float64) *Image {
	img.src = imaging.AdjustGamma(img.src, gamma)
	return img
}

// AdjustSaturation changes the saturation of the image using the percentage parameter
// and returns the adjusted image. The percentage must be in the range (-100, 100).
// The percentage = 0 gives the original image. The percentage = 100 gives the image
// with the saturation value doubled for each pixel. The percentage = -100 gives the image
// with the saturation value zeroed for each pixel (grayscale).
// Examples:
// Increase image saturation by 25%.
//  dstImage = imaging.AdjustSaturation(srcImage, 25)
// Decrease image saturation by 10%.
//  dstImage = imaging.AdjustSaturation(srcImage, -10)
// AdjustSaturation使用百分比参数更改图像的饱和度并返回调整后的图像。 百分比必须在范围（-100,100）内。 百分比= 0给出原始图像。
// 百分比= 100表示每个像素的饱和度值加倍的图像。 百分比= -100给出的图像的饱和度值为每个像素（灰度）。
func (img *Image) AdjustSaturation(percentage float64) *Image {
	img.src = imaging.AdjustSaturation(img.src, percentage)
	return img
}

// AdjustSigmoid changes the contrast of the image using a sigmoidal function and returns the adjusted image. It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail. The midpoint parameter is the midpoint of contrast that must be between 0 and 1, typically 0.5. The factor parameter indicates how much to increase or decrease the contrast, typically in range (-10, 10). If the factor parameter is positive the image contrast is increased otherwise the contrast is decreased.
// Examples:
// Increase the contrast.
//  dstImage = imaging.AdjustSigmoid(srcImage, 0.5, 3.0)
// Decrease the contrast.
//  dstImage = imaging.AdjustSigmoid(srcImage, 0.5, -3.0)
// AdjustSigmoid使用S形函数更改图像的对比度并返回调整后的图像。
// 这是一种非线性对比度变化，可用于照片调整，因为它可以保留高光和阴影细节。
// 中点参数是对比度的中点，必须介于0和1之间，通常为0.5。 因子参数表示增加或减少对比度的程度，
// 通常在范围（-10,10）内。 如果因子参数为正，则图像对比度增加，否则对比度降低
func (img *Image) AdjustSigmoid(midpoint float64, factor float64) *Image {
	img.src = imaging.AdjustSigmoid(img.src, midpoint, factor)
	return img
}

// Rotate rotates an image by the given angle counter-clockwise .
// The angle parameter is the rotation angle in degrees.
// The bgColor parameter specifies the color of the uncovered zone after the rotation.
//
// 旋转将图像逆时针旋转给定角度。 角度参数是以度为单位的旋转角度。 bgColor参数指定旋转后未覆盖区域的颜色。
func (img *Image) Rotate(angle float64, bgColor color.Color) *Image {
	img.src = imaging.Rotate(img.src, angle, bgColor)
	return img
}

// Transverse flips the image vertically.
func (img *Image) Transverse() *Image {
	img.src = imaging.Rotate90(imaging.Transverse(img.src))
	return img
}

// Fit scales down the image using the specified resample filter to fit the specified maximum width and height and returns the transformed image.
// Example:
//  dstImage := imaging.Fit(srcImage, 800, 600, imaging.Lanczos)
// 使用指定的重采样滤镜按比例缩小图像以适合指定的最大宽度和高度，并返回变换后的图像。
func (img *Image) Fit(width int, height int, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Fit(img.src, width, height, filter)
	return img
}

// Fill creates an image with the specified dimensions and fills it with the scaled source image.
// To achieve the correct aspect ratio without stretching, the source image will be cropped.
// Example:
//  dstImage := imaging.Fill(srcImage, 800, 600, imaging.Center, imaging.Lanczos)
// Fill创建具有指定尺寸的图像，并使用缩放的源图像填充它。 要在不拉伸的情况下获得正确的宽高比，将裁剪源图像。
func (img *Image) Fill(width int, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) *Image {
	img.src = imaging.Fill(img.src, width, height, anchor, filter)
	return img
}

// Paste pastes the an image to this image at the specified position.
// Example:
//  imaging.Paste(src1, src2, image.Pt(10, 100))
// 粘贴将img图像粘贴到指定位置的背景图像并返回组合图像。
func (img *Image) Paste(top *Image, pos image.Point) *Image {
	img.src = imaging.Paste(img.src, top.src, pos)
	return img
}

// Overlay draws an image over the background image at given position.
// Opacity parameter is the opacity of the img
// image layer, used to compose the images, it must be from 0.0 to 1.0.
//
// Examples:
//
//	dstImage := imaging.Overlay(backgroundImage, spriteImage, image.Pt(50, 50), 1.0)
//
//	dstImage := imaging.Overlay(imageOne, imageTwo, image.Pt(0, 0), 0.5)
// 叠加在给定位置的背景图像上绘制img图像并返回合成图像。 不透明度参数是img图像层的不透明度，用于构成图像，它必须从0.0到1.0。
func (img *Image) Overlay(top *Image, pos image.Point, opacity float64) *Image {
	img.src = imaging.Overlay(img.src, top.src, pos, opacity)
	return img
}

func (img *Image) AddWaterMark(watermark *Image, anchor imaging.Anchor, marginX int, marginY int, opacity float64) *Image {
	pot := CalculatePt(img.src.Bounds().Size(), watermark.GetSource().Bounds().Size(), anchor, marginX, marginY)
	// render watermark.
	img.src = imaging.Overlay(img.src, watermark.src, pot, opacity)
	return img
}

// JPEGQuality compressions jpeg without change the image size.
//
// Quality ranges from 1 to 100 inclusive, higher is better.
//
// 在不改变图片尺寸的情况下压缩JPEG图像。图像质量为1-100
func (img *Image) Compress(quality int) *Image {
	var buffer bytes.Buffer
	jpeg.Encode(&buffer, img.src, &jpeg.Options{Quality: quality})
	encoded, _ := imaging.Decode(&buffer)
	img.src = encoded
	return img
}

// DrawText draws text on the image.
// anchor is text align,
// marginX and marginY is the margin to nearest border, if the nearest border is not clear, such as imaging.Center,
// marginX and marginY always reference to the left border or the top border.
// example:
//  im, _ := img.OpenLocalFile("E:\\test\\1.jpg") // 1900x1283
//	fo, _ := fontx.LoadFont("E:\\test\\Inkfree.ttf")
//	fc := &fontx.FontConfig{
//		Font:     fo.Font,
//		FontSize: 200,
//		Color:    color.Black,
//	}
//	metrics := fo.GetMetrics(fc)
//	im.DrawText("Hello", fc, metrics, imaging.BottomRight, 500, 700)
// 尝试在image上绘制文字。
// 参数anchor是文字的对齐方式，marginX和marginY是文字的边距。
// 例如，当anchor=imaging.BottomRight，此时marginX是距离底边的距离，marginY是距离有边框的距离。
// 当anchor=imaging.BottomRight，此时marginX是距离底边的距离，marginY是距离有边框的距离。
func (img *Image) DrawText(content string, fc *fontx.FontConfig, m font.Metrics, anchor imaging.Anchor, marginX int, marginY int) (*Image, error) {
	ctx := freetype.NewContext()
	ctx.SetDPI(DefaultDPI)
	ctx.SetFont(fc.Font)
	ctx.SetFontSize(fc.FontSize)
	ctx.SetClip(img.src.Bounds())
	//overPaintImage := image.NewRGBA(img.src.Bounds())
	//draw.Draw(overPaintImage, img.src.Bounds(), img.src, image.ZP, draw.Over)
	ctx.SetDst((img.src.(interface{})).(draw.Image))
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

// Paste pastes the img image to the background image at the specified position and returns the combined image.
func Paste(background Image, img Image, pos image.Point) *Image {
	return &Image{imaging.Paste(background.src, img.src, pos)}
}

// Overlay draws the img image over the background image at given position and returns the combined image.
// Opacity parameter is the opacity of the img image layer, used to compose the images, it must be from 0.0 to 1.0.
func Overlay(background Image, img Image, pos image.Point, opacity float64) *Image {
	return &Image{imaging.Overlay(background.src, img.src, pos, opacity)}
}

// CalculatePt calculates point according to the given point.
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

// CalculatePt calculates point according to the given point.
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

// SaveToFile saves img.Image to file.
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

// Save writes img.Image to a writer.
func Save(img *Image, out io.Writer, format imaging.Format, opts ...imaging.EncodeOption) error {
	return imaging.Encode(out, img.src, format, opts...)
}

// NewRGBA creates a new image.RGBA
func NewRGBA(r image.Rectangle) *image.RGBA {
	return image.NewRGBA(r)
}

// NewImage creates a new img.Image
func NewImage(src image.Image) *Image {
	return &Image{src}
}
