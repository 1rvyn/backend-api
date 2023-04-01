package models

type TestResult struct {
	Output   string `json:"output"`
	Success  bool   `json:"success"`
	TestName string `json:"test_name"`
}
