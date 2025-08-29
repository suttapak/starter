package model

type (
	Image struct {
		CommonModel
		Path   string  `json:"path"`
		Url    string  `json:"url"`
		Size   float64 `json:"size"`
		Width  uint    `json:"width"`
		Height uint    `json:"height"`
		Type   string  `json:"type"`
		UserID uint    `json:"user_id"`
		User   User    `json:"user"`
	}
)
