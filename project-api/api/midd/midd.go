package midd

import (
	"context"
	"github.com/gin-gonic/gin"
	"mirey7/project-api/api/user"
	common "mirey7/project-common"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/user/login"
	"net/http"
)

func TokenVerify() func(*gin.Context) {
	return func(c *gin.Context) {
		// 1. 从 header 中获取 token
		token := c.GetHeader("Authorization")
		// 2. 调用 user 服务 进行token 认证
		//ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

		resp, err := user.LoginServiceClient.TokenVerify(context.Background(), &login.LoginMessage{Token: token})
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			result := common.Result{}
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}
		member := resp.Member

		// 3. 处理结果， 认证通过  将信息放入上下文， 失败返回未登录
		c.Set("memberId", member.Id)
		c.Next()
	}
}
