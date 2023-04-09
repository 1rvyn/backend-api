package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func wrapGoCodeWithTests(userID, code, filext, testCases string) string {
	testCode := `
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

` + code + `

func main() {
	testCases := ` + testCases + `

	results := []bool{}
	for _, tc := range testCases {
		result := twoSum(tc.inputNums, tc.inputTarget)
		if reflect.DeepEqual(result, tc.expected) {
			results = append(results, true)
		} else {
			results = append(results, false)
		}
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error encoding results to JSON")
		return
	}
	fmt.Println(string(resultsJSON))
}
`
	return testCode
}

func getGoTestCases(questionID string) string {
	// two_sum test cases
	if questionID == "1" {
		return `[]struct {
		inputNums   []int
		inputTarget int
		expected    []int
	}{
		{
			inputNums:   []int{2, 7, 11, 15},
			inputTarget: 9,
			expected:    []int{0, 1},
		},
	
		{	
			inputNums:   []int{3, 2, 4},
			inputTarget: 6,
			expected:    []int{1, 2},
		},

		{
			inputNums:   []int{3, 3},
			inputTarget: 6,
			expected:    []int{0, 1},
		},
	}`
		// Reverseinteger test cases
	} else if questionID == "2" {
		return `[]struct {
		inputNum int
		expected int
	}{
		{
			inputNum: 123,
			expected: 321,
		},
	
		{	
			inputNum: -123,
			expected: -321,
		},

		{
			inputNum: 120,
			expected: 21,
		},

		{
			inputNum: 0,
			expected: 0,
		},
	}`
	} else if questionID == "3" {
		return `[]struct {
		inputNum int
		expected bool
	}{
		{
			inputNum: 121,
			expected: true,
		},
	
		{	
			inputNum: -121,
			expected: false,
		},

		{
			inputNum: 10,
			expected: false,
		},

		{
			inputNum: -101,
			expected: false,
		},
	}`
		// RemoveDuplicates test cases
	} else if questionID == "4" {
		return `[]struct {
		inputNums []int
		expected  int
	}{
		{
			inputNums: []int{1, 1, 2},
			expected:  2,
		},
	
		{	
			inputNums: []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4},
			expected:  5,
		},
		
		{
			inputNums: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 3},
			expected:  3,
		},
		
		{
			inputNums: []int{5, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected:  5,
		},
	}`

	} else {
		return ""
	}
}

func runGoCode(userID, code, questionID, filext string) (string, string) {

	testCases := getGoTestCases(questionID)
	testCode := wrapGoCodeWithTests(userID, code, filext, testCases)

	file, err := os.Create("./remotecode/" + userID + "." + filext)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString(testCode)

	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "run", "./remotecode/"+userID+"."+filext)

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	output := outBuf.String()
	errorOutput := errBuf.String()

	return output, errorOutput
}

func runPythonCode(userID, code, filext string) (string, string) {

	file, err := os.Create("./remotecode/" + userID + "." + filext)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString(code)

	if err != nil {
		panic(err)
	}

	cmd := exec.Command("python", "./remotecode/"+userID+"."+filext)

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	output := outBuf.String()
	errorOutput := errBuf.String()

	return output, errorOutput
}

func runJSCode(userID, code, filext string) (string, string) {

	file, err := os.Create("./remotecode/" + userID + "." + filext)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString(code)

	if err != nil {
		panic(err)
	}
	cmd := exec.Command("node", "./remotecode/"+userID+"."+filext)

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	output := outBuf.String()
	errorOutput := errBuf.String()

	return output, errorOutput
}

func Marking(code, questionID, userID, language string) string {
	// This is the function I will use to mark the code snippet
	fmt.Println("the code is: ", code)
	fmt.Println("the question ID is: ", questionID)
	//file, err := os.Create("./remotecode/" + userID + "." + language)

	var output, errorOutput string

	switch language {
	case "python":
		output, errorOutput = runPythonCode(userID, code, "py")
	case "go":
		output, errorOutput = runGoCode(userID, code, questionID, "go")
	case "javascript":
		output, errorOutput = runJSCode(userID, code, "js")
		// Add other cases for each supported language (e.g., "javascript", "golang", ...)
	default:
		errorOutput = "Unsupported language"
	}

	if errorOutput != "" {
		fmt.Println("Error: ", errorOutput)
		return errorOutput
	} else {
		var results []bool
		err := json.Unmarshal([]byte(output), &results)
		if err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return "Error decoding JSON"
		}
		resultsJSON, _ := json.Marshal(results)
		return string(resultsJSON)
	}
}
