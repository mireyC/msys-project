package database

import (
	"context"
	"gorm.io/gorm"
)

type DBConn interface {
	Begin()
	Rollback()
	Commit()
	Session(ctx context.Context) *gorm.DB
	Tx(ctx context.Context) *gorm.DB
}
