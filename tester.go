package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Running(cppCode string, input string, output string) (string, error) {
	// 임시 파일 생성
	tmpFile, err := os.CreateTemp("", "example-*.cpp")
	if err != nil {
		return "Failed to create temporary file", err
	}
	defer os.Remove(tmpFile.Name()) // 사용 후 파일 삭제

	// C++ 코드 쓰기
	if _, err := tmpFile.Write([]byte(cppCode)); err != nil {
		return "Failed to write to temporary file", err
	}
	tmpFile.Close()

	// 출력 바이너리 파일 경로
	outputFile := tmpFile.Name() + ".out"

	defer os.Remove(outputFile) // 사용 후 파일 삭제

	// g++ 컴파일 명령 실행
	fmt.Println("Compiling the C++ code...")
	cmd := exec.Command("g++", tmpFile.Name(), "-o", outputFile)
	// CombinedOutput을 사용하여 표준 출력과 표준 에러를 모두 캡처
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(compileOutput), err
	}

	// 컴파일된 바이너리 실행
	runCmd := exec.Command(outputFile)
	// input을 표준 입력으로 전달
	runCmd.Stdin = strings.NewReader(input)

	runOutput, runErr := runCmd.CombinedOutput()
	if runErr != nil {
		return string(runOutput), runErr
	}

	fmt.Println("Output: ", string(runOutput))

	if strings.TrimSpace(string(runOutput)) != strings.TrimSpace(output) {
		return fmt.Sprintf("<< Expected output >>\n%s\n<< Your output >>\n%s", output, runOutput), fmt.Errorf("Output mismatch")
	}

	return string(runOutput), nil
}
