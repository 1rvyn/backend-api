package models

import (
	"encoding/json"
)

type Question struct {
	ID                int             `json:"id" gorm:"primaryKey"`
	Problem           string          `json:"problem"`
	ExampleAnswer     string          `json:"example_answer"`
	ExampleInput      string          `json:"example_input"`
	ProblemType       string          `json:"problem_type"`
	ProblemDifficulty string          `json:"problem_difficulty"`
	TemplateCode      json.RawMessage `json:"template_code"`
}
