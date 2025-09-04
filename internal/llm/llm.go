package llm

type PracticeOverview struct {
	Overview     string   `json:"overview"`
	Suggesttions []string `json:"suggestions"`
	Motivation   string   `json:"motivation"`
}

type MistakeCountByDomain map[string]int

type PracticeOverviewParams struct {
	CorrectAnswers int
	Mistakes       MistakeCountByDomain
}

type LLM interface {
	GeneratePracticeOverview(params *PracticeOverviewParams) (*PracticeOverview, error)
}

const THE_PRACTICE_PROMPT = `You are an expert level SAT tutor. You have at least 10 years of experience
A student has taken a DSAT practice test. They made %d mistakes and %d correct answers out of all 98 questions.
There are two sections Maths and English. Each have their own domains so dont mix up. 
They made mistakes in the following domains of the test.
%s

use a teen-undertandable language and dont minimize the buzzwords and technical terms to an absolute minimum.

Return a JSON payload ONLY! Nothing else. The JSON payload should look like this.
{
	overview:this is where you give a small overview of how the student did with honesty, and no sugar coating. Only string,
	suggessions:this field is an array where you give breif and small bullet point -like suggestions on how to improve their score. They are just array of strings. Example would be something like this Standard English Conventions (25 errors): suggestion to improve this domain goes here. you see? you show the domain and count the mistakes and give the suggestion 
	motivation:this is where you give student a little motivations according to their score on the practice test with some facts. Only string
}
`
