package repo

import (
	"context"
	"mirey7/project-user/internal/data/organization"
	"mirey7/project-user/internal/database"
)

type OrganizationRepo interface {
	SaveOrganization(conn database.DBConn, ctx context.Context, org *organization.Organization) error
	FindOrganizationByMenId(ctx context.Context, memId int64) ([]*organization.Organization, error)
}
