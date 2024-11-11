package dao

import (
	"context"
	"mirey7/project-project/internal/data/pro"
	"mirey7/project-project/internal/database"
	"mirey7/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn database.DBConn
}

func (p ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	db := session.Raw("select * from ms_project a, ms_project_member b where a.id=b.project_code and b.member_code=? limit ?, ?", memId, index, pageSize)
	var pms []*pro.ProjectAndMember
	err := db.Scan(&pms).Error
	var total int64
	err = session.Model(&pro.ProjectMember{}).Where("member_code=?", memId).Count(&total).Error
	return pms, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.New(),
	}
}
