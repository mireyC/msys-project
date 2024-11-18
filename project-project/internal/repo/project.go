package repo

import (
	"context"
	"mirey7/project-project/internal/data/pro"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64, condition string) ([]*pro.ProjectAndMember, int64, error)
	FindCollectProjectByMenId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error)
	FindProjectTemplateSystem(ctx context.Context, page int64, pageSize int64, system int) ([]*pro.ProjectTemplate, int64, error)
	FindProjectTemplateCustom(ctx context.Context, page int64, pageSize int64, organizationId int64, memberId int64) ([]*pro.ProjectTemplate, int64, error)
	FindProjectTemplateAll(ctx context.Context, page int64, pageSize int64, organizationId int64) ([]*pro.ProjectTemplate, int64, error)
}
