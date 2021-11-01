package service

import (
	cv "gocv-example/computerVision"
	"gocv-example/dto"
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
	video.Create(req.VideoPath)

	return &dto.SplitVideoResponse{VideoURL: req.VideoPath}, nil
}

func NewSplitVideoService() DefaultSplitVideoService {
	return DefaultSplitVideoService{}
}
