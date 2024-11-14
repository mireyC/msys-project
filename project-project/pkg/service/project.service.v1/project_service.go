package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"mirey7/project-common/encrypts"
	"mirey7/project-common/errs"
	"mirey7/project-common/tms"
	"mirey7/project-grpc/project"
	"mirey7/project-project/internal/data/menu"
	"mirey7/project-project/internal/data/pro"
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
	projectRepo repo.ProjectRepo
}

func New() *ProjectService {
	return &ProjectService{
		cacheRepo:   dao.Rc,
		transaction: dao.NewTransaction(),
		menuRepo:    dao.NewMenuDao(),
		projectRepo: dao.NewProjectDao(),
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

func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectRpcResponse, error) {
	memberId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	selectBy := msg.SelectBy

	var proList []*pro.ProjectAndMember
	var total int64
	var err error
	if selectBy == "" || selectBy == "my" {
		proList, total, err = p.projectRepo.FindProjectByMemId(context.Background(), memberId, page, pageSize, "")
		if err != nil {
			zap.L().Error("project FindProjectByMemId db error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}

	if selectBy == "archive" {
		proList, total, err = p.projectRepo.FindProjectByMemId(context.Background(), memberId, page, pageSize, "and archive = 1 ")
		if err != nil {
			zap.L().Error("project FindProjectByMemId db error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if selectBy == "deleted" {
		proList, total, err = p.projectRepo.FindProjectByMemId(context.Background(), memberId, page, pageSize, "and deleted = 1 ")
		if err != nil {
			zap.L().Error("project FindProjectByMemId db error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if selectBy == "collect" {
		proList, total, err = p.projectRepo.FindCollectProjectByMenId(ctx, memberId, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectByMemId db error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}

	if proList == nil {
		return &project.ProjectRpcResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}
	var pm []*project.ProjectMessage
	pam := pro.ToMap(proList)
	copier.Copy(&pm, proList)
	for _, v := range pm {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)

		v.AccessControlType = pam[v.Id].GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam[v.Id].OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam[v.Id].JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam[v.Id].Sort)
		v.CreateTime = tms.FormatByMill(pam[v.Id].CreateTime)
	}

	return &project.ProjectRpcResponse{Pm: pm, Total: total}, nil
}
