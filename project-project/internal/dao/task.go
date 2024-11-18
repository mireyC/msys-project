package dao

import (
	"context"
	"fmt"
	"mirey7/project-project/internal/data/task"
	"mirey7/project-project/internal/database"
	"mirey7/project-project/internal/database/gorms"
)

type TaskDao struct {
	conn database.DBConn
}

func (t *TaskDao) FindTaskByIds(ctx context.Context, ids []int64) ([]*task.MsTaskStagesTemplate, error) {
	session := t.conn.Session(ctx)

	sql := fmt.Sprintf("select * from ms_task_stages_template where id in ?")
	db := session.Raw(sql, ids)
	var tasks []*task.MsTaskStagesTemplate
	err := db.Scan(&tasks).Error
	return tasks, err
}

func NewTaskDao() *TaskDao {
	return &TaskDao{conn: gorms.NewTran()}
}
