package login_service_v1

import (
	"context"
	"go.uber.org/zap"
	"log"
	common "mirey7/project-common"
	"mirey7/project-common/errs"
	"mirey7/project-user/config"
	"mirey7/project-user/pkg/dao"
	"mirey7/project-user/pkg/model"
	"mirey7/project-user/pkg/repo"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	cache repo.Cache
}

func New() *LoginService {
	return &LoginService{
		cache: dao.Rc,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {
	// 1. 获取参数
	mobile := msg.Mobile
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	// 3. 生成验证码下·
	code := "123123"
	// 4，调用短信平台（）
	go func() {
		time.Sleep(2 * time.Second)
		//log.Println("短信平台调用成功，发送短信")
		zap.L().Info("短信平台调用成功，发送短信")
		// redis 假设后续可能存在 mongo 当中，也可能存在 memcache 当中
		// 5. 存储验证码 redis 中， 15分钟
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		er := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if er != nil {
			log.Printf("验证码存入 redis 失败，cause by: %v \n", er)
		}
	}()

	return &CaptchaResponse{Code: config.C.GC.Addr}, nil
}
