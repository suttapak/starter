package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrNotFound                     = New(http.StatusNotFound, "ไม่พบข้อมูล")
	ErrInternal                     = New(http.StatusInternalServerError, "บางอย่างผิดพลาด")
	ErrInvalid                      = New(http.StatusUnprocessableEntity, "ข้อมูลไม่ถูกต้อง")
	ErrConflict                     = New(http.StatusConflict, "ข้อมูลซ้ำกัน")
	ErrForbidden                    = New(http.StatusForbidden, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrUnauthorized                 = New(http.StatusUnauthorized, "ไม่มีสิทธิ์ในการเข้าถึงข้อมูลนี้")
	ErrBadRequest                   = New(http.StatusBadRequest, "ข้อมูลไม่ถูกต้อง")
	ErrDuplicateUsername            = New(http.StatusBadRequest, "username นี้มีคนใช้งานแล้ว กรุณาใช้ username ใหม่")
	ErrDuplicateEmail               = New(http.StatusBadRequest, "email นี้มีคนใช้งานแล้ว กรุณาใช้ email ใหม่")
	ErrHashPassword                 = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการเข้ารหัส password")
	ErrRegisterUsername             = New(http.StatusBadRequest, "บางอย่างผิดพลาดในการบันทึกข้อมูล")
	ErrGenerateJWTFail              = New(http.StatusInternalServerError, "บางอย่างผิดพลาดในการออก token")
	ErrUsernameOrPasswordIncorrect  = New(http.StatusUnauthorized, "username หรือ password ไม่ถูกต้อง")
	ErrVerifyEmail                  = New(http.StatusBadRequest, "ข้อมูลการยืนยันตัวตนไม่ถูกต้อง")
	ErrTeamUsernameIsUsed           = New(http.StatusBadRequest, "username นี้มีการใช้งานแล้วกรุณาใช้ username อื่น")
	ErrDuplicatedKey                = New(http.StatusBadRequest, "มีข้อมูลนี้แล้ว ข้อมูลซ่ำโปรดลองกรอกข้อมูลอื่นๆ")
	ErrNotActiveTeamId              = New(http.StatusUnauthorized, "กรุณาระบุแผนกที่ทำงานอยู่")
	ErrEmailNotVerify               = New(http.StatusUnauthorized, "ผู้ใช้งานยังไม่ได้ยืนยันอีเมล์")
	ErrSendEmail                    = New(http.StatusInternalServerError, "ไม่สามารถส่งอีเมล์ได้")
	ErrUserAlreadyInTeam            = New(http.StatusBadRequest, "ผู้ใช้งานอยู่ในแผนกแล้ว")
	ErrForeignKeyViolation          = New(http.StatusBadRequest, "ไม่สามารถลบข้อมูลได้เนื่องจากมีข้อมูลที่อ้างอิงถึงข้อมูลนี้")
	ErrUniqueViolation              = New(http.StatusBadRequest, "ไม่สามารถเพิ่มข้อมูลได้เนื่องจากมีข้อมูลนี้แล้ว")
	ErrNotNullViolation             = New(http.StatusBadRequest, "ไม่สามารถเพิ่มข้อมูลได้เนื่องจากข้อมูลไม่ครบถ้วน")
	ErrIntegrityConstraintViolation = New(http.StatusBadRequest, "ไม่สามารถเพิ่มข้อมูลได้เนื่องจากข้อมูลไม่ถูกต้อง")
	ErrProductLotOutOfStock         = New(http.StatusBadRequest, "สินค้าหมดสต็อก กรุณาติดต่อผู้ดูแลระบบ")
	ErrTransactionNotPending        = New(http.StatusBadRequest, "ไม่สามารถดำเนินการได้เนื่องจากสถานะของธุรกรรมไม่ใช่ 'pending' หรือ 'draft' กรุณาตรวจสอบสถานะธุรกรรมอีกครั้ง")
	ErrProductLotNotFound           = New(http.StatusBadRequest, "ไม่พบข้อมูล Lot สินค้า กรุณาตรวจสอบอีกครั้ง")
	ErrDisabledFunction             = New(http.StatusBadRequest, "ฟังก์ชันนี้ถูกปิดใช้งาน")
	ErrReturnMoreThanLotItem        = New(http.StatusBadRequest, "ไม่สามารถระบุจำนวนคืนเกินจำนวนที่เบิก")
	ErrDoNotHaveLotItem             = New(http.StatusBadRequest, "ไม่พบข้อมูล Lot ที่ต้องการคืน")
	ErrTransactionNotComplete       = New(http.StatusBadRequest, "รายการที่อ้างอิงยังไม่ได้อนุมัติ")
	ErrTransactionNotTypeSale       = New(http.StatusBadRequest, "รายการอ้างอิงไม่ใช่ประเภทขาย/เบิก")
	ErrChildTransactionNotApprove   = New(http.StatusBadRequest, "มีรายการคืนที่อ้างอิงจากรายการนี้บางรายการยังไม่ได้อนุมัติ")
	ErrFileUploadNotImage           = New(http.StatusBadRequest, "ไฟล์ที่อัปโหลดไม่ใช่รูปภาพที่รองรับ")
	ErrFileImageCanNotGetStats      = New(http.StatusBadRequest, "ไม่สามารถดึงข้อมูลของรูปภาพได้")
	ErrFileImageCanNotSaveToDisk    = New(http.StatusBadRequest, "ไม่สามารถบันทึกไฟล์รูปภาพลงดิสก์ได้")
	ErrFileUploadNotFound           = New(http.StatusBadRequest, "ไม่พบไฟล์ที่อัปโหลด")
	ErrFileUploadNoFile             = New(http.StatusBadRequest, "ไม่พบไฟล์ที่อัปโหลด")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func HandleGorm(err error, defaultError ...AppError) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDuplicatedKey
	default:
		if errPg, ok := err.(*pgconn.PgError); ok {
			switch errPg.Code {
			case "23505":
				return ErrUniqueViolation
			case "23503":
				return ErrForeignKeyViolation
			case "23502":
				return ErrNotNullViolation
			}
		}

		if len(defaultError) > 0 {
			return defaultError[0]
		}
		return ErrInternal
	}
}

func New(code int, message string) error {
	return AppError{
		Code:    code,
		Message: message,
	}
}
