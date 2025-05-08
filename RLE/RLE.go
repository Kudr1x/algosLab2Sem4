package RLE

import (
	"fmt"
	"math"
)

func RunLengthEncode(ac []int) [][]int {
	if len(ac) == 0 {
		return [][]int{}
	}
	result := make([][]int, 0)
	zeroCount := 0
	for i := 0; i < len(ac); i++ {
		if ac[i] == 0 {
			zeroCount++
			if zeroCount == 16 {
				result = append(result, []int{15, 0})
				zeroCount = 0
			}
		} else {
			category := 0
			absVal := int(math.Abs(float64(ac[i])))
			for tmpVal := absVal; tmpVal > 0; tmpVal >>= 1 {
				category++
			}
			result = append(result, []int{zeroCount, category})
			result = append(result, []int{ac[i]})
			zeroCount = 0
		}
	}
	if zeroCount > 0 { // EOB
		result = append(result, []int{0, 0})
	}
	return result
}

func RunLengthDecode(rleData [][]int) []int {
	numACCoeffs := 63
	acCoefficients := make([]int, 0, numACCoeffs)
	valueExpected := false

	for i := 0; i < len(rleData); i++ {
		if len(acCoefficients) >= numACCoeffs {
			break
		}

		pair := rleData[i]
		if valueExpected {
			if len(pair) > 0 {
				acCoefficients = append(acCoefficients, pair[0])
			} else {
				fmt.Println("Предупреждение: Ожидалось значение, но найден пустой массив")
			}
			valueExpected = false
			continue
		}

		if len(pair) < 2 {
			fmt.Println("Предупреждение: Некорректная RLE пара, пропускается")
			continue
		}

		run := pair[0]
		size := pair[1]

		if run == 0 && size == 0 {
			break
		}

		if run == 15 && size == 0 {
			zerosToAdd := 16
			if len(acCoefficients)+zerosToAdd > numACCoeffs {
				fmt.Printf("Предупреждение: ZRL приведет к превышению лимита (%d + %d > %d). Добавляем только %d нулей.\n",
					len(acCoefficients), zerosToAdd, numACCoeffs, numACCoeffs-len(acCoefficients))
				zerosToAdd = numACCoeffs - len(acCoefficients)
			}
			for k := 0; k < zerosToAdd; k++ {
				acCoefficients = append(acCoefficients, 0)
			}
		} else {
			if run < 0 || run > 15 {
				fmt.Printf("Предупреждение: Некорректный run_length=%d в RLE паре. Игнорируется пара.\n", run)
				continue
			}

			zerosToAdd := run
			if len(acCoefficients)+zerosToAdd > numACCoeffs {
				fmt.Printf("Предупреждение: Run-length %d привел бы к превышению лимита (%d + %d > %d). Добавляем только %d нулей.\n",
					run, len(acCoefficients), zerosToAdd, numACCoeffs, numACCoeffs-len(acCoefficients))
				zerosToAdd = numACCoeffs - len(acCoefficients)
			}
			for k := 0; k < zerosToAdd; k++ {
				acCoefficients = append(acCoefficients, 0)
			}

			if len(acCoefficients) < numACCoeffs {
				valueExpected = true
			} else if size != 0 {
				fmt.Printf("Предупреждение: После добавления %d нулей не осталось места для значения (достигнут лимит %d). Значение пропущено.\n",
					run, numACCoeffs)
			}
		}
	}

	if len(acCoefficients) < numACCoeffs {
		remaining := numACCoeffs - len(acCoefficients)
		for i := 0; i < remaining; i++ {
			acCoefficients = append(acCoefficients, 0)
		}
	}

	return acCoefficients[:numACCoeffs]
}
