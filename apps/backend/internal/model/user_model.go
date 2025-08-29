package model

type (
	User struct {
		CommonModel
		Username      string         `json:"username" gorm:"uniqueIndex"` // username is unique of user in system
		Password      string         `json:"-"`                           // password is user login password
		Email         string         `json:"email" gorm:"uniqueIndex"`    // email field for register and forgot password processed
		EmailVerifyed bool           `json:"email_verifyed"`
		FullName      string         `json:"full_name"`
		ProfileImage  []ProfileImage `json:"profile_image"`
		Role          Role           `json:"role"`
		RoleID        uint           `json:"role_id"`
		TeamMembers   []TeamMember   `json:"team_members"`
	}

	ProfileImage struct {
		CommonModel
		UserID  uint  `json:"user_id"`
		ImageID uint  `json:"image_id"`
		Image   Image `json:"image"`
	}

	Role struct {
		CommonModel
		Name string `json:"name" gorm:"uniqueIndex:idx_role_name"`
	}
)
