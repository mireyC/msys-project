package dao

import (
	"context"
	"gorm.io/gorm"
	"mirey7/project-user/internal/data/organization"
	"mirey7/project-user/internal/database"
	"mirey7/project-user/internal/database/gorms"
)

type Organization struct {
	conn *gorms.GormConn
}

func (o *Organization) FindOrganizationByMenId(ctx context.Context, memId int64) ([]*organization.Organization, error) {
	var orgs []*organization.Organization
	err := o.conn.Session(ctx).Where("member_id=?", memId).Order("id ASC").Find(&orgs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return orgs, err
}

func NewOrganization() *Organization {
	return &Organization{
		conn: gorms.New(),
	}
}

func (o *Organization) SaveOrganization(conn database.DBConn, ctx context.Context, org *organization.Organization) error {
	o.conn = conn.(*gorms.GormConn)
	return o.conn.Tx(ctx).Create(org).Error
}
