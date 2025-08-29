package repository

import "gorm.io/gorm"

type (
	DatabaseTransaction interface {
		BeginTx() *gorm.DB
		CommitTx(tx *gorm.DB) error
		RollbackTx(tx *gorm.DB) error
	}

	databaseTransaction struct {
		db *gorm.DB
	}
)

// BeginTx implements DatabaseTransaction.
func (d *databaseTransaction) BeginTx() *gorm.DB {
	return d.db.Begin()
}

// CommitTx implements DatabaseTransaction.
func (d *databaseTransaction) CommitTx(tx *gorm.DB) error {
	if tx != nil {
		return tx.Commit().Error
	}
	return nil
}

// RollbackTx implements DatabaseTransaction.
func (d *databaseTransaction) RollbackTx(tx *gorm.DB) error {
	if tx != nil {
		return tx.Rollback().Error
	}
	return nil
}

func NewDatabaseTransaction(db *gorm.DB) DatabaseTransaction {
	return &databaseTransaction{
		db: db,
	}
}
