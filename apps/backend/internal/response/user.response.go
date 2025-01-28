package response

type (
	// swagger:response UserResponse
	UserResponse struct {
		CommonModel
		Username      string `json:"username"` // username is uniqe of user in system
		Email         string `json:"email"`    // email feild for register and forgot password processed
		EmailVerifyed bool   `json:"email_verifyed"`
		FullName      string `json:"full_name"`
		RoleID        uint   `json:"role_id"`
	}
)
