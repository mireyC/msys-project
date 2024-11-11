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
	//log.Println(resp)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//rsp := &project.IndexResponse{}
	//jsonData, err := protojson.Marshal(indexResponse)
	c.JSON(http.StatusOK, result.Success(indexResponse.Menus))
}

func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberIdStr, _ := c.Get("memberId")
	page := model.Page{}
	page.Bind(c)
	memberId := memberIdStr.(int64)
	msg := &project.ProjectRpcMessage{MemberId: memberId, Page: page.Page, PageSize: page.PageSize}
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

func New() *HandlerProject {
	return &HandlerProject{}
}
