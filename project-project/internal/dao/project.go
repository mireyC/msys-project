package dao

import (
	"context"
	"fmt"
	"mirey7/project-project/internal/data/pro"
	"mirey7/project-project/internal/database"
	"mirey7/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn database.DBConn
}

func (p *ProjectDao) FindCollectByProjectCodeAndMemberId(ctx context.Context, projectCode int64, memberId int32) (bool, error) {
	session := p.conn.Session(ctx)
	sql := fmt.Sprintf("select count(*) from ms_project_collection a where a.project_code=? and a.member_code=?")
	db := session.Raw(sql, projectCode, memberId)
	var count int
	err := db.Scan(&count).Error

	return count > 0, err
}

func (p *ProjectDao) FindProjectAndMember(ctx context.Context, projectCode int64, memberId int32) (*pro.ProjectAndMember, error) {
	session := p.conn.Session(ctx)

	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id=b.project_code and b.id=?")
	db := session.Raw(sql, projectCode)
	var projectAndMember *pro.ProjectAndMember
	err := db.Scan(&projectAndMember).Error
	return projectAndMember, err
}

func (p *ProjectDao) FindWonerProject(ctx context.Context, ownerId int64) (*pro.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectDao) SaveProject(ctx context.Context, conn database.DBConn, pr *pro.Project) error {
	conn = conn.(*gorms.GormConn)
	return conn.Tx(ctx).Save(&pr).Error
}

func (p *ProjectDao) SaveProjectMember(ctx context.Context, conn database.DBConn, pm *pro.ProjectMember) error {
	conn = conn.(*gorms.GormConn)
	return conn.Tx(ctx).Save(&pm).Error
}

func (p ProjectDao) FindProjectTemplateSystem(ctx context.Context, page int64, pageSize int64, system int) ([]*pro.ProjectTemplate, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	if index < 0 {
		index = 0
	}
	sql := fmt.Sprintf("select * from ms_project_template where is_system=1 order by sort limit ?, ?")
	db := session.Raw(sql, index, pageSize)
	var pts []*pro.ProjectTemplate
	err := db.Scan(&pts).Error

	sql1 := fmt.Sprintf("select count(*) from ms_project_template where is_system=1")
	db1 := session.Raw(sql1)
	var total int64
	err = db1.Scan(&total).Error
	return pts, total, err
}

func (p ProjectDao) FindProjectTemplateCustom(ctx context.Context, page int64, pageSize int64, organizationId int64, memberId int64) ([]*pro.ProjectTemplate, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	if index < 0 {
		index = 0
	}
	sql := fmt.Sprintf("select * from ms_project_template where organization_code=? and member_code=? and is_system=0 order by sort limit ?, ?")
	db := session.Raw(sql, organizationId, memberId, index, pageSize)
	var pts []*pro.ProjectTemplate
	err := db.Scan(&pts).Error

	sql1 := fmt.Sprintf("select count(*) from ms_project_template where organization_code=? and member_code=? and is_system=0")
	var total int64
	db1 := session.Raw(sql1, organizationId, memberId)
	err = db1.Scan(&total).Error
	return pts, total, err
}

func (p ProjectDao) FindProjectTemplateAll(ctx context.Context, page int64, pageSize int64, organizationId int64) ([]*pro.ProjectTemplate, int64, error) {
	session := p.conn.Session(ctx)

	index := (page - 1) * pageSize
	if index < 0 {
		index = 0
	}
	sql := fmt.Sprintf("select * from ms_project_template where organization_code=? order by sort limit ?, ?")
	db := session.Raw(sql, organizationId, index, pageSize)
	var pts []*pro.ProjectTemplate
	err := db.Scan(&pts).Error

	sql1 := fmt.Sprintf("select count(*) from ms_project_template where organization_code=?")
	var total int64
	db1 := session.Raw(sql1, organizationId)
	err = db1.Scan(&total).Error
	return pts, total, err
}

func (p ProjectDao) FindCollectProjectByMenId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	if index < 0 {
		index = 0
	}
	sql := fmt.Sprintf("select * from ms_project where id in (select project_code from ms_project_collection where member_code=? ) order by sort limit ?, ?")
	db := session.Raw(sql, memId, index, pageSize)
	var pm []*pro.ProjectAndMember
	err := db.Scan(&pm).Error
	var total int64
	query := fmt.Sprintf("member_code=?")
	session.Model(&pro.CollectionProject{}).Where(query, memId).Count(&total)
	return pm, total, err
}

func (p ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64, condition string) ([]*pro.ProjectAndMember, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	if index < 0 {
		index = 0
	}
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id=b.project_code and b.member_code=? %s order by sort limit ?,?", condition)
	db := session.Raw(sql, memId, index, pageSize)
	var mp []*pro.ProjectAndMember
	err := db.Scan(&mp).Error
	var total int64
	sql2 := fmt.Sprintf("select count(*) from ms_project a, ms_project_member b where a.id=b.project_code and b.member_code=? %s", condition)
	db2 := session.Raw(sql2, memId)
	err = db2.Scan(&total).Error
	return mp, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.NewTran(),
	}
}
