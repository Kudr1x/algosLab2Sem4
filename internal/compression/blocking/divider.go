package blocking

import (
	"AlgosSem4Lab2Neo/internal/domain/models"
	"errors"
)

func DivideIntoBlocks(data []byte, width, height, blockSize int, channel string) ([]*models.Block, error) {
	if blockSize <= 0 {
		return nil, errors.New("block size must be positive")
	}

	blocksX := (width + blockSize - 1) / blockSize
	blocksY := (height + blockSize - 1) / blockSize

	blocks := make([]*models.Block, blocksX*blocksY)

	for by := 0; by < blocksY; by++ {
		for bx := 0; bx < blocksX; bx++ {
			block := &models.Block{
				Data:    make([]float64, blockSize*blockSize),
				Size:    blockSize,
				Channel: channel,
			}

			for y := 0; y < blockSize; y++ {
				for x := 0; x < blockSize; x++ {
					imgX := bx*blockSize + x
					imgY := by*blockSize + y

					blockIdx := y*blockSize + x

					if imgX < width && imgY < height {
						imgIdx := imgY*width + imgX
						block.Data[blockIdx] = float64(data[imgIdx])
					} else {
						block.Data[blockIdx] = 0
					}
				}
			}

			blocks[by*blocksX+bx] = block
		}
	}

	return blocks, nil
}

// ReconstructFromBlocks восстанавливает данные канала из блоков
func ReconstructFromBlocks(blocks []*models.Block, width, height int) ([]byte, error) {
	if len(blocks) == 0 {
		return nil, errors.New("no blocks provided")
	}

	blockSize := blocks[0].Size

	blocksX := (width + blockSize - 1) / blockSize
	blocksY := (height + blockSize - 1) / blockSize

	result := make([]byte, width*height)

	for by := 0; by < blocksY; by++ {
		for bx := 0; bx < blocksX; bx++ {
			block := blocks[by*blocksX+bx]

			for y := 0; y < blockSize; y++ {
				for x := 0; x < blockSize; x++ {
					imgX := bx*blockSize + x
					imgY := by*blockSize + y

					if imgX >= width || imgY >= height {
						continue
					}

					blockIdx := y*blockSize + x

					imgIdx := imgY*width + imgX

					value := int(block.Data[blockIdx] + 0.5)
					if value < 0 {
						value = 0
					} else if value > 255 {
						value = 255
					}

					result[imgIdx] = byte(value)
				}
			}
		}
	}

	return result, nil
}
