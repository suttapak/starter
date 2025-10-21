package model

type (
	User struct {
		CommonModel
		Username      string         `db:"username" json:"username"` // username is unique of user in system
		Password      string         `db:"password" json:"-"`        // password is user login password
		Email         string         `db:"email" json:"email"`       // email field for register and forgot password processed
		EmailVerifyed bool           `db:"email_verifyed" json:"email_verifyed"`
		FullName      string         `db:"full_name" json:"full_name"`
		RoleID        uint           `db:"role_id" json:"role_id"`
		ProfileImage  []ProfileImage `db:"-" json:"profile_image,omitempty"`
		Role          *Role          `db:"-" json:"role,omitempty"`
	}

	ProfileImage struct {
		CommonModel
		UserID  uint   `db:"user_id" json:"user_id"`
		ImageID uint   `db:"image_id" json:"image_id"`
		Image   *Image `db:"-" json:"image,omitempty"`
	}

	Role struct {
		CommonModel
		Name string `db:"name" json:"name"`
	}
)
