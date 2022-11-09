package dto

import "time"

type PictureResponse struct {
	PostedBy string
	Title    string
	Url      string
	Created  time.Time
	Expires  time.Time
}
