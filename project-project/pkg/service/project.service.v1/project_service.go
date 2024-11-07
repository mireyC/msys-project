package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/project"
	"mirey7/project-project/internal/data/menu"
	"mirey7/project-project/pkg/model"

	"mirey7/project-project/internal/dao"
	"mirey7/project-project/internal/database/tran"
	"mirey7/project-project/internal/repo"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cacheRepo   repo.CacheRepo
	transaction tran.Transaction
	menuRepo    repo.MenuRepo
}

func New() *ProjectService {
	return &ProjectService{
		cacheRepo:   dao.Rc,
		transaction: dao.NewTransaction(),
		menuRepo:    dao.NewMenuDao(),
	}
}

func (p *ProjectService) Index(ctx context.Context, msg *project.IndexMessage) (*project.IndexResponse, error) {
	c := context.Background()
	menus, err := p.menuRepo.FindMenus(c)
	if err != nil {
		zap.L().Error("menuRepo FindMenus db error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	res := &menu.ProjectMenuTree{}

	menu.BuildProjectMenuTree(menus, 0, res)
	//jsonData, _ := json.Marshal(res)
	//log.Println("res json: ", string(jsonData))

	var mms []*project.MenuMessage
	copier.Copy(&mms, res.Children)
	return &project.IndexResponse{
		Menus: mms,
	}, nil

}
