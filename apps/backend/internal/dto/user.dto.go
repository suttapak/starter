package dto

type (
	UserRegisterDto struct {
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
	}

	LoginDto struct {
		UserNameEmail string `json:"username_email" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
	}

	VerifyEmailDto struct {
		Token string `json:"token" form:"token"`
	}
)
