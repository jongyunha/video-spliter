package app

import (
	"fmt"
	"gocv-example/dto"
	"gocv-example/service"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
)

type SplitVideoHandler struct {
	Service service.SplitVideoService
}

func (h SplitVideoHandler) SplitVideo(c echo.Context) error {
	rows, _ := strconv.Atoi(c.FormValue("rows"))
	columns, _ := strconv.Atoi(c.FormValue("columns"))
	file, err := c.FormFile("video_file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Destination
	uuid4, _ := uuid.NewV4()
	extension := filepath.Ext(file.Filename)
	videoFileName := fmt.Sprintf("%s.%s", uuid4, extension)
	dst, err := os.Create(videoFileName)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer dst.Close()
	defer os.Remove(videoFileName)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	dir, _ := os.Getwd()

	req := dto.SplitVideoRequest{
		Rows:      rows,
		Columns:   columns,
		VideoPath: fmt.Sprintf("%s/%s", dir, videoFileName),
		VideoName: videoFileName,
	}

	res, err := h.Service.SplitVideo(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
