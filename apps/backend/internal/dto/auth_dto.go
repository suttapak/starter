package dto

type (
	UserRegisterDto struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name" binding:"required"`
	}

	LoginDto struct {
		UserNameEmail string `json:"username" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
	}

	VerifyEmailDto struct {
		Token string `json:"token" form:"token"`
	}
	SendVerifyEmailDto struct {
		UserID uint `json:"user_id"`
	}
)
