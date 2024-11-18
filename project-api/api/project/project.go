package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"mirey7/project-api/pkg/model"
	"mirey7/project-api/pkg/model/pro"
	common "mirey7/project-common"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/project"
	"net/http"
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

func New() *HandlerProject {
	return &HandlerProject{}
}
