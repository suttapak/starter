package model

type (
	User struct {
		CommonModel
		Username      string         `json:"username"` // username is uniqe of user in system
		Password      string         `json:"-"`        // passowrd is user login password
		Email         string         `json:"email"`    // email feild for register and forgot password processed
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
		Name string `json:"name"`
	}
)
