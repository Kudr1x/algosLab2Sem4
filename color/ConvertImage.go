package color

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type YCbCrImage struct {
	Y, Cb, Cr     [][]uint8
	Width, Height int
}

func ConvertImage(fileName string) ([]image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось декодировать PNG: %v", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	imgGS := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))
			imgGS.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	imgBW := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := imgGS.GrayAt(x, y).Y
			if gray > 127 {
				imgBW.SetGray(x, y, color.Gray{Y: 255})
			} else {
				imgBW.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	imgBWNoDither := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := imgGS.GrayAt(x, y).Y
			if gray > 127 {
				imgBWNoDither.SetGray(x, y, color.Gray{Y: 255})
			} else {
				imgBWNoDither.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return []image.Image{img, imgGS, imgBW, imgBWNoDither}, nil
}

func RGBToYCbCr(r, g, b uint8) (y, cb, cr uint8) {
	yFloat := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	cbFloat := -0.1687*float64(r) - 0.3313*float64(g) + 0.5*float64(b) + 128
	crFloat := 0.5*float64(r) - 0.4187*float64(g) - 0.0813*float64(b) + 128

	yVal := math.Max(0, math.Min(255, yFloat))
	cbVal := math.Max(0, math.Min(255, cbFloat))
	crVal := math.Max(0, math.Min(255, crFloat))

	return uint8(yVal), uint8(cbVal), uint8(crVal)
}

func ConvertToYCbCr(img image.Image) *YCbCrImage {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	y := make([][]uint8, height)
	cb := make([][]uint8, height)
	cr := make([][]uint8, height)

	for i := range y {
		y[i] = make([]uint8, width)
		cb[i] = make([]uint8, width)
		cr[i] = make([]uint8, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			y[i][j], cb[i][j], cr[i][j] = RGBToYCbCr(r8, g8, b8)
		}
	}

	return &YCbCrImage{
		Y:      y,
		Cb:     cb,
		Cr:     cr,
		Width:  width,
		Height: height,
	}
}

func YCbCrToRGB(y, cb, cr uint8) (r, g, b uint8) {
	rFloat := float64(y) + 1.402*(float64(cr)-128.0)
	gFloat := float64(y) - 0.34414*(float64(cb)-128.0) - 0.71414*(float64(cr)-128.0)
	bFloat := float64(y) + 1.772*(float64(cb)-128.0)

	r = uint8(math.Max(0, math.Min(255, math.Round(rFloat))))
	g = uint8(math.Max(0, math.Min(255, math.Round(gFloat))))
	b = uint8(math.Max(0, math.Min(255, math.Round(bFloat))))
	return
}
