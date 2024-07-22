package utils

import "math"

func ToAmericanOdds(value int) float64 {
	decimalOdds := float64(value) / 1000.0
	if decimalOdds < 1.00 {
		return 0
	}

	var americanOdds float64
	if decimalOdds >= 2.00 {
		americanOdds = (decimalOdds - 1) * 100
	} else {
		americanOdds = -100 / (decimalOdds - 1)
	}

	return math.Round(americanOdds*100) / 100 // Round to 2 decimal places
}
