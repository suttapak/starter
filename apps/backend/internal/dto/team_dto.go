package dto

type (
	CreateTeamDto struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}
	UpdateMemberRoleDto struct {
		UserId uint `json:"user_id"`
		RoleId uint `json:"role_id"`
	}
	CreateTeamPendingTeamMemberDto struct {
		Username string `json:"username"`
	}
	AcceptTeamMemberDto struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}
)
