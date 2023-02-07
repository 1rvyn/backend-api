package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Marking(code string) string {
	// This is the function I will use to mark the code snippet
	file, err := os.Create("./remotecode/code.py")
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

	cmd := exec.Command("python3", "./remotecode/code.py")

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	output := outBuf.String()
	errorOutput := errBuf.String()

	if errorOutput != "" {
		fmt.Println("Error: ", errorOutput)
		return errorOutput
	} else {
		return output
	}
}
