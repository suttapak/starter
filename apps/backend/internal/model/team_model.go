package model

type (
	Team struct {
		CommonModel
		Name        string       `json:"name"`
		Username    string       `json:"username" gorm:"unique_index"`
		Description string       `json:"description"`
		TeamMembers []TeamMember `json:"team_members"`
	}

	TeamMember struct {
		TeamID     uint     `json:"team_id" gorm:"primaryKey"`
		UserID     uint     `json:"user_id" gorm:"primaryKey"`
		TeamRoleID uint     `json:"team_role_id"`
		Team       Team     `json:"team"`
		User       User     `json:"user"`
		TeamRole   TeamRole `json:"team_role"`
	}

	TeamRole struct {
		CommonModel
		Name        string       `json:"name"`
		TeamMembers []TeamMember `json:"team_members"`
	}
)
