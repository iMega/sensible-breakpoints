package main

import (
	"math"
)

const MIN_STRONG_VALUES = 8

func getNumberPartsFromOriginal(p, w int) int {
	return int(math.Floor(float64(w) / float64(p)))
}

func thereSufficientNumberStrongValues(numberParts int) bool {
	return numberParts >= MIN_STRONG_VALUES
}

func createStrongValues(min, max int) []int {
	var result []int
	diff := max - min

	p := int(math.Floor(float64(diff) / float64(MIN_STRONG_VALUES)))

	for i := 0; i < MIN_STRONG_VALUES; i++ {
		val := min + p*i
		if isEven(val) == false {
			val--
		}
		result = append(result, val)
	}

	logger("it created strong-values %d", result)

	return result
}

func calcStrongValues(min, max int) []int {
	var result []int
	points := getStrongValuesLessHalf(max, min)
	result = append(result, points...)

	minWidth := getMinimalWidth(points)
	result = append(result, getStrongValuesOverHalf(max, max/2, minWidth)...)

	return result
}

func getStrongValuesLessHalf(max, min int) []int {
	var (
		ret         []int
		strongValue = max
	)

	for {
		strongValue = strongValue / 2

		if strongValue <= min {
			return ret
		}

		if isEven(strongValue) == false {
			strongValue--
		}

		if strongValue > 0 {
			ret = append(ret, strongValue)
		}
	}
}

func getMinimalWidth(points []int) int {
	min := points[0]

	for _, p := range points {
		if p < min {
			min = p
		}
	}

	return int(min)
}

func getStrongValuesOverHalf(max, min, step int) []int {
	var (
		ret         []int
		strongValue = max
		i           int
	)

	for {
		i++
		strongValue = strongValue - step*i
		if strongValue < min {
			return ret
		}

		if isEven(strongValue) == false {
			strongValue--
		}

		ret = append(ret, strongValue)
	}
}

func isEven(n int) bool {
	return n%2 == 0
}
