package login_service_v1

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	common "mirey7/project-common"
	"mirey7/project-common/encrypts"
	"mirey7/project-common/errs"
	"mirey7/project-common/jwts"
	"mirey7/project-grpc/user/login"
	"mirey7/project-user/config"
	"mirey7/project-user/internal/dao"
	"mirey7/project-user/internal/data/member"
	"mirey7/project-user/internal/data/organization"
	"mirey7/project-user/internal/database"
	"mirey7/project-user/internal/database/tran"
	"mirey7/project-user/internal/repo"
	"mirey7/project-user/pkg/model"
	"strconv"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cacheRepo        repo.CacheRepo
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cacheRepo:        dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganization(),
		transaction:      dao.NewTransaction(),
	}
}

func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	//log.Printf("msg %v \n", msg)
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

	log.Println("Register password", msg.Password)
	// 1.去数据库查询 账号密码是否正确
	//pwd := encrypts.Md5(msg.Password)
	log.Println("Register password", pwd)

	err = ls.transaction.Action(func(conn database.DBConn) error {

		err = ls.memberRepo.SaveMember(conn, c, mem)
		if err != nil {
			zap.L().Error("register SaveMember db error ", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}

		// 存入组织
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   int32(model.Personal),
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}

		return nil
	})
	if err != nil {
		return &login.RegisterResponse{}, err
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

func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()
	log.Println("Login password", msg.Password)
	// 1.去数据库查询 账号密码是否正确
	pwd := encrypts.Md5(msg.Password)
	log.Println("Login password", pwd)
	mem, err := ls.memberRepo.FindMember(c, msg.Account, pwd)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountAndPwdError)
	}

	memMessage := &login.MemberMessage{}
	err = copier.Copy(memMessage, mem)
	memMessage.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)
	//memMessage.Code = str
	// 2.根据用户id 查询组织
	orgs, err := ls.organizationRepo.FindOrganizationByMenId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login db FindOrganizationByMenId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if orgs == nil {
		return nil, errs.GrpcError(model.OrganizationNoFound)
	}

	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, &orgs)
	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
	}

	// 3.用 jwt 生成 token
	memId := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.C.JC.AccessExp*3600*24) * time.Second
	refreshExp := time.Duration(config.C.JC.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memId, exp, refreshExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}

	return &login.LoginResponse{
		Member:           memMessage,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}

func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if token == "" {
		return nil, errs.GrpcError(model.NoLogin)
	}
	parseToken, err := jwts.ParseJwt(token, config.C.JC.AccessSecret)
	if err != nil {
		zap.L().Error("TokenVerify ParseToken err", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	memId, err := strconv.ParseInt(parseToken, 10, 64)
	if err != nil {
		zap.L().Error("TokenVerify ParseInt err", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	c := context.Background()
	mem, err := ls.memberRepo.FindMemberById(c, memId)
	if err != nil {
		zap.L().Error("Login db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if mem == nil {
		zap.L().Error("TokenVerify member is nil")
		return nil, errs.GrpcError(model.NoLogin)
	}
	memMsg := &login.MemberMessage{}
	_ = copier.Copy(memMsg, mem)
	return &login.LoginResponse{Member: memMsg}, nil
}
