package colorspace

import (
	"image"
	"image/color"
)

func YCbCrToRGB(y, cb, cr []byte, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y_idx := 0; y_idx < height; y_idx++ {
		for x := 0; x < width; x++ {
			idx := y_idx*width + x

			yVal := float64(y[idx])
			cbVal := float64(cb[idx]) - 128
			crVal := float64(cr[idx]) - 128

			r := yVal + 1.402*crVal
			g := yVal - 0.344136*cbVal - 0.714136*crVal
			b := yVal + 1.772*cbVal

			img.Set(x, y_idx, color.RGBA{
				R: clamp(r),
				G: clamp(g),
				B: clamp(b),
				A: 255,
			})
		}
	}

	return img
}
