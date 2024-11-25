package repo

import (
	"context"
	"mirey7/project-project/internal/data/pro"
	"mirey7/project-project/internal/database"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64, condition string) ([]*pro.ProjectAndMember, int64, error)
	FindCollectProjectByMenId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error)
	FindProjectTemplateSystem(ctx context.Context, page int64, pageSize int64, system int) ([]*pro.ProjectTemplate, int64, error)
	FindProjectTemplateCustom(ctx context.Context, page int64, pageSize int64, organizationId int64, memberId int64) ([]*pro.ProjectTemplate, int64, error)
	FindProjectTemplateAll(ctx context.Context, page int64, pageSize int64, organizationId int64) ([]*pro.ProjectTemplate, int64, error)
	SaveProject(ctx context.Context, conn database.DBConn, pr *pro.Project) error
	SaveProjectMember(ctx context.Context, conn database.DBConn, pm *pro.ProjectMember) error
	FindProjectAndMember(ctx context.Context, projectCode int64, memberId int32) (*pro.ProjectAndMember, error)
	FindCollectByProjectCodeAndMemberId(ctx context.Context, projectCode int64, memberId int32) (bool, error)
	UpdateDeleteProject(ctx context.Context, projectId int64, isDeleted int32) error
	DelProject(ctx context.Context, projectCode int64) error
	ProjectCollect(ctx context.Context, memberCode int64, projectCode int64, isCollect int32) error
	ProjectEdit(ctx context.Context, project *pro.Project) error
}
