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


func getChoiceIDByLabel(choices []store.AnswerChoice, label string) *int32 {
	for i := 0; i < len(choices); i++ {
		if choices[i].Label == label {
			return &choices[i].ID
		}
	}
	return nil
}

func Check(studentResponse []string, correctAnswers []store.CorrectAnswerWithAnswerChoices) *store.Result {
	var result store.Result
	var answers []store.ResultAnswer

	rwCorrect := 0
	mathCorrect := 0
	for i := 0; i < len(studentResponse); i++ {
		answer := store.ResultAnswer{
			UserAnswerId: getChoiceIDByLabel(correctAnswers[i].AnswerChoices, studentResponse[i]),
			QuestionId: correctAnswers[i].QuestionID,
		}

		if correctAnswers[i].CorrectAnswer == studentResponse[i] {
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

		answers = append(answers, answer)
	}

	rwScore, mathScore, totalScore := Score(rwCorrect, mathCorrect)
	result.MathScore = int32(mathScore)
	result.EnglishScore = int32(rwScore)
	result.TotalScore = int32(totalScore)
	result.Answers = answers
	return &result
}
