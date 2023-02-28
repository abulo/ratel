package img

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/img/fontx"
	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
)

func TestOpenLocalFile(t *testing.T) {
	type args struct {
		filename string
		opts     []imaging.DecodeOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenLocalFile(tt.args.filename, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenLocalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenLocalFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_GetSource(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want image.Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.GetSource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.GetSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_SetSource(t *testing.T) {
	type args struct {
		src image.Image
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.SetSource(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.SetSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Width(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Width(); got != tt.want {
				t.Errorf("Image.Width() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Height(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Height(); got != tt.want {
				t.Errorf("Image.Height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenReader(t *testing.T) {
	type args struct {
		src  io.Reader
		opts []imaging.DecodeOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenReader(tt.args.src, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Clone(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Resize(t *testing.T) {
	type args struct {
		width  int
		height int
		filter imaging.ResampleFilter
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Resize(tt.args.width, tt.args.height, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Resize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Crop(t *testing.T) {
	type args struct {
		width  int
		height int
		anchor imaging.Anchor
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Crop(tt.args.width, tt.args.height, tt.args.anchor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Crop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Blur(t *testing.T) {
	type args struct {
		sigma float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Blur(tt.args.sigma); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Blur() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Gray(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Gray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Gray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AdjustContrast(t *testing.T) {
	type args struct {
		percentage float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AdjustContrast(tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AdjustContrast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Sharpen(t *testing.T) {
	type args struct {
		percentage float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Sharpen(tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Sharpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Invert(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Invert(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Invert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Convolve3x3(t *testing.T) {
	type args struct {
		kernel [9]float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Convolve3x3(tt.args.kernel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Convolve3x3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Convolve5x5(t *testing.T) {
	type args struct {
		kernel [25]float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Convolve5x5(tt.args.kernel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Convolve5x5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AdjustBrightness(t *testing.T) {
	type args struct {
		percentage float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AdjustBrightness(tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AdjustBrightness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AdjustGamma(t *testing.T) {
	type args struct {
		gamma float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AdjustGamma(tt.args.gamma); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AdjustGamma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AdjustSaturation(t *testing.T) {
	type args struct {
		percentage float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AdjustSaturation(tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AdjustSaturation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AdjustSigmoid(t *testing.T) {
	type args struct {
		midpoint float64
		factor   float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AdjustSigmoid(tt.args.midpoint, tt.args.factor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AdjustSigmoid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Rotate(t *testing.T) {
	type args struct {
		angle   float64
		bgColor color.Color
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Rotate(tt.args.angle, tt.args.bgColor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Rotate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Transverse(t *testing.T) {
	tests := []struct {
		name string
		img  *Image
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Transverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Transverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Fit(t *testing.T) {
	type args struct {
		width  int
		height int
		filter imaging.ResampleFilter
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Fit(tt.args.width, tt.args.height, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Fit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Fill(t *testing.T) {
	type args struct {
		width  int
		height int
		anchor imaging.Anchor
		filter imaging.ResampleFilter
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Fill(tt.args.width, tt.args.height, tt.args.anchor, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Fill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Paste(t *testing.T) {
	type args struct {
		top *Image
		pos image.Point
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Paste(tt.args.top, tt.args.pos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Paste() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Overlay(t *testing.T) {
	type args struct {
		top     *Image
		pos     image.Point
		opacity float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.Overlay(tt.args.top, tt.args.pos, tt.args.opacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Overlay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_AddWaterMark(t *testing.T) {
	type args struct {
		watermark *Image
		anchor    imaging.Anchor
		marginX   int
		marginY   int
		opacity   float64
	}
	tests := []struct {
		name string
		img  *Image
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img.AddWaterMark(tt.args.watermark, tt.args.anchor, tt.args.marginX, tt.args.marginY, tt.args.opacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.AddWaterMark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_Compress(t *testing.T) {
	type args struct {
		quality int
	}
	tests := []struct {
		name    string
		img     *Image
		args    args
		want    *Image
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.img.Compress(tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_DrawText(t *testing.T) {
	type args struct {
		content string
		fc      *fontx.FontConfig
		m       font.Metrics
		anchor  imaging.Anchor
		marginX int
		marginY int
	}
	tests := []struct {
		name    string
		img     *Image
		args    args
		want    *Image
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.img.DrawText(tt.args.content, tt.args.fc, tt.args.m, tt.args.anchor, tt.args.marginX, tt.args.marginY)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.DrawText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.DrawText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaste(t *testing.T) {
	type args struct {
		background Image
		img        Image
		pos        image.Point
	}
	tests := []struct {
		name string
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Paste(tt.args.background, tt.args.img, tt.args.pos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Paste() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverlay(t *testing.T) {
	type args struct {
		background Image
		img        Image
		pos        image.Point
		opacity    float64
	}
	tests := []struct {
		name string
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Overlay(tt.args.background, tt.args.img, tt.args.pos, tt.args.opacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Overlay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePt(t *testing.T) {
	type args struct {
		targetSize image.Point
		watermark  image.Point
		anchor     imaging.Anchor
		marginX    int
		marginY    int
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePt(tt.args.targetSize, tt.args.watermark, tt.args.anchor, tt.args.marginX, tt.args.marginY); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePt2(t *testing.T) {
	type args struct {
		targetSize image.Point
		watermark  image.Point
		anchor     imaging.Anchor
		marginX    int
		marginY    int
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePt2(tt.args.targetSize, tt.args.watermark, tt.args.anchor, tt.args.marginX, tt.args.marginY); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePt2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveToFile(t *testing.T) {
	type args struct {
		img      *Image
		filename string
		opts     []imaging.EncodeOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveToFile(tt.args.img, tt.args.filename, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSave(t *testing.T) {
	type args struct {
		img    *Image
		format imaging.Format
		opts   []imaging.EncodeOption
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := Save(tt.args.img, out, tt.args.format, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Save() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestNewRGBA(t *testing.T) {
	type args struct {
		r image.Rectangle
	}
	tests := []struct {
		name string
		args args
		want *image.RGBA
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRGBA(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRGBA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewImage(t *testing.T) {
	type args struct {
		src image.Image
	}
	tests := []struct {
		name string
		args args
		want *Image
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImage(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
