package dao

import (
	"context"
	"mirey7/project-project/internal/data/menu"
	"mirey7/project-project/internal/database"
	"mirey7/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn database.DBConn
}

func (m *MenuDao) FindMenus(ctx context.Context) ([]*menu.ProjectMenu, error) {
	var pms []*menu.ProjectMenu
	err := m.conn.Session(ctx).Order("pid,sort asc, id asc").Find(&pms).Error
	return pms, err
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		conn: gorms.New(),
	}
}
