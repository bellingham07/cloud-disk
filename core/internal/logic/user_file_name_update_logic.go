package logic

import (
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest, userIdentity string) (resp *types.UserFileNameUpdateReply, err error) {
	// 判断当前名称是否已经在该目录下存在
	count, err := l.svcCtx.Engine.Where("name=? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity=?)", req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该名称已存在")
	}
	// 直接进行文件修改
	data := &models.UserRepository{Name: req.Name}
	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Update(data)
	if err != nil {
		return nil, err
	}
	return
}
