package grading

import (
	"math"

	"github.com/myselfBZ/sat-jade/internal/store"
)


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


// this function HAS to be guaranteed to recieve an intiger between 0 and 97
func getModuleNameFromIdx(responseIdx int) string {
	if responseIdx < 27 {
		return  "Reading And Writing 1"
	}

	if responseIdx < 54 && responseIdx > 26 {
		return "Reading And Writing 2"
	}


	if responseIdx < 76 && responseIdx > 53 {
		return "Math 1"
	}

	if responseIdx < 98 && responseIdx > 75 {
		return "Math 1"
	}

	return ""
}


func Check(studentResponse []string, correctAnswers []string) *store.Result {
	var result store.Result
	var answers []*store.ResultAnswer

	rwCorrect := 0
	mathCorrect := 0
	for i, correct := range correctAnswers {
		answer := store.ResultAnswer{
			UserAnswer:    studentResponse[i],
			CorrectAnswer: correct,
			Module:        getModuleNameFromIdx(i),
		}

		if correct == studentResponse[i] {
			answer.Status = "correct"
			// counting correct answers by domain on the fly
			if i < 54 {
				rwCorrect++
			} else {
				mathCorrect++
			}
		} else if studentResponse[i]  == "" {
			answer.Status = "omitted"
		} else {
			answer.Status = "incorrect"
		}

		answers = append(answers, &answer)
	}

	rwScore, mathScore, totalScore := Score(rwCorrect, mathCorrect)
	result.MathScore = int32(mathScore)
	result.EnglishScore = int32(rwScore)
	result.TotalScore = int32(totalScore)
	result.Answers = answers
	return &result
}
