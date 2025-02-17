package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"context"
	"time"
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

	// 3초 후 타임아웃 되는 context 생성
	ctx3, cancel3 := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel3()

	// g++ 컴파일 명령 실행
	fmt.Println("Compiling the C++ code...")
	cmd := exec.CommandContext(ctx3, "g++", tmpFile.Name(), "-o", outputFile)
	// CombinedOutput을 사용하여 표준 출력과 표준 에러를 모두 캡처
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(compileOutput), err
	}

	// 1초 후 타임아웃 되는 context 생성
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel1()

	// 컴파일된 바이너리 실행
	fmt.Println("Running the compiled binary...")
	runCmd := exec.CommandContext(ctx1, "prlimit",
	"--as=134217728",  // 메모리 128MB
	"--fsize=0",  // 만들 수 있는 파일 크기 0MB
	"--nofile=4",  // 사용하는 파일 4개
	"--nproc=0",  // 자식 프로세스 0개
	"--", outputFile)

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
