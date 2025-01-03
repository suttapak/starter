package errs

import "net/http"

var (
	ErrNotFound          = New(http.StatusNotFound, "ไม่พบข้อมูล")
	ErrInternal          = New(http.StatusInternalServerError, "บางอย่างผิดพลาด")
	ErrInvalid           = New(http.StatusUnprocessableEntity, "ข้อมูลไม่ถูกต้อง")
	ErrConflict          = New(http.StatusConflict, "ข้อมูลซ้ำกัน")
	ErrForbidden         = New(http.StatusForbidden, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrUnauthorized      = New(http.StatusUnauthorized, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrBadRequest        = New(http.StatusBadRequest, "ข้อมูลไม่ถูกต้อง")
	ErrDuplicateUsername = New(http.StatusBadRequest, "username นี้มีคนใช้งานแล้ว กรุณาใช้ username ใหม่")
	ErrDuplicateEmail    = New(http.StatusBadRequest, "email นี้มีคนใช้งานแล้ว กรุณาใช้ email ใหม่")
	ErrHashPassword      = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการเข้ารหัส password")
	ErrRegisterUsername  = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการบันทึกข้อมูล")
	ErrGenerateJWTFail   = New(http.StatusInternalServerError, "บางอย่างผิดพลาดในการออก token")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func New(code int, message string) error {
	return AppError{
		Code:    code,
		Message: message,
	}
}
