package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"mirey7/project-api/pkg/model"
	"mirey7/project-api/pkg/model/pro"
	common "mirey7/project-common"
	"mirey7/project-common/encrypts"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/project"
	"net/http"
	"strconv"
	"time"
)

type HandlerProject struct {
}

func (p *HandlerProject) Index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msg := &project.IndexMessage{}
	indexResponse, err := ProjectSvcClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var menus []*pro.Menu
	copier.Copy(&menus, indexResponse.Menus)
	c.JSON(http.StatusOK, result.Success(menus))
}

func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		Page:       page.Page,
		PageSize:   page.PageSize,
		SelectBy:   selectBy}
	projectRpcResponse, err := ProjectSvcClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 坑：
	if projectRpcResponse.Pm == nil {
		projectRpcResponse.Pm = []*project.ProjectMessage{}
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, projectRpcResponse.Pm)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": projectRpcResponse.Total,
	}))
}

func (p *HandlerProject) projectTemplate(c *gin.Context) {
	result := common.Result{}
	//var total int
	req := &pro.ProjectTemplateQueryReq{}
	c.ShouldBind(req)
	//list := []pro.ProjectTemplate{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := ProjectSvcClient.FindProjectTemplateList(ctx, &project.ProjectTemplateMessage{
		Page:           int64(req.Page),
		PageSize:       int64(req.PageSize),
		ViewType:       req.ViewType,
		MemberId:       c.GetInt64("memberId"),
		OrganizationId: c.GetInt64("organizationId"),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	var list []*pro.ProjectTemplate
	copier.Copy(&list, resp.Pts)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": resp.Total,
		"list":  list,
		"page":  req.Page,
	}))
	return
}

func (p *HandlerProject) projectSave(c *gin.Context) {
	result := common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationId := c.GetString("organizationId")
	var req *pro.ProjectSaveReq
	c.ShouldBind(&req)
	msg := &project.ProjectSaveRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationId,
		TemplateCode:     req.TemplateCode,
		Description:      req.Description,
		Id:               int64(req.Id),
		Name:             req.Name,
	}

	saveProject, err := ProjectSvcClient.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	var resp *pro.ProjectSaveResp
	copier.Copy(&resp, saveProject)
	c.JSON(http.StatusOK, result.Success(resp))

}

func (p *HandlerProject) projectRead(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.ProjectMessage{ProjectCode: projectCode, MemberId: int32(c.GetInt64("memberId"))}
	res, err := ProjectSvcClient.ReadProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	resp := &pro.ProjectDetail{}
	copier.Copy(resp, res)
	resp.Code = projectCode
	c.JSON(http.StatusOK, result.Success(resp))
}

func (p *HandlerProject) projectRecycle(c *gin.Context) {
	result := &common.Result{}
	idmsg := c.PostForm("projectCode")
	projectCodeStr, _ := encrypts.Decrypt(idmsg, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectSvcClient.UpdateDeleteProject(ctx, &project.ProjectMessage{Id: projectCode, IsDeleted: 1})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success([]int{}))
	return
}

func (p *HandlerProject) delProject(c *gin.Context) {
	result := &common.Result{}
	projectCodeStr, _ := encrypts.Decrypt(c.PostForm("projectCode"), model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectSvcClient.DelProject(ctx, &project.ProjectMessage{Id: projectCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) projectCollect(c *gin.Context) {
	result := &common.Result{}

	memberId := c.GetInt64("memberId")
	projectCodeStr, _ := encrypts.Decrypt(c.PostForm("projectCode"), model.AESKey)
	//projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collteType := c.PostForm("type")
	var colleted int32
	if collteType == "collect" {
		colleted = 1
	} else {
		colleted = 0
	}
	msg := &project.ProjectMessage{MemberCode: memberId, ProjectCode: projectCodeStr, IsCollected: colleted}
	_, err := ProjectSvcClient.ProjectCollect(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) projectRecovery(c *gin.Context) {
	result := &common.Result{}
	projectCodestr, _ := encrypts.Decrypt(c.PostForm("projectCode"), model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodestr, 10, 64)

	msg := &project.ProjectMessage{
		Id:        projectCode,
		IsDeleted: 0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectSvcClient.UpdateDeleteProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) projectEdit(c *gin.Context) {
	result := &common.Result{}
	req := &pro.Project{}
	_ = c.ShouldBind(req)

	msg := &project.ProjectMessage{}
	copier.Copy(msg, req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectSvcClient.ProjectEdit(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success([]int{}))
}

func New() *HandlerProject {
	return &HandlerProject{}
}
