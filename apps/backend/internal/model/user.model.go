package model

type (
	User struct {
		CommonModel
		Username      string         `json:"username"` // username is uniqe of user in system
		Password      string         `json:"-"`        // passowrd is user login password
		Email         string         `json:"email"`    // email feild for register and forgot password processed
		EmailVerifyed bool           `json:"email_verifyed"`
		FirstName     string         `json:"first_name"`
		LastName      string         `json:"last_name"`
		ProfileImage  []ProfileImage `json:"profile_image"`
		Role          Role           `json:"role"`
		RoleID        uint           `json:"role_id"`
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
