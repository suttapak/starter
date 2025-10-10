package model

type (
	Team struct {
		CommonModel
		Name        string       `db:"name" json:"name"`
		Address     *string      `db:"address" json:"address"`
		Phone       *string      `db:"phone" json:"phone"`
		Email       *string      `db:"email" json:"email"`
		Username    string       `db:"username" json:"username"`
		Description *string      `db:"description" json:"description"`
		TeamMembers []TeamMember `db:"-" json:"team_members,omitempty"`
		Products    []Product    `db:"-" json:"products,omitempty"`
	}

	TeamMember struct {
		CommonModel
		TeamID     uint      `db:"team_id" json:"team_id"`
		UserID     uint      `db:"user_id" json:"user_id"`
		TeamRoleID uint      `db:"team_role_id" json:"team_role_id"`
		IsActive   bool      `db:"is_active" json:"is_active"`
		Team       *Team     `db:"-" json:"team,omitempty"`
		User       *User     `db:"-" json:"user,omitempty"`
		TeamRole   *TeamRole `db:"-" json:"team_role,omitempty"`
	}

	TeamRole struct {
		CommonModel
		Name        string       `db:"name" json:"name"`
		TeamMembers []TeamMember `db:"-" json:"team_members,omitempty"`
	}
)
