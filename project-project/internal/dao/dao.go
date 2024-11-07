package dao

import (
	"mirey7/project-project/internal/database"
	"mirey7/project-project/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DBConn
}

func (t *TransactionImpl) Action(f func(conn database.DBConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}

	t.conn.Commit()
	return nil
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTran(),
	}
}
