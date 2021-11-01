package service

import (
	cv "gocv-example/computerVision"
	"gocv-example/dto"
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
	res, err := video.Split(req)
	if err != nil {
		return nil, err
	}
	elapsedTime := time.Since(startTime)
	res.ElapsedTime = elapsedTime

	return res, nil
}

func NewSplitVideoService() DefaultSplitVideoService {
	return DefaultSplitVideoService{}
}
