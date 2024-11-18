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

func (p ProjectDao) FindProjectTemplateSystem(ctx context.Context, page int64, pageSize int64, system int) ([]*pro.ProjectTemplate, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
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
	sql := fmt.Sprintf("select * from ms_project_template where organization_code=? and is_system=0 order by sort limit ?, ?")
	db := session.Raw(sql, organizationId, index, pageSize)
	var pts []*pro.ProjectTemplate
	err := db.Scan(&pts).Error

	sql1 := fmt.Sprintf("select count(*) from ms_proje_template where organization_code=? and is_system=0")
	var total int64
	db1 := session.Raw(sql1, organizationId)
	err = db1.Scan(&total).Error
	return pts, total, err
}

func (p ProjectDao) FindCollectProjectByMenId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
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
		conn: gorms.New(),
	}
}
