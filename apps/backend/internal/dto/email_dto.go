package dto

type (
	VerifyEmailTemplateDataDto struct {
		Email           string
		VerifyEmailLink string
	}
	InviteTeamMemberTemplateDataDto struct {
		TeamName     string
		JoinTeamLink string
	}
)
