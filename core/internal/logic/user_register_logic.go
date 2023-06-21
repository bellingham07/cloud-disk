package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"errors"
	"log"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	// 判断验证码是否一致
	code, err := l.svcCtx.RDB.Get(l.ctx, define.LoginCodePrefix+req.Email).Result()
	if err != nil {
		err = errors.New("验证码错误")
		return nil, err
	}
	if code != req.Code {
		err = errors.New("验证码错误")
		return
	}
	// 判断用户名是否已存在
	count, err := l.svcCtx.Engine.Where("name=?", req.Name).Count(new(models.UserBasic))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		err = errors.New("用户名已存在")
		return
	}
	//一致 将用户信息注册进入数据库
	user := models.UserBasic{
		Identity: helper.GetUUID(),
		Name:     req.Name,
		Password: helper.Md5(req.Password),
		Email:    req.Email,
	}
	insert, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	log.Println("insert user row:", insert)
	return
}
