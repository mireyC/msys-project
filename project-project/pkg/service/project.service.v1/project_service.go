package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	"mirey7/project-common/encrypts"
	"mirey7/project-common/errs"
	"mirey7/project-common/tms"
	"mirey7/project-grpc/project"
	"mirey7/project-project/internal/data/menu"
	"mirey7/project-project/internal/data/pro"
	"mirey7/project-project/internal/data/task"
	"mirey7/project-project/internal/database"
	"mirey7/project-project/pkg/model"
	"strconv"
	"time"

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
	taskRepo    repo.TaskRepo
}

func New() *ProjectService {
	return &ProjectService{
		cacheRepo:   dao.Rc,
		transaction: dao.NewTransaction(),
		menuRepo:    dao.NewMenuDao(),
		projectRepo: dao.NewProjectDao(),
		taskRepo:    dao.NewTaskDao(),
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

func (p *ProjectService) FindProjectTemplateList(ctx context.Context, msg *project.ProjectTemplateMessage) (*project.ProjectTemplateResp, error) {
	viewType := msg.ViewType
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var pts []*pro.ProjectTemplate
	var total int64
	var err error
	if viewType == "1" { // 系统模板
		pts, total, err = p.projectRepo.FindProjectTemplateSystem(c, msg.Page, msg.PageSize, 1)
	}

	if viewType == "0" { // 自定义模板
		pts, total, err = p.projectRepo.FindProjectTemplateCustom(c, msg.Page, msg.PageSize, msg.OrganizationId, msg.MemberId)
	}

	if viewType == "-1" { // -1 所有模板
		pts, total, err = p.projectRepo.FindProjectTemplateAll(c, msg.Page, msg.PageSize, msg.OrganizationId)
	}

	if err != nil {
		zap.L().Error("project FindProjectTemplateList db error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	if pts == nil {
		return &project.ProjectTemplateResp{}, nil
	}

	ids := pro.ToProjectTemplateIds(pts)
	//查询task stages数据库
	taskStages, err := p.taskRepo.FindTaskByIds(c, ids)
	if err != nil {
		zap.L().Error("project FindProjectTemplateList FindTaskByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	tm := task.CovertProjectMap(taskStages)
	var proTemplates []*pro.ProjectTemplateAll
	for _, v := range pts {
		proTemplates = append(proTemplates, v.Convert(tm[v.Id]))
	}

	var proTemplateList []*project.ProjectTemplate
	copier.Copy(&proTemplateList, proTemplates)

	return &project.ProjectTemplateResp{
		Pts:   proTemplateList,
		Total: total,
	}, nil
}

func (ps *ProjectService) SaveProject(ctx context.Context, msg *project.ProjectSaveRpcMessage) (*project.ProjectSaveRespMessage, error) {

	organizationCode, _ := strconv.ParseInt(msg.OrganizationCode, 10, 64)
	templateCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, model.AESKey)
	templateCode, _ := strconv.ParseInt(templateCodeStr, 10, 64)
	pr := &pro.Project{
		Name:              msg.Name,
		Description:       msg.Description,
		TemplateCode:      int(templateCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://img2.baidu.com/it/u=792555388,2449797505&fm=253&fmt=auto&app=138&f=JPEG?w=667&h=500",
		Deleted:           model.NoDeleted,
		Archive:           model.NoArchive,
		OrganizationCode:  organizationCode,
		AccessControlType: model.Open,
		TaskBoardTheme:    model.Simple,
	}
	var rsp *project.ProjectSaveRespMessage
	err := ps.transaction.Action(func(conn database.DBConn) error {
		err := ps.projectRepo.SaveProject(ctx, conn, pr)
		if err != nil {
			zap.L().Error("SaveProject Save error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		pm := &pro.ProjectMember{
			ProjectCode: pr.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   "",
		}
		err = ps.projectRepo.SaveProjectMember(ctx, conn, pm)
		if err != nil {
			zap.L().Error("SaveProject SaveProjectMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		code, _ := encrypts.EncryptInt64(pr.Id, model.AESKey)
		organizationCodeStr, _ := encrypts.EncryptInt64(organizationCode, model.AESKey)
		rsp = &project.ProjectSaveRespMessage{
			Id:               pr.Id,
			Code:             code,
			OrganizationCode: organizationCodeStr,
			Name:             pr.Name,
			Cover:            pr.Cover,
			CreateTime:       tms.FormatByMill(pr.CreateTime),
			TaskBoardTheme:   pr.TaskBoardTheme,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// ReadProject
// 项目信息，是否收藏，姓名，拥有者头像
func (p *ProjectService) ReadProject(c context.Context, msg *project.ProjectMessage) (*project.ProjectMessage, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ctx = context.Background()
	projectAndMember, err := p.projectRepo.FindProjectAndMember(ctx, projectCode, msg.MemberId)
	if err != nil {
		zap.L().Error("ReadProject FindProjectAndMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	ownerId := projectAndMember.IsOwner
	log.Println(ownerId)
	collected, err := p.projectRepo.FindCollectByProjectCodeAndMemberId(ctx, projectCode, msg.MemberId)
	if err != nil {
		zap.L().Error("ReadProject FindCollectByProjectCodeAndMemberId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	resp := &project.ProjectMessage{}
	copier.Copy(resp, projectAndMember)
	resp.OwnerAvatar = ""
	if collected {
		resp.Collected = 1
	}

	return resp, nil
}
