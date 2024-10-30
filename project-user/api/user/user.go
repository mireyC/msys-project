package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	common "mirey7/project-common"
	"mirey7/project-user/pkg/dao"
	"mirey7/project-user/pkg/model"
	"mirey7/project-user/pkg/repo"
	"net/http"
	"time"
)

type HandlerUser struct {
	cache repo.Cache
}

func New() *HandlerUser {
	return &HandlerUser{
		cache: dao.Rc,
	}
}

func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	resp := &common.Result{}
	// 1. 获取参数
	mobile := ctx.PostForm("mobile")
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, resp.Fail(model.NoLegalMobile, "手机号不合法"))
		return
	}
	// 3. 生成验证码
	code := "123123"
	// 4，调用短信平台（）
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("短信平台调用成功，发送短信")

		// redis 假设后续可能存在 mongo 当中，也可能存在 memcache 当中
		// 5. 存储验证码 redis 中， 15分钟
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		er := h.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if er != nil {
			log.Printf("验证码存入 redis 失败，cause by: %v \n", er)
		}
	}()

	ctx.JSON(200, resp.Success(code))
}
