package repo

import (
	"context"
	"mirey7/project-project/internal/data/menu"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*menu.ProjectMenu, error)
}
