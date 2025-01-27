package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func decodeEscapedString(input string) string {
	// C++ 코드에서 따옴표 제거
	input = strings.Trim(input, "\"")
	// \n -> 개행 문자, \t -> 탭 문자, \\ -> 역슬래시
	replacer := strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\\\", "\\",
		"\\\"", "\"",
	)
	return replacer.Replace(input)
}

func Running(cppCode string) (string, error) {
	// C++ 코드 디코딩
	cppCode = decodeEscapedString(cppCode)
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	// 컴파일된 바이너리 실행
	fmt.Println("Compilation successful! Executing the binary...")
	runCmd := exec.Command(outputFile)
	runOutput, runErr := runCmd.CombinedOutput()
	if runErr != nil {
		return string(runOutput), runErr
	}

	return string(runOutput), nil
}
