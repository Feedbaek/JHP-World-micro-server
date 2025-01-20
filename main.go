package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 요청 데이터를 담을 구조체 정의
type CodeRequest struct {
	Code string `json:"code"`
}

func main() {
	// Echo 인스턴스 생성
	e := echo.New()

	// 기본 라우트 설정
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/run", func(c echo.Context) error {
		// C++ 코드 가져오기
		var req CodeRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		res, err := Running(req.Code)

		if err != nil {
			fmt.Println("Error:", err)
			return c.String(http.StatusBadRequest, res)
		}

		fmt.Println(res)

		return c.String(http.StatusOK, res)
	})

	// 서버 실행
	e.Logger.Fatal(e.Start(":8080"))
}
