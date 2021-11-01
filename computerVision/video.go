package cv

import (
	"fmt"
	"gocv-example/dto"
	"image"
	"log"
	"os"
	"sync"

	"github.com/gofrs/uuid"
	"gocv.io/x/gocv"
)

var wg sync.WaitGroup

type Video interface {
}

type DefaultVideo struct {
}

type ImageCV struct {
	mat gocv.Mat
}

type CropRect struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

type SplitDefine struct {
	n         int
	vws       []gocv.VideoWriter
	cropRects []CropRect
}

func (icv *ImageCV) Crop(cropRect CropRect) *ImageCV {
	croppedMat := icv.mat.Region(image.Rect(cropRect.x0, cropRect.y0, cropRect.x1, cropRect.y1))
	resultMat := croppedMat.Clone()
	return &ImageCV{mat: resultMat}
}

func getWidths(frameWidth, rows int) []int {
	widths := make([]int, 0)
	sWidth := frameWidth / rows
	for i := 1; i <= rows; i++ {
		widths = append(widths, sWidth*i)
	}
	return widths
}

func getHeights(frameHeight, columns int) []int {
	heights := make([]int, 0)
	sHeight := frameHeight / columns
	for i := 1; i <= columns; i++ {
		heights = append(heights, sHeight*i)
	}
	return heights
}

func createVideoWriters(videoNames []string, fps float64, w, h int) ([]gocv.VideoWriter, error) {
	vws := make([]gocv.VideoWriter, 0)
	for _, name := range videoNames {
		vw, err := gocv.VideoWriterFile(name, "H264", fps, w, h, true)
		if err != nil {
			return nil, err
		}
		vws = append(vws, *vw)
	}
	return vws, nil
}

func createVideoNames(rows, columns int) []string {
	videoNames := make([]string, 0)
	dir, _ := os.Getwd()
	for i := 0; i < rows*columns; i++ {
		uuid4, _ := uuid.NewV4()
		videoNames = append(videoNames, fmt.Sprintf("%s/%s.mov", dir, uuid4))
	}
	return videoNames
}

func getCropRects(sWidth, sHeight int, widths, heights []int) []CropRect {
	cropRects := make([]CropRect, 0)
	for _, w := range widths {
		for _, h := range heights {
			cropRect := CropRect{
				x0: w - sWidth,
				y0: h - sHeight,
				x1: w,
				y1: h,
			}
			cropRects = append(cropRects, cropRect)
		}
	}
	return cropRects
}

func (v DefaultVideo) Split(req dto.SplitVideoRequest) error {
	cap, err := gocv.VideoCaptureFile(req.VideoPath)

	if err != nil {
		log.Println("Cannot read videofile")
		return err
	}

	fps := cap.Get(gocv.VideoCaptureFPS)
	frameWidth := int(cap.Get(gocv.VideoCaptureFrameWidth))
	frameHeight := int(cap.Get(gocv.VideoCaptureFrameHeight))
	widths := getWidths(frameWidth, req.Rows)
	heights := getHeights(frameHeight, req.Columns)
	videoNames := createVideoNames(req.Rows, req.Columns)
	videoWriters, err := createVideoWriters(videoNames, fps, frameWidth/req.Rows, frameHeight/req.Columns)

	if err != nil {
		return err
	}
	cropRects := getCropRects(frameWidth/req.Rows, frameHeight/req.Columns, widths, heights)
	splitDefine := SplitDefine{
		n:         req.Rows * req.Columns,
		vws:       videoWriters,
		cropRects: cropRects,
	}

	imcv := ImageCV{
		mat: gocv.NewMat(),
	}
	defer imcv.mat.Close()
	for {
		cap.Read(&imcv.mat)
		if imcv.mat.Empty() {
			break
		}

		wg.Add(len(splitDefine.cropRects))
		go func(c *ImageCV, splitDefine SplitDefine) {
			for i := 0; i < len(splitDefine.cropRects); i++ {
				cropped := imcv.Crop(splitDefine.cropRects[i])
				splitDefine.vws[i].Write(cropped.mat)
				wg.Done()
			}
		}(&imcv, splitDefine)
		wg.Wait()
	}

	for _, vw := range videoWriters {
		vw.Close()
	}

	return nil
}

func NewVideo() DefaultVideo {
	return DefaultVideo{}
}
