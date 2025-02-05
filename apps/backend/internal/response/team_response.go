package response

type (
	CreateTeamResponse struct {
		CommonModel
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
	}
	TeamResponse struct {
		CommonModel
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
	}
)
