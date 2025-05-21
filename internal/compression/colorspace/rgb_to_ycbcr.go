package colorspace

import (
	"AlgosSem4Lab2Neo/internal/domain/models"
	"image"
	"math"
)

func RGBToYCbCr(img image.Image) (*models.YCbCrData, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	y := make([]byte, width*height)
	cb := make([]byte, width*height)
	cr := make([]byte, width*height)

	for y_idx := 0; y_idx < height; y_idx++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y_idx+bounds.Min.Y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			yVal := 0.299*r8 + 0.587*g8 + 0.114*b8
			cbVal := -0.168736*r8 - 0.331264*g8 + 0.5*b8 + 128
			crVal := 0.5*r8 - 0.418688*g8 - 0.081312*b8 + 128

			y[y_idx*width+x] = clamp(yVal)
			cb[y_idx*width+x] = clamp(cbVal)
			cr[y_idx*width+x] = clamp(crVal)
		}
	}

	return &models.YCbCrData{
		Y:      y,
		Cb:     cb,
		Cr:     cr,
		Width:  width,
		Height: height,
	}, nil
}

func clamp(value float64) byte {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return byte(math.Round(value))
}
