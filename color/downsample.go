package color

import "math"

func Upsample(channel [][]uint8, targetWidth, targetHeight int) [][]uint8 {
	srcHeight := len(channel)
	if srcHeight == 0 {
		return make([][]uint8, targetHeight)
	}
	srcWidth := len(channel[0])

	upsampled := make([][]uint8, targetHeight)
	for i := range upsampled {
		upsampled[i] = make([]uint8, targetWidth)
	}

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcY := float64(y) / 2.0
			srcX := float64(x) / 2.0

			y0 := int(math.Floor(srcY))
			y1 := int(math.Min(float64(srcHeight-1), math.Ceil(srcY)))
			x0 := int(math.Floor(srcX))
			x1 := int(math.Min(float64(srcWidth-1), math.Ceil(srcX)))

			// Билинейная интерполяция
			if y0 == y1 && x0 == x1 {
				upsampled[y][x] = channel[y0][x0]
			} else if y0 == y1 {
				dx := srcX - float64(x0)
				upsampled[y][x] = uint8(float64(channel[y0][x0])*(1-dx) +
					float64(channel[y0][x1])*dx)
			} else if x0 == x1 {
				dy := srcY - float64(y0)
				upsampled[y][x] = uint8(float64(channel[y0][x0])*(1-dy) +
					float64(channel[y1][x0])*dy)
			} else {
				dx := srcX - float64(x0)
				dy := srcY - float64(y0)

				top := float64(channel[y0][x0])*(1-dx) + float64(channel[y0][x1])*dx
				bottom := float64(channel[y1][x0])*(1-dx) + float64(channel[y1][x1])*dx

				upsampled[y][x] = uint8(top*(1-dy) + bottom*dy)
			}
		}
	}
	return upsampled
}

func Downsample(channel [][]uint8, k int) [][]uint8 {
	height := len(channel)
	if height == 0 {
		return nil
	}
	width := len(channel[0])

	newHeight := (height + k - 1) / k
	newWidth := (width + k - 1) / k

	result := make([][]uint8, newHeight)
	for i := range result {
		result[i] = make([]uint8, newWidth)
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			var sum int
			var count int

			for iy := 0; iy < k; iy++ {
				sy := y*k + iy
				if sy >= height {
					continue
				}

				for ix := 0; ix < k; ix++ {
					sx := x*k + ix
					if sx >= width {
						continue
					}

					sum += int(channel[sy][sx])
					count++
				}
			}

			if count > 0 {
				result[y][x] = uint8(sum / count)
			}
		}
	}

	return result
}
