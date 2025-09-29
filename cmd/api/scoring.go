package main

import "math"

const (
	RWTotal   = 54
	MathTotal = 44
)

func roundToNearest10(x float64) int {
	return int(math.Round(x/10.0) * 10)
}

func scaleToSection(rawCorrect int, sectionTotal int) int {
	if sectionTotal <= 0 {
		return 200
	}
	pct := float64(rawCorrect) / float64(sectionTotal)
	scaled := 200.0 + pct*600.0
	return roundToNearest10(scaled)
}

func Score(rawCorrectRW, rawCorrectMath int) (int, int, int) {
	if rawCorrectRW < 0 {
		rawCorrectRW = 0
	}
	if rawCorrectMath < 0 {
		rawCorrectMath = 0
	}
	if rawCorrectRW > RWTotal {
		rawCorrectRW = RWTotal
	}
	if rawCorrectMath > MathTotal {
		rawCorrectMath = MathTotal
	}

	scaledRW := scaleToSection(rawCorrectRW, RWTotal)
	scaledMath := scaleToSection(rawCorrectMath, MathTotal)
	total := scaledRW + scaledMath
	return scaledRW, scaledMath, total
}
