package DCT

import "math"

func DCT2D(block [][]float64) [][]float64 {
	result := make([][]float64, 8)
	for i := range result {
		result[i] = make([]float64, 8)
	}

	for u := 0; u < 8; u++ {
		for v := 0; v < 8; v++ {
			var sum float64
			cu := 1.0
			if u == 0 {
				cu = 1.0 / math.Sqrt(2.0)
			}
			cv := 1.0
			if v == 0 {
				cv = 1.0 / math.Sqrt(2.0)
			}
			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					sum += block[x][y] *
						math.Cos((float64(2*x+1)*float64(u)*math.Pi)/16.0) *
						math.Cos((float64(2*y+1)*float64(v)*math.Pi)/16.0)
				}
			}
			result[u][v] = 0.25 * cu * cv * sum
		}
	}
	return result
}

func IDCT2D(block [][]float64) [][]float64 {
	result := make([][]float64, 8)
	for i := range result {
		result[i] = make([]float64, 8)
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			var sum float64
			for u := 0; u < 8; u++ {
				for v := 0; v < 8; v++ {
					cu := 1.0
					if u == 0 {
						cu = 1.0 / math.Sqrt(2.0)
					}
					cv := 1.0
					if v == 0 {
						cv = 1.0 / math.Sqrt(2.0)
					}
					sum += cu * cv * block[u][v] *
						math.Cos((float64(2*x+1)*float64(u)*math.Pi)/16.0) *
						math.Cos((float64(2*y+1)*float64(v)*math.Pi)/16.0)
				}
			}
			result[x][y] = 0.25 * sum
		}
	}
	return result
}

func Idct2DBlocks(dctBlocks [][][]float64) [][][]float64 {
	blocks := make([][][]float64, len(dctBlocks))
	for i, block := range dctBlocks {
		blocks[i] = IDCT2D(block)
	}
	return blocks
}
