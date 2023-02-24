package models

type Question struct {
	ID                int    `json:"id" gorm:"primaryKey"`
	Problem           string `json:"problem"`
	ExampleAnswer     string `json:"example_answer"`
	ExampleInput      string `json:"example_input"`
	ExampleOutput     string `json:"example_output"`
	ProblemType       string `json:"problem_type"`
	ProblemDifficulty string `json:"problem_difficulty"`
}
