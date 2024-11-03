package dao

import (
	"context"
	"mirey7/project-user/internal/data/organization"
	"mirey7/project-user/internal/database/gorms"
)

type Organization struct {
	conn *gorms.GormConn
}

func NewOrganization() *Organization {
	return &Organization{
		conn: gorms.New(),
	}
}

func (o *Organization) SaveOrganization(ctx context.Context, org *organization.Organization) error {
	return o.conn.Session(ctx).Create(org).Error
}
