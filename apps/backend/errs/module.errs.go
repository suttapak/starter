package errs

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrNotFound                   = New(http.StatusNotFound, "ไม่พบข้อมูล")
	ErrInternal                   = New(http.StatusInternalServerError, "บางอย่างผิดพลาด")
	ErrInvalid                    = New(http.StatusUnprocessableEntity, "ข้อมูลไม่ถูกต้อง")
	ErrConflict                   = New(http.StatusConflict, "ข้อมูลซ้ำกัน")
	ErrForbidden                  = New(http.StatusForbidden, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrUnauthorized               = New(http.StatusUnauthorized, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrBadRequest                 = New(http.StatusBadRequest, "ข้อมูลไม่ถูกต้อง")
	ErrDuplicateUsername          = New(http.StatusBadRequest, "username นี้มีคนใช้งานแล้ว กรุณาใช้ username ใหม่")
	ErrDuplicateEmail             = New(http.StatusBadRequest, "email นี้มีคนใช้งานแล้ว กรุณาใช้ email ใหม่")
	ErrHashPassword               = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการเข้ารหัส password")
	ErrRegisterUsername           = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการบันทึกข้อมูล")
	ErrGenerateJWTFail            = New(http.StatusInternalServerError, "บางอย่างผิดพลาดในการออก token")
	ErrUsernameOrPasswordIncorect = New(http.StatusUnauthorized, "username หรือ password ไม่ถูกต้อง")
	ErrVerifyEmail                = New(http.StatusBadRequest, "ข้อมูลการยืนยันตัวตนไม่ถูกต้อง")
	ErrTeamUsernameIsUsed         = New(http.StatusBadRequest, "username นี้มีการใช้งานแล้วกรุณาใช้ username อื่น")
	ErrDuplicatedKey              = New(http.StatusBadRequest, "มีข้อมูลนี้แล้ว ข้อมูลซ่ำโปรดลองกรอกข้อมูลอื่นๆ")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func HandeGorm(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicatedKey
	default:
		return ErrInternal
	}
}

func New(code int, message string) error {
	return AppError{
		Code:    code,
		Message: message,
	}
}
