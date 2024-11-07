package database

type DBConn interface {
	Begin()
	Rollback()
	Commit()
}
