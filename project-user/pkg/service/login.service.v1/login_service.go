package login_service_v1

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"log"
	common "mirey7/project-common"
	"mirey7/project-common/encrypts"
	"mirey7/project-common/errs"
	"mirey7/project-grpc/user/login"
	"mirey7/project-user/config"
	"mirey7/project-user/internal/dao"
	"mirey7/project-user/internal/data/member"
	"mirey7/project-user/internal/data/organization"
	"mirey7/project-user/internal/repo"
	"mirey7/project-user/pkg/model"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cacheRepo        repo.CacheRepo
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
}

func New() *LoginService {
	return &LoginService{
		cacheRepo:        dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganization(),
	}
}

func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	log.Printf("msg %v \n", msg)
	// 1. 获取参数，校验参数
	c := context.Background()
	redisCode, err := ls.cacheRepo.Get(c, model.RegisterRedisKey+msg.Mobile)

	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}

	if err != nil {
		zap.L().Error("Register redis get error ", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}

	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	// 2. 校验业务逻辑（邮箱是否被注册，手机号是否被注册，昵称是否已存在）
	exits, err := ls.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("db error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exits {
		return nil, errs.GrpcError(model.EmailExist)
	}

	exits, err = ls.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("db error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exits {
		return nil, errs.GrpcError(model.MobileExist)
	}

	exits, err = ls.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("db error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exits {
		return nil, errs.GrpcError(model.AccountExist)
	}

	// 3. 执行业务 将数据存入 member表， 生成 一个数据 存入组织表 organization
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = ls.memberRepo.SaveMember(c, mem)
	if err != nil {
		zap.L().Error("db SaveMember error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	// 存入组织
	org := &organization.Organization{
		Name:       mem.Name + "个人组织",
		MemberId:   mem.Id,
		CreateTime: time.Now().UnixMilli(),
		Personal:   int32(model.Personal),
		Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
	}
	err = ls.organizationRepo.SaveOrganization(c, org)
	if err != nil {
		zap.L().Error("register SaveOrganization db err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &login.RegisterResponse{}, nil
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	// 1. 获取参数
	mobile := msg.Mobile
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	// 3. 生成验证码下·
	code := config.C.GC.Addr

	// 4，调用短信平台（）
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("短信平台调用成功，发送短信")
		//zap.L().Info("短信平台调用成功，发送短信")
		// redis 假设后续可能存在 mongo 当中，也可能存在 memcache 当中
		// 5. 存储验证码 redis 中， 15分钟
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		er := ls.cacheRepo.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if er != nil {
			//log.Printf("验证码存入 redis 失败，cause by: %v \n", er)
			zap.L().Error("Put redis Captcha error ", zap.Error(er))
		}
	}()

	return &login.CaptchaResponse{Code: code}, nil
}
