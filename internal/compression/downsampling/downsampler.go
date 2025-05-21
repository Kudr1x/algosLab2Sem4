package downsampling

import (
	"AlgosSem4Lab2Neo/internal/domain/models"
	"image"
	"math"
)

func Downsample(ycbcr *models.YCbCrData, ratio string) error {
	width := ycbcr.Width
	height := ycbcr.Height

	switch ratio {
	case "4:4:4":
		ycbcr.CbWidth = width
		ycbcr.CbHeight = height
		ycbcr.CrWidth = width
		ycbcr.CrHeight = height
		ycbcr.SubsamplingRatio = "4:4:4"

	case "4:2:2":
		ycbcr.CbWidth = int(math.Ceil(float64(width) / 2))
		ycbcr.CbHeight = height
		ycbcr.CrWidth = int(math.Ceil(float64(width) / 2))
		ycbcr.CrHeight = height
		ycbcr.SubsamplingRatio = "4:2:2"

		newCb := make([]byte, ycbcr.CbWidth*ycbcr.CbHeight)
		newCr := make([]byte, ycbcr.CrWidth*ycbcr.CrHeight)

		for y := 0; y < height; y++ {
			for x := 0; x < ycbcr.CbWidth; x++ {
				x2 := x * 2
				if x2+1 < width {
					cbSum := int(ycbcr.Cb[y*width+x2]) + int(ycbcr.Cb[y*width+x2+1])
					crSum := int(ycbcr.Cr[y*width+x2]) + int(ycbcr.Cr[y*width+x2+1])
					newCb[y*ycbcr.CbWidth+x] = byte(cbSum / 2)
					newCr[y*ycbcr.CrWidth+x] = byte(crSum / 2)
				} else {
					newCb[y*ycbcr.CbWidth+x] = ycbcr.Cb[y*width+x2]
					newCr[y*ycbcr.CrWidth+x] = ycbcr.Cr[y*width+x2]
				}
			}
		}

		ycbcr.Cb = newCb
		ycbcr.Cr = newCr

	case "4:2:0":
		ycbcr.CbWidth = int(math.Ceil(float64(width) / 2))
		ycbcr.CbHeight = int(math.Ceil(float64(height) / 2))
		ycbcr.CrWidth = int(math.Ceil(float64(width) / 2))
		ycbcr.CrHeight = int(math.Ceil(float64(height) / 2))
		ycbcr.SubsamplingRatio = "4:2:0"

		newCb := make([]byte, ycbcr.CbWidth*ycbcr.CbHeight)
		newCr := make([]byte, ycbcr.CrWidth*ycbcr.CrHeight)

		for y := 0; y < ycbcr.CbHeight; y++ {
			for x := 0; x < ycbcr.CbWidth; x++ {
				x2 := x * 2
				y2 := y * 2

				cbSum := 0
				crSum := 0
				count := 0

				for dy := 0; dy < 2; dy++ {
					if y2+dy < height {
						for dx := 0; dx < 2; dx++ {
							if x2+dx < width {
								cbSum += int(ycbcr.Cb[(y2+dy)*width+(x2+dx)])
								crSum += int(ycbcr.Cr[(y2+dy)*width+(x2+dx)])
								count++
							}
						}
					}
				}

				newCb[y*ycbcr.CbWidth+x] = byte(cbSum / count)
				newCr[y*ycbcr.CrWidth+x] = byte(crSum / count)
			}
		}

		ycbcr.Cb = newCb
		ycbcr.Cr = newCr

	default:
		return image.ErrFormat
	}

	return nil
}

func Upsample(ycbcr *models.YCbCrData) error {
	width := ycbcr.Width
	height := ycbcr.Height

	switch ycbcr.SubsamplingRatio {
	case "4:4:4":
		return nil

	case "4:2:2":
		newCb := make([]byte, width*height)
		newCr := make([]byte, width*height)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				srcX := x / 2
				newCb[y*width+x] = ycbcr.Cb[y*ycbcr.CbWidth+srcX]
				newCr[y*width+x] = ycbcr.Cr[y*ycbcr.CrWidth+srcX]
			}
		}

		ycbcr.Cb = newCb
		ycbcr.Cr = newCr
		ycbcr.CbWidth = width
		ycbcr.CbHeight = height
		ycbcr.CrWidth = width
		ycbcr.CrHeight = height

	case "4:2:0":
		newCb := make([]byte, width*height)
		newCr := make([]byte, width*height)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				srcX := x / 2
				srcY := y / 2
				newCb[y*width+x] = ycbcr.Cb[srcY*ycbcr.CbWidth+srcX]
				newCr[y*width+x] = ycbcr.Cr[srcY*ycbcr.CrWidth+srcX]
			}
		}

		ycbcr.Cb = newCb
		ycbcr.Cr = newCr
		ycbcr.CbWidth = width
		ycbcr.CbHeight = height
		ycbcr.CrWidth = width
		ycbcr.CrHeight = height

	default:
		return image.ErrFormat
	}

	return nil
}
