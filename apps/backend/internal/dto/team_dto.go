package dto

type (
	CreateTeamDto struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
	}
)
