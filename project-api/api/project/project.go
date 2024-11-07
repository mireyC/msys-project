package project

import (
	"context"
	"github.com/gin-gonic/gin"
	common "mirey7/project-common"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/project"
	"net/http"
	"time"
)

type HandlerProject struct {
}

func (p HandlerProject) Index(c *gin.Context) {
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

func New() *HandlerProject {
	return &HandlerProject{}
}
