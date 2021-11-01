package app

import (
	"gocv-example/service"

	"github.com/labstack/echo"
)

func Routes() *echo.Echo {
	e := echo.New()
	sh := SplitVideoHandler{Service: service.NewSplitVideoService()}

	e.POST("/api/v1/video/split", sh.SplitVideo)

	return e
}
