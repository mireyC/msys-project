package tran

import "mirey7/project-project/internal/database"

// Transaction 事务操作，需要注入数据库连接 gorm.db
type Transaction interface {
	Action(func(conn database.DBConn) error) error
}
