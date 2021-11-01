package cv

import (
	"log"

	"gocv.io/x/gocv"
)

type Video interface {
	Read(videoPath string)
}

type DefaultVideo struct {
}

func (v DefaultVideo) Create(videoPath string) {
	cap, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		log.Println("Cannot read videofile")
	}
	fps := cap.Get(gocv.VideoCaptureFPS)
	width := cap.Get(gocv.VideoCaptureFrameWidth)
	height := cap.Get(gocv.VideoCaptureFrameHeight)
	vw, err := gocv.VideoWriterFile("new_video.avi", "x264", fps, int(width), int(height), true)

	if err != nil {
		log.Println("Cannot create video writer")
	}
	defer vw.Close()

	img := gocv.NewMat()
	defer img.Close()
	for {
		cap.Read(&img)
		if img.Empty() {
			log.Println("Can not get frame")
			break
		}
		err = vw.Write(img)
		if err != nil {
			log.Println("Can not write frame")
			break
		}
		gocv.WaitKey(1)
		// croppedMat := img.Region(image.Rect(0, 0, 400, 400))
	}
}

func NewVideo() DefaultVideo {
	return DefaultVideo{}
}
