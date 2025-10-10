package model

type (
	Image struct {
		CommonModel
		Path   string  `db:"path" json:"path"`
		Url    string  `db:"url" json:"url"`
		Size   float64 `db:"size" json:"size"`
		Width  uint    `db:"width" json:"width"`
		Height uint    `db:"height" json:"height"`
		Type   string  `db:"type" json:"type"`
		UserID uint    `db:"user_id" json:"user_id"`
		User   *User   `db:"-" json:"user,omitempty"`
	}
)
