package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {
	// 判断改邮箱未被注册
	count, err := l.svcCtx.Engine.Where("email=?", req.Email).Count(new(models.UserBasic))
	if err != nil {
		return
	}
	if count > 0 {
		// 邮箱已存在
		err = errors.New("此邮箱已被注册")
		return
	}
	// 未被注册 生成验证码
	code := helper.RandCode()
	// 存储验证码
	key := define.LoginCodePrefix + req.Email
	// 五分钟有效
	l.svcCtx.RDB.Set(l.ctx, key, code, define.CodeExpireTime)
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		return nil, err
	}
	return nil, err
}
