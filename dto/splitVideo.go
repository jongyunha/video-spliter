package dto

import "errors"

type SplitVideoRequest struct {
	Rows      int
	Columns   int
	VideoPath string
	VideoName string
}

type SplitVideoResponse struct {
	VideoURL string `json:"video_url"`
}

func (r SplitVideoRequest) Validate() error {
	if r.Rows <= 0 {
		return errors.New("The value of rows cannot be less than or equal to 0")
	} else if r.Columns <= 0 {
		return errors.New("The value of columns cannot be less than or equal to 0")
	}

	return nil
}
