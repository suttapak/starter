package model

type (
	Image struct {
		CommonModel
		Path   string  `json:"path"`
		Url    string  `json:"url"`
		Size   float64 `json:"size"`
		Width  uint    `json:"with"`
		Height uint    `json:"height"`
		Type   string  `json:"type"`
	}
)
