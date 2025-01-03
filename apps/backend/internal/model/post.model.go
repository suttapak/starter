package model

type (
	Post struct {
		CommonModel
		UserID uint   `json:"user_id"`
		User   User   `json:"user"`
		Data   string `json:"data"`
	}
)
