package repo

import (
	"context"
	"mirey7/project-project/internal/data/task"
)

type TaskRepo interface {
	FindTaskByIds(ctx context.Context, ids []int64) ([]*task.MsTaskStagesTemplate, error)
}
