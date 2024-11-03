package repo

import (
	"context"
	"mirey7/project-user/internal/data/organization"
)

type OrganizationRepo interface {
	SaveOrganization(ctx context.Context, org *organization.Organization) error
}
