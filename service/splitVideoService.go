package service

import (
	cv "gocv-example/computerVision"
	"gocv-example/dto"
	"log"
	"time"
)

type SplitVideoService interface {
	SplitVideo(req dto.SplitVideoRequest) (*dto.SplitVideoResponse, error)
}

type DefaultSplitVideoService struct {
}

func (s DefaultSplitVideoService) SplitVideo(req dto.SplitVideoRequest) (*dto.SplitVideoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	video := cv.NewVideo()

	startTime := time.Now()
	video.Split(req)
	elapsedTime := time.Since(startTime)
	log.Printf("실행시간: %s\n", elapsedTime)

	return &dto.SplitVideoResponse{VideoURL: req.VideoPath}, nil
}

func NewSplitVideoService() DefaultSplitVideoService {
	return DefaultSplitVideoService{}
}
