package repository

import "gorm.io/gorm"

type (
	CommonDB interface {
		BeginTx() *gorm.DB
		CommitTx(tx *gorm.DB)
		RollbackTx(tx *gorm.DB)
	}
)
